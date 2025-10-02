package main

import (
	"flag"
	"fmt"
	"log"
	"mealymoore/pkg/mealymoore"
	"os"
)

func main() {
	inputFile := flag.String("in", "", "Путь к входному .dot файлу")
	outputFile := flag.String("out", "", "Путь к выходному .dot файлу")
	conversionType := flag.String("type", "", "Тип преобразования: 'moore-to-mealy' или 'mealy-to-moore'")

	flag.Parse()

	if *inputFile == "" || *outputFile == "" || *conversionType == "" {
		fmt.Println("Использование:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch *conversionType {
	case "moore-to-mealy":
		mooreMachine, err := mealymoore.ParseMooreMachine(*inputFile)
		if err != nil {
			log.Fatalf("Ошибка разбора автомата Мура: %v", err)
		}
		mealyMachine := mealymoore.MooreToMealy(mooreMachine)
		err = mealymoore.WriteMealyMachine(mealyMachine, *outputFile)
		if err != nil {
			log.Fatalf("Ошибка записи автомата Мили: %v", err)
		}
		fmt.Printf("Успешно преобразован автомат Мура из '%s' в автомат Мили в '%s'\n", *inputFile, *outputFile)

	case "mealy-to-moore":
		mealyMachine, err := mealymoore.ParseMealyMachine(*inputFile)
		if err != nil {
			log.Fatalf("Ошибка разбора автомата Мили: %v", err)
		}
		mooreMachine := mealymoore.MealyToMoore(mealyMachine)
		err = mealymoore.WriteMooreMachine(mooreMachine, *outputFile)
		if err != nil {
			log.Fatalf("Ошибка записи автомата Мура: %v", err)
		}
		fmt.Printf("Успешно преобразован автомат Мили из '%s' в автомат Мура в '%s'\n", *inputFile, *outputFile)

	default:
		log.Fatalf("Неверный тип преобразования: %s. Используйте 'moore-to-mealy' или 'mealy-to-moore'", *conversionType)
	}
}
