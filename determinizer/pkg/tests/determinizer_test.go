package tests

import (
	"determinizer/pkg/determinizer"
	"determinizer/pkg/writer"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"determinizer/pkg/model"
	"determinizer/pkg/parser"
)

const (
	testDataDir   = "../tests/data"
	testOutputDir = "../tests/output"
)

var (
	testCases = []struct {
		name         string
		inputFile    string
		expectedFile string
	}{
		{
			name:         "Simple NFA",
			inputFile:    "simple_nfa.dot",
			expectedFile: "simple_dfa_expected.dot",
		},
		{
			name:         "Hard NFA",
			inputFile:    "hard_nfa.dot",
			expectedFile: "hard_dfa_expected.dot",
		},
		{
			name:         "NFA with merging epsilon paths",
			inputFile:    "merge_paths_dfa.dot",
			expectedFile: "merge_paths_dfa_expected.dot",
		},
		{
			name:         "NFA with overlapping paths to accepting state",
			inputFile:    "overlapping_paths_nfa.dot",
			expectedFile: "overlapping_paths_nfa_expected.dot",
		},
	}
)

func TestDeterminize(t *testing.T) {
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		t.Fatalf("Could not create output directory %s: %v", testOutputDir, err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nfa := parseNFAFile(t, tc.inputFile)

			d := determinizer.NewDeterminizer(nfa)
			actualDFA := d.Run()

			resultFile := strings.Split(tc.inputFile, ".")[0] + "_result.dot"
			outputFilename := filepath.Join(testOutputDir, resultFile)
			w := writer.NewWriter()
			if err := w.WriteToFile(actualDFA, outputFilename); err != nil {
				t.Fatalf("Failed to write actual DFA to %s: %v", outputFilename, err)
			}

			expectedDFA := parseDFAFile(t, tc.expectedFile)

			assertDFAEqual(t, expectedDFA, actualDFA)
		})
	}
}

func parseNFAFile(t *testing.T, filename string) *model.NFA {
	t.Helper()
	path := filepath.Join(testDataDir, filename)
	dotBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	nfa, err := parser.ParseNFA(string(dotBytes))
	if err != nil {
		t.Fatalf("Failed to parse NFA from %s: %v", filename, err)
	}
	return nfa
}

func parseDFAFile(t *testing.T, filename string) *model.DFA {
	t.Helper()
	path := filepath.Join(testDataDir, filename)
	dotBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	dfa, err := parser.ParseDFA(string(dotBytes))
	if err != nil {
		t.Fatalf("Failed to parse expected DFA from %s: %v", filename, err)
	}
	return dfa
}

func assertDFAEqual(t *testing.T, expected, actual *model.DFA) {
	t.Helper()

	if expected.StartState != actual.StartState {
		t.Errorf("StartState mismatch: want %q, got %q", expected.StartState, actual.StartState)
	}

	if !reflect.DeepEqual(expected.States, actual.States) {
		t.Errorf("States mismatch: \nwant %v, \ngot  %v", expected.States, actual.States)
	}

	if !reflect.DeepEqual(expected.Alphabet, actual.Alphabet) {
		t.Errorf("Alphabet mismatch: \nwant %v, \ngot  %v", expected.Alphabet, actual.Alphabet)
	}

	if !reflect.DeepEqual(expected.AcceptingStates, actual.AcceptingStates) {
		t.Errorf("AcceptingStates mismatch: \nwant %v, \ngot  %v", expected.AcceptingStates, actual.AcceptingStates)
	}

	if !reflect.DeepEqual(expected.Transitions, actual.Transitions) {
		t.Errorf("Transitions mismatch: \nwant %+v, \ngot  %+v", expected.Transitions, actual.Transitions)
	}
}
