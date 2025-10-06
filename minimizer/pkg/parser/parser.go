package parser

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"minimisation/pkg/model"
)

var (
	acceptingStateRegex = regexp.MustCompile(`node\s*\[\s*shape\s*=\s*doublecircle\s*\];\s*([^;]+);`)
	startStateRegex     = regexp.MustCompile(`start\s*->\s*(\w+);`)
	transitionRegex     = regexp.MustCompile(`(\w+)\s*->\s*(\w+)\s*\[\s*label\s*=\s*"([^"]+)"\s*\];`)
)

type Parser struct {
	dfa         *model.DFA
	allStates   map[string]bool
	alphabetSet map[string]bool
	scanner     *bufio.Scanner
}

func NewParser(dotString string) *Parser {
	return &Parser{
		dfa:         model.NewDFA(),
		allStates:   make(map[string]bool),
		alphabetSet: make(map[string]bool),
		scanner:     bufio.NewScanner(strings.NewReader(dotString)),
	}
}

func (p *Parser) Parse() (*model.DFA, error) {
	for p.scanner.Scan() {
		p.processLine(strings.TrimSpace(p.scanner.Text()))
	}
	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строки: %w", err)
	}
	return p.finalizeDFA()
}

func (p *Parser) processLine(line string) {
	if p.parseAcceptingStates(line) {
		return
	}
	if p.parseStartState(line) {
		return
	}
	p.parseTransition(line)
}

func (p *Parser) parseAcceptingStates(line string) bool {
	matches := acceptingStateRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return false
	}
	for _, state := range strings.Fields(matches[1]) {
		p.dfa.AcceptingStates[state] = true
		p.allStates[state] = true
	}

	return true
}

func (p *Parser) parseStartState(line string) bool {
	matches := startStateRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return false
	}
	p.dfa.StartState = matches[1]
	p.allStates[p.dfa.StartState] = true

	return true
}

func (p *Parser) parseTransition(line string) bool {
	matches := transitionRegex.FindStringSubmatch(line)
	if len(matches) < 4 {
		return false
	}
	from, to, symbol := matches[1], matches[2], matches[3]

	if _, ok := p.dfa.Transitions[from]; !ok {
		p.dfa.Transitions[from] = make(map[string]string)
	}
	p.dfa.Transitions[from][symbol] = to

	p.allStates[from] = true
	p.allStates[to] = true
	p.alphabetSet[symbol] = true

	return true
}

func (p *Parser) finalizeDFA() (*model.DFA, error) {
	if p.dfa.StartState == "" {
		return nil, errors.New("начальное состояние не найдено")
	}
	for state := range p.allStates {
		p.dfa.States = append(p.dfa.States, state)
	}
	sort.Strings(p.dfa.States)
	for symbol := range p.alphabetSet {
		p.dfa.Alphabet = append(p.dfa.Alphabet, symbol)
	}
	sort.Strings(p.dfa.Alphabet)

	return p.dfa, nil
}
