package mealymoore

import (
	"fmt"
	"mealymoore/pkg/mealymoore/model"
)

func MooreToMealy(moore *model.MooreMachine) *model.MealyMachine {
	mealy := &model.MealyMachine{
		States:      make(map[string]bool),
		Transitions: make(map[string]map[string]model.MealyTransition),
	}

	for srcState, transitions := range moore.Transitions {
		mealy.States[srcState] = true
		if _, ok := mealy.Transitions[srcState]; !ok {
			mealy.Transitions[srcState] = make(map[string]model.MealyTransition)
		}
		for input, dstStateName := range transitions {
			mealy.States[dstStateName] = true
			output := moore.States[dstStateName].Output
			mealy.Transitions[srcState][input] = model.MealyTransition{
				DestinationState: dstStateName,
				Output:           output,
			}
		}
	}
	return mealy
}

func MealyToMoore(mealy *model.MealyMachine) *model.MooreMachine {
	moore := &model.MooreMachine{
		States:      make(map[string]model.MooreState),
		Transitions: make(map[string]map[string]string),
	}

	stateMap := make(map[string]map[string]string)
	for mealyStateName := range mealy.States {
		stateMap[mealyStateName] = make(map[string]string)
	}

	outputsForState := make(map[string]map[string]bool)
	for mealyStateName := range mealy.States {
		outputsForState[mealyStateName] = make(map[string]bool)
	}
	for _, transitions := range mealy.Transitions {
		for _, transition := range transitions {
			outputsForState[transition.DestinationState][transition.Output] = true
		}
	}

	for mealyStateName, outputs := range outputsForState {
		if len(outputs) == 0 {
			mooreStateName := mealyStateName
			moore.States[mooreStateName] = model.MooreState{Name: mooreStateName, Output: "λ"}
			stateMap[mealyStateName]["λ"] = mooreStateName
		} else {
			for output := range outputs {
				mooreStateName := fmt.Sprintf("%s%s", mealyStateName, output)
				moore.States[mooreStateName] = model.MooreState{Name: mooreStateName, Output: output}
				stateMap[mealyStateName][output] = mooreStateName
			}
		}
	}

	for srcMealyState, transitions := range mealy.Transitions {
		for input, transition := range transitions {
			dstMealyState := transition.DestinationState
			output := transition.Output
			dstMooreState := stateMap[dstMealyState][output]

			if len(stateMap[srcMealyState]) == 0 {
				srcMooreState := srcMealyState
				if _, ok := moore.Transitions[srcMooreState]; !ok {
					moore.Transitions[srcMooreState] = make(map[string]string)
				}
				moore.Transitions[srcMooreState][input] = dstMooreState
			} else {
				for _, srcMooreState := range stateMap[srcMealyState] {
					if _, ok := moore.Transitions[srcMooreState]; !ok {
						moore.Transitions[srcMooreState] = make(map[string]string)
					}
					moore.Transitions[srcMooreState][input] = dstMooreState
				}
			}
		}
	}
	return moore
}
