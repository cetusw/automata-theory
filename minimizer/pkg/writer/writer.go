package writer

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"minimisation/pkg/model"
)

const (
	digraphHeader        = "digraph FiniteStateMachine {\n"
	digraphFooter        = "}\n"
	digraphDirection     = "\trankdir=LR;\n"
	acceptingStates      = "\tnode [shape = doublecircle]; %s;\n"
	nodes                = "\tnode [shape = circle];\n"
	startState           = "\tstart [shape=point, style=invis];\n"
	startStateTransition = "\tstart -> %s;\n"
	transition           = "\t%s -> %s [label = \"%s\"];\n"
)

type Writer struct {
	builder strings.Builder
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteToFile(dfa *model.DFA, filePath string) error {
	dotString := w.generateDOTString(dfa)
	return os.WriteFile(filePath, []byte(dotString), 0644)
}

func (w *Writer) generateDOTString(dfa *model.DFA) string {
	w.builder.Reset()
	w.writeHeader()
	w.writeAcceptingStates(dfa)
	w.writeAllNodes()
	w.writeStartState(dfa)
	w.writeTransitions(dfa)
	w.writeFooter()
	return w.builder.String()
}

func (w *Writer) writeHeader() {
	w.builder.WriteString(digraphHeader)
	w.builder.WriteString(digraphDirection)
}

func (w *Writer) writeAcceptingStates(dfa *model.DFA) {
	var accepting []string
	for state := range dfa.AcceptingStates {
		accepting = append(accepting, state)
	}
	if len(accepting) > 0 {
		sort.Strings(accepting)
		line := fmt.Sprintf(acceptingStates, strings.Join(accepting, " "))
		w.builder.WriteString(line)
	}
}

func (w *Writer) writeAllNodes() {
	w.builder.WriteString(nodes)
}

func (w *Writer) writeStartState(dfa *model.DFA) {
	w.builder.WriteString(startState)
	line := fmt.Sprintf(startStateTransition, dfa.StartState)
	w.builder.WriteString(line)
}

func (w *Writer) writeTransitions(dfa *model.DFA) {
	sortedStates := make([]string, 0, len(dfa.Transitions))
	for from := range dfa.Transitions {
		sortedStates = append(sortedStates, from)
	}
	sort.Strings(sortedStates)

	for _, from := range sortedStates {
		transitions := dfa.Transitions[from]
		for symbol, to := range transitions {
			line := fmt.Sprintf(transition, from, to, symbol)
			w.builder.WriteString(line)
		}
	}
}

func (w *Writer) writeFooter() {
	w.builder.WriteString(digraphFooter)
}
