package tests

import (
	"fmt"
	"mealymoore/pkg/mealymoore"
	"mealymoore/pkg/mealymoore/model"
	"reflect"
	"testing"
)

func compareMealyMachines(got, want *model.MealyMachine) (bool, string) {
	if !reflect.DeepEqual(got.States, want.States) {
		return false, fmt.Sprintf("наборы состояний не совпадают.\nПолучено: %v\nОжидалось: %v", got.States, want.States)
	}
	if !reflect.DeepEqual(got.Transitions, want.Transitions) {
		return false, fmt.Sprintf("переходы не совпадают.\nПолучено: %v\nОжидалось: %v", got.Transitions, want.Transitions)
	}
	return true, ""
}

func compareMooreMachines(got, want *model.MooreMachine) (bool, string) {
	if !reflect.DeepEqual(got.States, want.States) {
		return false, fmt.Sprintf("состояния не совпадают.\nПолучено: %v\nОжидалось: %v", got.States, want.States)
	}
	if !reflect.DeepEqual(got.Transitions, want.Transitions) {
		return false, fmt.Sprintf("переходы не совпадают.\nПолучено: %v\nОжидалось: %v", got.Transitions, want.Transitions)
	}
	return true, ""
}

func TestMealyToMoore(t *testing.T) {
	inputFile := "data/mealy.dot"
	mealy, err := mealymoore.ParseMealyMachine(inputFile)
	if err != nil {
		t.Fatalf("Не удалось разобрать автомат Мили из %s: %v", inputFile, err)
	}

	resultMoore := mealymoore.MealyToMoore(mealy)

	expectedMoore := &model.MooreMachine{
		States: map[string]model.MooreState{
			"S0":   {Name: "S0", Output: "λ"},
			"S1y2": {Name: "S1y2", Output: "y2"},
			"S2y1": {Name: "S2y1", Output: "y1"},
			"S2y2": {Name: "S2y2", Output: "y2"},
			"S2y3": {Name: "S2y3", Output: "y3"},
			"S3y1": {Name: "S3y1", Output: "y1"},
			"S3y2": {Name: "S3y2", Output: "y2"},
		},
		Transitions: map[string]map[string]string{
			"S0":   {"x1": "S2y1", "x2": "S3y1"},
			"S1y2": {"x1": "S2y2", "x2": "S3y2"},
			"S2y1": {"x1": "S1y2", "x2": "S3y1"},
			"S2y2": {"x1": "S1y2", "x2": "S3y1"},
			"S2y3": {"x1": "S1y2", "x2": "S3y1"},
			"S3y1": {"x1": "S2y3", "x2": "S3y2"},
			"S3y2": {"x1": "S2y3", "x2": "S3y2"},
		},
	}

	if ok, reason := compareMooreMachines(resultMoore, expectedMoore); !ok {
		t.Errorf("Тест TestMealyToMoore_Final провален: %s", reason)
	}
}

func TestMooreToMealy_Final(t *testing.T) {
	inputFile := "data/moore.dot"
	moore, err := mealymoore.ParseMooreMachine(inputFile)
	if err != nil {
		t.Fatalf("Не удалось разобрать автомат Мура из %s: %v", inputFile, err)
	}

	resultMealy := mealymoore.MooreToMealy(moore)

	expectedMealy := &model.MealyMachine{
		States: map[string]bool{
			"S0_": true, "S1y2": true, "S2y1": true, "S2y2": true,
			"S2y3": true, "S3y1": true, "S3y2": true,
		},
		Transitions: map[string]map[string]model.MealyTransition{
			"S0_":  {"x1": {DestinationState: "S2y1", Output: "y1"}, "x2": {DestinationState: "S3y1", Output: "y1"}},
			"S1y2": {"x1": {DestinationState: "S2y2", Output: "y2"}, "x2": {DestinationState: "S3y2", Output: "y2"}},
			"S2y1": {"x1": {DestinationState: "S1y2", Output: "y2"}, "x2": {DestinationState: "S3y1", Output: "y1"}},
			"S2y2": {"x1": {DestinationState: "S1y2", Output: "y2"}, "x2": {DestinationState: "S3y1", Output: "y1"}},
			"S2y3": {"x1": {DestinationState: "S1y2", Output: "y2"}, "x2": {DestinationState: "S3y1", Output: "y1"}},
			"S3y1": {"x1": {DestinationState: "S2y3", Output: "y3"}, "x2": {DestinationState: "S3y2", Output: "y2"}},
			"S3y2": {"x1": {DestinationState: "S2y3", Output: "y3"}, "x2": {DestinationState: "S3y2", Output: "y2"}},
		},
	}

	if ok, reason := compareMealyMachines(resultMealy, expectedMealy); !ok {
		t.Errorf("Тест TestMooreToMealy_Final провален: %s", reason)
	}
}
