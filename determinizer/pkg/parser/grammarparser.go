package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"determinizer/pkg/model"
)

const (
	finalStateName = "H"
	newStartState  = "F"
)

type production struct {
	terminal      string
	nonTerminal   string
	isRightLinear bool
}

type grammarParser struct {
	rules         map[string][]production
	nonTerminals  map[string]bool
	startSymbol   string
	hasLeftLinear bool
}

func ParseGrammarToNFA(grammarString string) (*model.NFA, error) {
	p := &grammarParser{
		rules:        make(map[string][]production),
		nonTerminals: make(map[string]bool),
	}

	if err := p.parseAndAnalyze(grammarString); err != nil {
		return nil, err
	}

	return p.buildNFA()
}

func (p *grammarParser) parseAndAnalyze(grammarString string) error {
	preScanner := bufio.NewScanner(strings.NewReader(grammarString))
	for preScanner.Scan() {
		line := strings.TrimSpace(preScanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "->")
		if len(parts) == 2 {
			nonTerminal := strings.TrimSpace(parts[0])
			p.nonTerminals[nonTerminal] = true
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(grammarString))
	isFirstRule := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if isFirstRule {
			parts := strings.Split(line, "->")
			if len(parts) == 2 {
				p.startSymbol = strings.TrimSpace(parts[0])
				isFirstRule = false
			}
		}

		if err := p.parseLine(line); err != nil {
			return err
		}
	}

	if p.startSymbol == "" {
		return errors.New("empty grammar")
	}
	return nil
}

func (p *grammarParser) parseLine(line string) error {
	parts := strings.Split(line, "->")
	if len(parts) != 2 {
		return fmt.Errorf("некорректная строка: %s", line)
	}

	nonTerminal := strings.TrimSpace(parts[0])

	productionsStr := strings.Split(parts[1], "|")
	for _, prodStr := range productionsStr {
		if err := p.addProduction(nonTerminal, strings.TrimSpace(prodStr)); err != nil {
			return err
		}
	}
	return nil
}

func (p *grammarParser) addProduction(from, prodStr string) error {
	if prodStr == "eps" {
		p.rules[from] = append(p.rules[from], production{terminal: "eps"})
		return nil
	}
	if len(prodStr) == 1 {
		p.rules[from] = append(p.rules[from], production{terminal: prodStr, isRightLinear: true})
		return nil
	}
	if len(prodStr) != 2 {
		return fmt.Errorf("rule must have 1 or 2 symbols: %s -> %s", from, prodStr)
	}

	if p.isNonTerminal(string(prodStr[1])) {
		p.rules[from] = append(p.rules[from], production{
			terminal:      string(prodStr[0]),
			nonTerminal:   string(prodStr[1]),
			isRightLinear: true,
		})
	} else if p.isNonTerminal(string(prodStr[0])) {
		p.hasLeftLinear = true
		p.rules[from] = append(p.rules[from], production{
			terminal:      string(prodStr[1]),
			nonTerminal:   string(prodStr[0]),
			isRightLinear: false,
		})
	} else {
		return fmt.Errorf("invalid rule: %s -> %s", from, prodStr)
	}
	return nil
}

func (p *grammarParser) isNonTerminal(s string) bool {
	_, exists := p.nonTerminals[s]
	return exists
}

func (p *grammarParser) buildNFA() (*model.NFA, error) {
	if p.hasLeftLinear {
		return p.buildNFAForLeftLinear()
	}
	return p.buildNFAForRightLinear()
}

func (p *grammarParser) buildNFAForRightLinear() (*model.NFA, error) {
	nfa := model.NewNFA()
	nfa.StartState = p.startSymbol
	nfa.AcceptingStates[finalStateName] = true
	alphabetSet := make(map[string]bool)

	for from, productions := range p.rules {
		for _, prod := range productions {
			if prod.terminal == "eps" {
				nfa.AcceptingStates[from] = true
				continue
			}

			to := prod.nonTerminal
			if to == "" {
				to = finalStateName
			}
			if _, ok := nfa.Transitions[from]; !ok {
				nfa.Transitions[from] = make(map[string][]string)
			}
			nfa.Transitions[from][prod.terminal] = append(nfa.Transitions[from][prod.terminal], to)

			alphabetSet[prod.terminal] = true
		}
	}
	for nt := range p.nonTerminals {
		nfa.States = append(nfa.States, nt)
	}
	nfa.States = append(nfa.States, finalStateName)

	for term := range alphabetSet {
		nfa.Alphabet = append(nfa.Alphabet, term)
	}

	return nfa, nil
}

func (p *grammarParser) buildNFAForLeftLinear() (*model.NFA, error) {
	nfa := model.NewNFA()
	nfa.StartState = newStartState
	nfa.AcceptingStates[p.startSymbol] = true
	alphabetSet := make(map[string]bool)

	for from, productions := range p.rules {
		for _, prod := range productions {
			if prod.terminal == "eps" {
				nfa.AcceptingStates[from] = true
				continue
			}

			if prod.isRightLinear && prod.nonTerminal != "" {
				return nil, fmt.Errorf("смешанные грамматики (с правилами вида aB и Ba) не поддерживаются")
			}

			fromState := prod.nonTerminal
			toState := from
			if fromState == "" {
				fromState = newStartState
			}

			if _, ok := nfa.Transitions[fromState]; !ok {
				nfa.Transitions[fromState] = make(map[string][]string)
			}
			nfa.Transitions[fromState][prod.terminal] = append(nfa.Transitions[fromState][prod.terminal], toState)

			alphabetSet[prod.terminal] = true
		}
	}

	for nt := range p.nonTerminals {
		nfa.States = append(nfa.States, nt)
	}
	nfa.States = append(nfa.States, newStartState)

	for term := range alphabetSet {
		nfa.Alphabet = append(nfa.Alphabet, term)
	}

	return nfa, nil
}
