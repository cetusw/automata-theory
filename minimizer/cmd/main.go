package main

import (
	"flag"
	"fmt"
	"os"

	"minimisation/pkg/minimizer"
	"minimisation/pkg/model"
	"minimisation/pkg/parser"
	"minimisation/pkg/writer"
)

func main() {
	inputFile, outputFile := parseInput()
	assertInput(inputFile, outputFile)

	originalDFA, err := parseDFAFromFile(*inputFile)
	if err != nil {
		fmt.Printf("Error parsing input file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully parsed original DFA with %d states.\n", len(originalDFA.States))

	m := minimizer.NewMinimizer(originalDFA)
	minimizedDFA := m.Minimize()
	fmt.Printf("Minimized DFA has %d states.\n", len(minimizedDFA.States))

	w := writer.NewWriter()
	err = w.WriteToFile(minimizedDFA, *outputFile)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully wrote minimized DFA to %s\n", *outputFile)
}

func parseDFAFromFile(filePath string) (*model.DFA, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error parsing input file: %v\n", err)
		os.Exit(1)
	}
	p := parser.NewParser(string(data))
	return p.Parse()
}

func parseInput() (*string, *string) {
	inputFile := flag.String("in", "", "Input file in .dot format")
	outputFile := flag.String("out", "", "Output file for the minimized DFA")
	flag.Parse()

	return inputFile, outputFile
}

func assertInput(inputFile *string, outputFile *string) {
	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: go run . -in <input.dot> -out <output.dot>")
		os.Exit(1)
	}
}
