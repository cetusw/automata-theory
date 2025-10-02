package model

type MooreState struct {
	Name   string
	Output string
}

type MooreMachine struct {
	States      map[string]MooreState
	Transitions map[string]map[string]string
}
