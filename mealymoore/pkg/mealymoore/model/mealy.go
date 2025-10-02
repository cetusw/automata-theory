package model

type MealyTransition struct {
	DestinationState string
	Output           string
}

type MealyMachine struct {
	States      map[string]bool
	Transitions map[string]map[string]MealyTransition
}
