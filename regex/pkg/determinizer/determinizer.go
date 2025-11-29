package determinizer

import (
	"sort"
	"strings"

	"regex/pkg/model"
)

type Determinizer struct {
	nfa        *model.NFA
	dfa        *model.DFA
	dStates    [][]string
	dStatesMap map[string]int
	queue      []int
}

func NewDeterminizer(nfa *model.NFA) *Determinizer {
	d := &Determinizer{
		nfa:        nfa,
		dfa:        model.NewDFA(),
		dStatesMap: make(map[string]int),
	}

	for _, symbol := range nfa.Alphabet {
		if symbol != model.Epsilon {
			d.dfa.Alphabet = append(d.dfa.Alphabet, symbol)
		}
	}
	sort.Strings(d.dfa.Alphabet)

	d.initializeStartState()
	return d
}

func (d *Determinizer) Run() *model.DFA {
	for len(d.queue) > 0 {
		currentIndex := d.queue[0]
		d.queue = d.queue[1:]
		d.processState(currentIndex)
	}
	d.finalizeDFA()

	return d.dfa
}

func (d *Determinizer) initializeStartState() {
	startSet := epsilonClosure([]string{d.nfa.StartState}, d.nfa)
	name, isNew := d.registerDFAState(startSet)

	d.dfa.StartState = name
	if isNew {
		d.queue = append(d.queue, d.dStatesMap[name])
	}
}

func (d *Determinizer) processState(stateIndex int) {
	T := d.dStates[stateIndex]
	fromStateName := makeStateName(T)

	for _, symbol := range d.dfa.Alphabet {
		moveResult := move(T, symbol, d.nfa)
		if len(moveResult) == 0 {
			continue
		}

		U := epsilonClosure(moveResult, d.nfa)
		toStateName, isNew := d.registerDFAState(U)

		if isNew {
			d.queue = append(d.queue, d.dStatesMap[toStateName])
		}
		d.addTransition(fromStateName, symbol, toStateName)
	}
}

func (d *Determinizer) registerDFAState(nfaStates []string) (string, bool) {
	name := makeStateName(nfaStates)
	if _, exists := d.dStatesMap[name]; exists {
		return name, false
	}

	index := len(d.dStates)
	d.dStatesMap[name] = index
	d.dStates = append(d.dStates, nfaStates)
	return name, true
}

func (d *Determinizer) addTransition(from, symbol, to string) {
	if _, ok := d.dfa.Transitions[from]; !ok {
		d.dfa.Transitions[from] = make(map[string]string)
	}
	d.dfa.Transitions[from][symbol] = to
}

func (d *Determinizer) finalizeDFA() {
	for name, index := range d.dStatesMap {
		d.dfa.States = append(d.dfa.States, name)
		nfaStates := d.dStates[index]

		for _, nfaState := range nfaStates {
			if d.nfa.AcceptingStates[nfaState] {
				d.dfa.AcceptingStates[name] = true
				break
			}
		}
	}
	sort.Strings(d.dfa.States)
}

func epsilonClosure(states []string, nfa *model.NFA) []string {
	closureSet := make(map[string]bool)
	stack := make([]string, 0, len(states))

	for _, state := range states {
		if !closureSet[state] {
			closureSet[state] = true
			stack = append(stack, state)
		}
	}

	for len(stack) > 0 {
		currentState := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if transitions, ok := nfa.Transitions[currentState]; ok {
			if dest, ok := transitions[model.Epsilon]; ok {
				for _, destState := range dest {
					if !closureSet[destState] {
						closureSet[destState] = true
						stack = append(stack, destState)
					}
				}
			}
		}
	}

	result := make([]string, 0, len(closureSet))
	for state := range closureSet {
		result = append(result, state)
	}
	sort.Strings(result)
	return result
}

func move(states []string, symbol string, nfa *model.NFA) []string {
	destinations := make(map[string]bool)
	for _, state := range states {
		if transitions, ok := nfa.Transitions[state]; ok {
			if destStates, ok := transitions[symbol]; ok {
				for _, dest := range destStates {
					destinations[dest] = true
				}
			}
		}
	}

	result := make([]string, 0, len(destinations))
	for state := range destinations {
		result = append(result, state)
	}
	sort.Strings(result)
	return result
}

func makeStateName(states []string) string {
	return strings.Join(states, "_")
}
