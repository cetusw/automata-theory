package regex

import (
	"fmt"

	"regex/pkg/model"
)

type Converter struct {
	stateCounter int
	stack        []*model.NfaFragment
}

func NewConverter() *Converter {
	return &Converter{
		stateCounter: 0,
		stack:        make([]*model.NfaFragment, 0),
	}
}

func (c *Converter) newState() *model.State {
	s := model.NewState(c.stateCounter)
	c.stateCounter++
	return s
}

func (c *Converter) ConvertToNFA(postfix string) (*model.NFA, error) {
	if postfix == "" {
		start := c.newState()
		start.IsAccepting = true
		return c.buildFinalNFA(start), nil
	}

	for _, r := range postfix {
		var err error
		switch r {
		case '.':
			err = c.handleConcatenation()
		case '|':
			err = c.handleAlternation()
		case '*':
			err = c.handleKleenStar()
		case '+':
			err = c.handleKleenPlus()
		default:
			err = c.handleOperand(r)
		}
		if err != nil {
			return nil, err
		}
	}

	if len(c.stack) != 1 {
		return nil, fmt.Errorf("error: stack must contain one NFA fragment, but contains %d: %v")
	}

	finalFragment := c.stack[0]
	finalFragment.EndState.IsAccepting = true
	return c.buildFinalNFA(finalFragment.StartState), nil
}

func (c *Converter) handleOperand(r rune) error {
	start := c.newState()
	end := c.newState()
	start.AddTransition(string(r), end)
	c.stack = append(c.stack, &model.NfaFragment{StartState: start, EndState: end})
	return nil
}

func (c *Converter) handleConcatenation() error {
	if len(c.stack) < 2 {
		return fmt.Errorf("concat error: not enough operands (at least 2 required): %v")
	}
	frag2 := c.stack[len(c.stack)-1]
	frag1 := c.stack[len(c.stack)-2]
	c.stack = c.stack[:len(c.stack)-2]

	frag1.EndState.AddTransition(model.Epsilon, frag2.StartState)
	frag1.EndState.IsAccepting = false
	c.stack = append(c.stack, &model.NfaFragment{StartState: frag1.StartState, EndState: frag2.EndState})
	return nil
}

func (c *Converter) handleAlternation() error {
	if len(c.stack) < 2 {
		return fmt.Errorf("model.Epsilon error: not enough operands (at least 2 required): %v")
	}
	frag2 := c.stack[len(c.stack)-1]
	frag1 := c.stack[len(c.stack)-2]
	c.stack = c.stack[:len(c.stack)-2]

	newStart := c.newState()
	newStart.AddTransition(model.Epsilon, frag1.StartState)
	newStart.AddTransition(model.Epsilon, frag2.StartState)

	newEnd := c.newState()
	frag1.EndState.AddTransition(model.Epsilon, newEnd)
	frag1.EndState.IsAccepting = false
	frag2.EndState.AddTransition(model.Epsilon, newEnd)
	frag2.EndState.IsAccepting = false
	c.stack = append(c.stack, &model.NfaFragment{StartState: newStart, EndState: newEnd})
	return nil
}

func (c *Converter) handleKleenStar() error {
	if len(c.stack) < 1 {
		return fmt.Errorf("kleen star error: not enough operands (at least 1 required): %v")
	}
	frag := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]

	newStart := c.newState()
	newEnd := c.newState()
	newStart.AddTransition(model.Epsilon, frag.StartState)
	newStart.AddTransition(model.Epsilon, newEnd)

	frag.EndState.AddTransition(model.Epsilon, frag.StartState)
	frag.EndState.AddTransition(model.Epsilon, newEnd)
	frag.EndState.IsAccepting = false
	c.stack = append(c.stack, &model.NfaFragment{StartState: newStart, EndState: newEnd})
	return nil
}

func (c *Converter) handleKleenPlus() error {
	if len(c.stack) < 1 {
		return fmt.Errorf("kleen plus error: not enough operands (at least 1 required): %v")
	}
	frag := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]

	newStart := c.newState()
	newEnd := c.newState()
	newStart.AddTransition(model.Epsilon, frag.StartState)

	frag.EndState.AddTransition(model.Epsilon, frag.StartState)
	frag.EndState.AddTransition(model.Epsilon, newEnd)
	frag.EndState.IsAccepting = false
	c.stack = append(c.stack, &model.NfaFragment{StartState: newStart, EndState: newEnd})
	return nil
}

func (c *Converter) buildFinalNFA(startState *model.State) *model.NFA {
	nfa := &model.NFA{
		Transitions:     make(map[string]map[string][]string),
		AcceptingStates: make(map[string]bool),
		StartState:      startState.GetID(),
	}
	alphabetSet := make(map[string]bool)

	visited := make(map[*model.State]bool)
	queue := []*model.State{startState}
	visited[startState] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		nfa.States = append(nfa.States, current.GetID())
		if current.IsAccepting {
			nfa.AcceptingStates[current.GetID()] = true
		}

		c.collectTransitions(current, nfa, alphabetSet, visited, &queue)
	}

	for alpha := range alphabetSet {
		nfa.Alphabet = append(nfa.Alphabet, alpha)
	}
	return nfa
}

func (c *Converter) collectTransitions(
	current *model.State,
	nfa *model.NFA,
	alphabetSet map[string]bool,
	visited map[*model.State]bool,
	queue *[]*model.State,
) {
	stateID := current.GetID()
	nfa.Transitions[stateID] = make(map[string][]string)

	for symbol, nextStates := range current.Transitions {
		if symbol != model.Epsilon {
			alphabetSet[symbol] = true
		}
		for _, next := range nextStates {
			nfa.Transitions[stateID][symbol] = append(nfa.Transitions[stateID][symbol], next.GetID())
			if !visited[next] {
				visited[next] = true
				*queue = append(*queue, next)
			}
		}
	}
}
