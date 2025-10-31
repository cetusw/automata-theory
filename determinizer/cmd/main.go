package main

import (
	"determinizer/pkg/determinizer"
	"determinizer/pkg/model"
	"determinizer/pkg/parser"
	"determinizer/pkg/writer"
	"flag"
	"fmt"
	"os"
)

func main() {
	inputFile, outputFile := parseInput()
	assertInput(inputFile, outputFile)

	originalNFA, err := parseNFAFromFile(*inputFile)
	if err != nil {
		fmt.Printf("Error parsing input file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully parsed original NFA with %d states.\n", len(originalNFA.States))

	m := determinizer.NewDeterminizer(originalNFA)
	newDFA := m.Run()

	w := writer.NewWriter()
	err = w.WriteToFile(newDFA, *outputFile)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully wrote minimized DFA to %s\n", *outputFile)
}

func assertInput(inputFile *string, outputFile *string) {
	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: go run . -in <input.dot> -out <output.dot>")
		os.Exit(1)
	}
}

func parseNFAFromFile(filePath string) (*model.NFA, error) {
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
