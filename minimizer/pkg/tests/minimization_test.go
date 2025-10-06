package tests

import (
	"reflect"
	"sort"
	"testing"

	"minimisation/pkg/minimizer"
	"minimisation/pkg/model"
)

type minimizerTestCase struct {
	name        string
	inputDFA    *model.DFA
	expectedDFA *model.DFA
}

var (
	simpleCaseInput = &model.DFA{
		States:   []string{"q0", "q1", "q2", "q3", "q4"},
		Alphabet: []string{"0", "1"},
		Transitions: map[string]map[string]string{
			"q0": {"0": "q1", "1": "q2"}, "q1": {"0": "q3", "1": "q4"},
			"q2": {"0": "q3", "1": "q4"}, "q3": {"0": "q3", "1": "q3"},
			"q4": {"0": "q4", "1": "q4"},
		},
		StartState: "q0", AcceptingStates: map[string]bool{"q3": true, "q4": true},
	}
	simpleCaseExpected = &model.DFA{
		States:   []string{"S0", "S1", "S2"},
		Alphabet: []string{"0", "1"},
		Transitions: map[string]map[string]string{
			"S0": {"0": "S1", "1": "S1"}, "S1": {"0": "S2", "1": "S2"},
			"S2": {"0": "S2", "1": "S2"},
		},
		StartState: "S0", AcceptingStates: map[string]bool{"S2": true},
	}
	alreadyMinimalDFA = &model.DFA{
		States:   []string{"S0", "S1", "S2"},
		Alphabet: []string{"0", "1"},
		Transitions: map[string]map[string]string{
			"S0": {"0": "S1", "1": "S2"}, "S1": {"0": "S2", "1": "S1"},
			"S2": {"0": "S2", "1": "S2"},
		},
		StartState: "S0", AcceptingStates: map[string]bool{"S1": true},
	}
	unreachableStateInput = &model.DFA{
		States:   []string{"q0", "q1", "q2", "q3", "q4", "q_unreachable"},
		Alphabet: []string{"0", "1"},
		Transitions: map[string]map[string]string{
			"q0": {"0": "q1", "1": "q2"}, "q1": {"0": "q3", "1": "q4"},
			"q2": {"0": "q3", "1": "q4"}, "q3": {"0": "q3", "1": "q3"},
			"q4": {"0": "q4", "1": "q4"}, "q_unreachable": {"0": "q0", "1": "q1"},
		},
		StartState: "q0", AcceptingStates: map[string]bool{"q3": true, "q4": true},
	}

	minimizerTestCases = []minimizerTestCase{
		{
			name:        "Simple Case With Equivalent States",
			inputDFA:    simpleCaseInput,
			expectedDFA: simpleCaseExpected,
		},
		{
			name:        "Already Minimal DFA",
			inputDFA:    alreadyMinimalDFA,
			expectedDFA: alreadyMinimalDFA,
		},
		{
			name:        "DFA With Unreachable State",
			inputDFA:    unreachableStateInput,
			expectedDFA: simpleCaseExpected,
		},
	}
)

func TestMinimizer(t *testing.T) {
	for _, tc := range minimizerTestCases {
		t.Run(tc.name, func(t *testing.T) {
			m := minimizer.NewMinimizer(tc.inputDFA)
			actualDFA := m.Minimize()

			assertDFAEqual(t, actualDFA, tc.expectedDFA)
		})
	}
}

func assertDFAEqual(t *testing.T, actual, expected *model.DFA) {
	t.Helper()

	sort.Strings(actual.States)
	sort.Strings(expected.States)
	sort.Strings(actual.Alphabet)
	sort.Strings(expected.Alphabet)

	if !reflect.DeepEqual(actual.States, expected.States) {
		t.Errorf("States mismatch:\n got: %v\n want: %v", actual.States, expected.States)
	}
	if !reflect.DeepEqual(actual.Alphabet, expected.Alphabet) {
		t.Errorf("Alphabet mismatch:\n got: %v\n want: %v", actual.Alphabet, expected.Alphabet)
	}
	if actual.StartState != expected.StartState {
		t.Errorf("StartState mismatch:\n got: %s\n want: %s", actual.StartState, expected.StartState)
	}
	if !reflect.DeepEqual(actual.AcceptingStates, expected.AcceptingStates) {
		t.Errorf("AcceptingStates mismatch:\n got: %v\n want: %v", actual.AcceptingStates, expected.AcceptingStates)
	}
	if !reflect.DeepEqual(actual.Transitions, expected.Transitions) {
		t.Errorf("Transitions mismatch:\n got: %v\n want: %v", actual.Transitions, expected.Transitions)
	}
}
