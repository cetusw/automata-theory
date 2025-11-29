package model

type NFA struct {
	States          []string
	Alphabet        []string
	Transitions     map[string]map[string][]string
	StartState      string
	AcceptingStates map[string]bool
}

func NewNFA() *NFA {
	return &NFA{
		Transitions:     make(map[string]map[string][]string),
		AcceptingStates: make(map[string]bool),
	}
}

type NfaFragment struct {
	StartState *State
	EndState   *State
}
