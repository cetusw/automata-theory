package mealymoore

import (
	"fmt"
	"os"
	"strings"

	"mealymoore/pkg/mealymoore/model"
)

func WriteMooreMachine(machine *model.MooreMachine, filename string) error {
	var builder strings.Builder
	builder.WriteString("digraph MooreMachine {\n")

	for name, state := range machine.States {
		builder.WriteString(fmt.Sprintf("  %s [label=\"%s/%s\"];\n", name, state.Name, state.Output))
	}

	for src, transitions := range machine.Transitions {
		for input, dst := range transitions {
			builder.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n", src, dst, input))
		}
	}

	builder.WriteString("}\n")
	return os.WriteFile(filename, []byte(builder.String()), 0644)
}

func WriteMealyMachine(machine *model.MealyMachine, filename string) error {
	var builder strings.Builder
	builder.WriteString("digraph MealyMachine {\n")

	for state := range machine.States {
		builder.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", state, state))
	}

	for src, transitions := range machine.Transitions {
		for input, transition := range transitions {
			builder.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s/%s\"];\n", src, transition.DestinationState, input, transition.Output))
		}
	}

	builder.WriteString("}\n")
	return os.WriteFile(filename, []byte(builder.String()), 0644)
}
