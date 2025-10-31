package parser

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"determinizer/pkg/model"
)

var (
	acceptingStateRegex = regexp.MustCompile(`node\s*\[\s*shape\s*=\s*doublecircle\s*\];\s*([^;]+);`)
	startStateRegex     = regexp.MustCompile(`start\s*->\s*(\w+);`)
	transitionRegex     = regexp.MustCompile(`([\w_]+)\s*->\s*([\w_]+)\s*\[\s*label\s*=\s*"([^"]+)"\s*\];`)
)

type parser struct {
	startState      string
	acceptingStates map[string]bool
	transitions     map[string]map[string][]string
	allStates       map[string]bool
	alphabetSet     map[string]bool
}

func newParser() *parser {
	return &parser{
		acceptingStates: make(map[string]bool),
		transitions:     make(map[string]map[string][]string),
		allStates:       make(map[string]bool),
		alphabetSet:     make(map[string]bool),
	}
}

func ParseNFA(dotString string) (*model.NFA, error) {
	p := newParser()
	if err := p.parse(dotString); err != nil {
		return nil, err
	}
	return p.buildNFA()
}

func ParseDFA(dotString string) (*model.DFA, error) {
	p := newParser()
	if err := p.parse(dotString); err != nil {
		return nil, err
	}
	return p.buildDFA()
}

func (p *parser) parse(dotString string) error {
	var fullText strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(dotString))
	for scanner.Scan() {
		fullText.WriteString(scanner.Text())
		fullText.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read string: %w", err)
	}

	text := fullText.String()
	p.parseAcceptingStates(text)
	p.parseAllTransitions(text)

	if err := p.parseStartState(text); err != nil {
		return err
	}

	return nil
}

func (p *parser) parseAcceptingStates(text string) {
	matches := acceptingStateRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return
	}
	for _, state := range strings.Fields(matches[1]) {
		p.acceptingStates[state] = true
		p.allStates[state] = true
	}
}

func (p *parser) parseStartState(text string) error {
	matches := startStateRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return errors.New("start state not found")
	}
	p.startState = matches[1]
	p.allStates[p.startState] = true
	return nil
}

func (p *parser) parseAllTransitions(text string) {
	matches := transitionRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		from, to, symbol := match[1], match[2], match[3]

		if _, ok := p.transitions[from]; !ok {
			p.transitions[from] = make(map[string][]string)
		}
		p.transitions[from][symbol] = append(p.transitions[from][symbol], to)

		p.allStates[from] = true
		p.allStates[to] = true
		p.alphabetSet[symbol] = true
	}
}

func (p *parser) buildNFA() (*model.NFA, error) {
	if p.startState == "" {
		return nil, errors.New("cannot finalize DFA: start state not found")
	}

	nfa := model.NewNFA()
	nfa.StartState = p.startState
	nfa.AcceptingStates = p.acceptingStates
	nfa.Transitions = p.transitions

	for state := range p.allStates {
		nfa.States = append(nfa.States, state)
	}
	sort.Strings(nfa.States)

	for symbol := range p.alphabetSet {
		nfa.Alphabet = append(nfa.Alphabet, symbol)
	}
	sort.Strings(nfa.Alphabet)

	return nfa, nil
}

func (p *parser) buildDFA() (*model.DFA, error) {
	if p.startState == "" {
		return nil, errors.New("cannot finalize DFA: start state not found")
	}

	dfa := model.NewDFA()
	dfa.StartState = p.startState
	dfa.AcceptingStates = p.acceptingStates

	for from, transitions := range p.transitions {
		dfa.Transitions[from] = make(map[string]string)
		for symbol, toStates := range transitions {
			if len(toStates) > 1 {
				return nil, fmt.Errorf("failed to parse DFA: nondeterministic transition from '%s' by symbol '%s'", from, symbol)
			}
			dfa.Transitions[from][symbol] = toStates[0]
		}
	}

	for state := range p.allStates {
		dfa.States = append(dfa.States, state)
	}
	sort.Strings(dfa.States)

	for symbol := range p.alphabetSet {
		dfa.Alphabet = append(dfa.Alphabet, symbol)
	}
	sort.Strings(dfa.Alphabet)

	return dfa, nil
}
