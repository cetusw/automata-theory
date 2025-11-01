package main

import (
	"flag"
	"fmt"
	"os"

	"determinizer/pkg/determinizer"
	"determinizer/pkg/model"
	"determinizer/pkg/parser"
	"determinizer/pkg/writer"
)

const (
	nfaType     = "nfa"
	grammarType = "grammar"
)

type config struct {
	input  *string
	output *string
	t      *string
}

func main() {
	c := parseCliFlags()
	assertInput(c.input, c.output)

	var originalNFA *model.NFA
	var err error

	data, err := os.ReadFile(*c.input)
	if err != nil {
		fmt.Printf("Ошибка чтения входного файла: %v\n", err)
		os.Exit(1)
	}
	inputString := string(data)

	switch *c.t {
	case nfaType:
		originalNFA, err = parser.ParseNFA(inputString)
		fmt.Println("Парсинг входного файла как NFA (.dot)...")
	case grammarType:
		originalNFA, err = parser.ParseGrammarToNFA(inputString)
		fmt.Println("Парсинг входного файла как грамматики...")
	default:
		fmt.Printf("Неизвестный тип входных данных: %s. Используйте 'nfa' или 'grammar'.\n", *c.t)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Ошибка парсинга входного файла: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Успешно построен НКА с %d состояниями.\n", len(originalNFA.States))

	d := determinizer.NewDeterminizer(originalNFA)
	newDFA := d.Run()

	w := writer.NewWriter()
	err = w.WriteToFile(newDFA, *c.output)
	if err != nil {
		fmt.Printf("Ошибка записи выходного файла: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Успешно записан ДКА в %s\n", *c.output)
}

func assertInput(inputFile *string, outputFile *string) {
	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Использование: go run . -in <input_file> -out <output_file> [-type <nfa|grammar>]")
		os.Exit(1)
	}
}

func parseCliFlags() *config {
	inputFile := flag.String("in", "", "Входной файл")
	outputFile := flag.String("out", "", "Выходной файл")
	inputType := flag.String("type", "nfa", "Тип входных данных: 'nfa' (файл .dot) или 'grammar' (файл с грамматикой)")
	flag.Parse()

	return &config{
		input:  inputFile,
		output: outputFile,
		t:      inputType,
	}
}
