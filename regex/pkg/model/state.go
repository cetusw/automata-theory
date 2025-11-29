package model

import "strconv"

type State struct {
	ID          int
	IsAccepting bool
	Transitions map[string][]*State
}

func NewState(id int) *State {
	s := &State{
		ID:          id,
		IsAccepting: false,
		Transitions: make(map[string][]*State),
	}
	return s
}

func (s *State) AddTransition(symbol string, to *State) {
	s.Transitions[symbol] = append(s.Transitions[symbol], to)
}

func (s *State) GetID() string {
	return "S" + strconv.Itoa(s.ID)
}
