package model

type DFA struct {
	States          []string
	Alphabet        []string
	Transitions     map[string]map[string]string
	StartState      string
	AcceptingStates map[string]bool
}

func NewDFA() *DFA {
	return &DFA{
		Transitions:     make(map[string]map[string]string),
		AcceptingStates: make(map[string]bool),
	}
}
