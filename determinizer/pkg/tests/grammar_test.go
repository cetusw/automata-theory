package tests

import (
	"os"
	"path/filepath"
	"testing"

	"determinizer/pkg/determinizer"
	"determinizer/pkg/parser"
)

const (
	testGrammarDataDir = "data/grammardata"
)

func TestGrammarToDFAPipeline(t *testing.T) {
	testCases := []struct {
		name         string
		inputFile    string
		expectedFile string
	}{
		{"Right-Linear Grammar", "right_linear_grammar.txt", "right_linear_dfa_expected.dot"},
		{"Left-Linear Grammar", "left_linear_grammar.txt", "left_linear_dfa_expected.dot"},
		{"Epsilon Grammar", "epsilon_grammar.txt", "epsilon_dfa_expected.dot"},
	}

	if err := os.MkdirAll(TestOutputDir, 0755); err != nil {
		t.Fatalf("Не удалось создать директорию %s: %v", TestOutputDir, err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grammarPath := filepath.Join(testGrammarDataDir, tc.inputFile)
			grammarBytes, err := os.ReadFile(grammarPath)
			if err != nil {
				t.Fatalf("Не удалось прочитать файл грамматики %s: %v", grammarPath, err)
			}

			nfa, err := parser.ParseGrammarToNFA(string(grammarBytes))
			if err != nil {
				t.Fatalf("Не удалось распарсить грамматику в НКА: %v", err)
			}

			d := determinizer.NewDeterminizer(nfa)
			actualDFA := d.Run()

			writeDFAOutput(t, actualDFA, tc.inputFile)

			expectedDFA := parseDFAFile(t, filepath.Join(testGrammarDataDir, tc.expectedFile))

			assertDFAEqual(t, expectedDFA, actualDFA)
		})
	}
}
