package writer

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"regex/pkg/model"
)

const (
	digraphHeader        = "digraph FiniteStateMachine {\n"
	digraphFooter        = "}\n"
	digraphDirection     = "\trankdir=LR;\n"
	acceptingNodeDecl    = "\tnode [shape = doublecircle]; %s;\n"
	nonAcceptingNodeDecl = "\tnode [shape = circle]; %s;\n"
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
	w.writeNodes(dfa)
	w.writeStartState(dfa)
	w.writeTransitions(dfa)
	w.writeFooter()
	return w.builder.String()
}

func (w *Writer) writeHeader() {
	w.builder.WriteString(digraphHeader)
	w.builder.WriteString(digraphDirection)
}

func (w *Writer) writeNodes(dfa *model.DFA) {
	var accepting, nonAccepting []string

	for _, state := range dfa.States {
		if dfa.AcceptingStates[state] {
			accepting = append(accepting, state)
		} else {
			nonAccepting = append(nonAccepting, state)
		}
	}

	sort.Strings(accepting)
	sort.Strings(nonAccepting)

	if len(accepting) > 0 {
		line := fmt.Sprintf(acceptingNodeDecl, strings.Join(accepting, " "))
		w.builder.WriteString(line)
	}

	if len(nonAccepting) > 0 {
		line := fmt.Sprintf(nonAcceptingNodeDecl, strings.Join(nonAccepting, " "))
		w.builder.WriteString(line)
	}
}

func (w *Writer) writeStartState(dfa *model.DFA) {
	if dfa.StartState == "" {
		return
	}
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
		sortedSymbols := make([]string, 0, len(transitions))
		for symbol := range transitions {
			sortedSymbols = append(sortedSymbols, symbol)
		}
		sort.Strings(sortedSymbols)

		for _, symbol := range sortedSymbols {
			to := transitions[symbol]
			line := fmt.Sprintf(transition, from, to, symbol)
			w.builder.WriteString(line)
		}
	}
}

func (w *Writer) writeFooter() {
	w.builder.WriteString(digraphFooter)
}
