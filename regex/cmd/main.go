package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regex/pkg/minimizer"
	"regex/pkg/postfix"
	"regex/pkg/regex"
	"strings"

	"regex/pkg/determinizer"
	"regex/pkg/writer"
)

type config struct {
	input  *string
	output *string
}

func main() {
	c := parseCliFlags()
	assertInput(c.input, c.output)

	data, err := os.ReadFile(*c.input)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		os.Exit(1)
	}

	inputString := strings.TrimSpace(string(data))

	inputPostfix, err := postfix.ToPostfix(inputString)
	log.Println(inputPostfix)
	if err != nil {
		fmt.Printf("Failed to convert input to postfix notation: %v\n", err)
		os.Exit(1)
	}

	regexConverter := regex.NewConverter()
	nfa, err := regexConverter.ConvertToNFA(inputPostfix)
	if err != nil {
		fmt.Printf("Failed to convert regex to NFA: %v\n", err)
		os.Exit(1)
	}

	d := determinizer.NewDeterminizer(nfa)
	dfa := d.Run()

	m := minimizer.NewMinimizer(dfa)
	minimizedDFA := m.Minimize()

	w := writer.NewWriter()
	err = w.WriteToFile(minimizedDFA, *c.output)
	if err != nil {
		fmt.Printf("Failed to write to output file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully converted Regular Expression to minimized DFA to %s\n", *c.output)
}

func assertInput(inputFile *string, outputFile *string) {
	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Использование: go run . -in <input_file> -out <output_file>")
		os.Exit(1)
	}
}

func parseCliFlags() *config {
	inputFile := flag.String("in", "", "Входной файл")
	outputFile := flag.String("out", "", "Выходной файл")
	flag.Parse()

	return &config{
		input:  inputFile,
		output: outputFile,
	}
}
