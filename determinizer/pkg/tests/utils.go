package tests

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"determinizer/pkg/model"
	"determinizer/pkg/parser"
	"determinizer/pkg/writer"
)

const (
	TestOutputDir = "output"
)

func parseDFAFile(t *testing.T, filePath string) *model.DFA {
	t.Helper()
	dotBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Не удалось прочитать файл %s: %v", filePath, err)
	}

	dfa, err := parser.ParseDFA(string(dotBytes))
	if err != nil {
		t.Fatalf("Не удалось распарсить ожидаемый ДКА из %s: %v", filePath, err)
	}
	return dfa
}

func assertDFAEqual(t *testing.T, expected, actual *model.DFA) {
	t.Helper()

	if expected.StartState != actual.StartState {
		t.Errorf("Несовпадение StartState: ожидалось %q, получено %q", expected.StartState, actual.StartState)
	}
	if !reflect.DeepEqual(expected.States, actual.States) {
		t.Errorf("Несовпадение States: \nожидалось %v, \nполучено  %v", expected.States, actual.States)
	}
	if !reflect.DeepEqual(expected.Alphabet, actual.Alphabet) {
		t.Errorf("Несовпадение Alphabet: \nожидалось %v, \nполучено  %v", expected.Alphabet, actual.Alphabet)
	}
	if !reflect.DeepEqual(expected.AcceptingStates, actual.AcceptingStates) {
		t.Errorf("Несовпадение AcceptingStates: \nожидалось %v, \nполучено  %v", expected.AcceptingStates, actual.AcceptingStates)
	}
	if !reflect.DeepEqual(expected.Transitions, actual.Transitions) {
		t.Errorf("Несовпадение Transitions: \nожидалось %+v, \nполучено  %+v", expected.Transitions, actual.Transitions)
	}
}

func writeDFAOutput(t *testing.T, dfa *model.DFA, originalInputFile string) {
	t.Helper()
	resultFile := strings.Split(originalInputFile, ".")[0] + "_result.dot"
	outputFilename := filepath.Join(TestOutputDir, resultFile)
	w := writer.NewWriter()
	if err := w.WriteToFile(dfa, outputFilename); err != nil {
		t.Fatalf("Не удалось записать ДКА в %s: %v", outputFilename, err)
	}
}
