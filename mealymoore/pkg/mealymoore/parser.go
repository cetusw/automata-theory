package mealymoore

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mealymoore/pkg/mealymoore/model"
)

func extractLabelContent(attr string) (string, error) {
	start := strings.Index(attr, "\"")
	end := strings.LastIndex(attr, "\"")
	if start == -1 || end == -1 || start == end {
		return "", fmt.Errorf("неверный формат атрибута метки: %s", attr)
	}
	return attr[start+1 : end], nil
}

func ParseMooreMachine(filename string) (*model.MooreMachine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	machine := &model.MooreMachine{
		States:      make(map[string]model.MooreState),
		Transitions: make(map[string]map[string]string),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "digraph") || strings.HasPrefix(line, "}") {
			continue
		}

		if strings.Contains(line, "->") {
			parts := strings.Split(line, "->")
			src := strings.TrimSpace(parts[0])
			destAndLabel := strings.Split(parts[1], "[")
			dst := strings.TrimSpace(destAndLabel[0])
			label, err := extractLabelContent(destAndLabel[1])
			if err != nil {
				return nil, err
			}
			input := strings.Trim(label, `";`)
			if _, ok := machine.Transitions[src]; !ok {
				machine.Transitions[src] = make(map[string]string)
			}
			machine.Transitions[src][input] = dst
		} else if strings.Contains(line, "[label=") {
			parts := strings.Split(line, "[")
			name := strings.TrimSpace(parts[0])
			label, err := extractLabelContent(parts[1])
			if err != nil {
				return nil, err
			}

			labelParts := strings.Split(label, "/")
			if len(labelParts) != 2 {
				return nil, fmt.Errorf("неверный формат метки для состояния Мура (ожидалось 'Имя/Выход'), получено: '%s'", label)
			}
			output := strings.TrimSpace(labelParts[1])
			machine.States[name] = model.MooreState{Name: name, Output: output}
		}
	}
	return machine, scanner.Err()
}

func ParseMealyMachine(filename string) (*model.MealyMachine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	machine := &model.MealyMachine{
		States:      make(map[string]bool),
		Transitions: make(map[string]map[string]model.MealyTransition),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "digraph") || strings.HasPrefix(line, "}") {
			continue
		}

		if strings.Contains(line, "->") {
			parts := strings.Split(line, "->")
			src := strings.TrimSpace(parts[0])
			destAndLabel := strings.Split(parts[1], "[")
			dst := strings.TrimSpace(destAndLabel[0])
			label, err := extractLabelContent(destAndLabel[1])
			if err != nil {
				return nil, err
			}
			labelParts := strings.Split(label, "/")
			if len(labelParts) != 2 {
				return nil, fmt.Errorf("неверный формат метки для перехода Мили: %s", label)
			}
			input := strings.TrimSpace(labelParts[0])
			output := strings.TrimSpace(labelParts[1])
			if _, ok := machine.Transitions[src]; !ok {
				machine.Transitions[src] = make(map[string]model.MealyTransition)
			}
			machine.Transitions[src][input] = model.MealyTransition{DestinationState: dst, Output: output}
			machine.States[src] = true
			machine.States[dst] = true
		} else if strings.Contains(line, "[label=") {
			parts := strings.Split(line, "[")
			name := strings.TrimSpace(parts[0])
			machine.States[name] = true
		}
	}
	return machine, scanner.Err()
}
