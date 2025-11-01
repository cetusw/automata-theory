package tests

import (
	"os"
	"path/filepath"
	"testing"

	"determinizer/pkg/determinizer"
	"determinizer/pkg/model"
	"determinizer/pkg/parser"
)

const (
	testNFADataDir = "data/nfadata"
)

func parseNFAFile(t *testing.T, filePath string) *model.NFA {
	t.Helper()
	dotBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Не удалось прочитать файл %s: %v", filePath, err)
	}

	nfa, err := parser.ParseNFA(string(dotBytes))
	if err != nil {
		t.Fatalf("Не удалось распарсить НКА из %s: %v", filePath, err)
	}
	return nfa
}

func TestDeterminize(t *testing.T) {
	testCases := []struct {
		name         string
		inputFile    string
		expectedFile string
	}{
		{"Simple NFA", "simple_nfa.dot", "simple_dfa_expected.dot"},
		{"Hard NFA", "hard_nfa.dot", "hard_dfa_expected.dot"},
		{"Merge Paths", "merge_paths_nfa.dot", "merge_paths_dfa_expected.dot"},
		{"Overlapping Paths", "overlapping_paths_nfa.dot", "overlapping_paths_dfa_expected.dot"},
	}

	if err := os.MkdirAll(TestOutputDir, 0755); err != nil {
		t.Fatalf("Не удалось создать директорию %s: %v", TestOutputDir, err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nfa := parseNFAFile(t, filepath.Join(testNFADataDir, tc.inputFile))

			d := determinizer.NewDeterminizer(nfa)
			actualDFA := d.Run()

			writeDFAOutput(t, actualDFA, tc.inputFile)

			expectedDFA := parseDFAFile(t, filepath.Join(testNFADataDir, tc.expectedFile))

			assertDFAEqual(t, expectedDFA, actualDFA)
		})
	}
}
