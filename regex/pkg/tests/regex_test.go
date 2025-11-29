package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"regex/pkg/determinizer"
	"regex/pkg/minimizer"
	"regex/pkg/postfix"
	"regex/pkg/regex"
	"regex/pkg/writer"
)

func runTest(t *testing.T, regexInput string, expectedOutput string) {
	postfixStr, err := postfix.ToPostfix(regexInput)
	if err != nil {
		t.Fatalf("Postfix conversion failed: %v", err)
	}

	regConv := regex.NewConverter()
	nfa, err := regConv.ConvertToNFA(postfixStr)
	if err != nil {
		t.Fatalf("NFA conversion failed: %v", err)
	}

	det := determinizer.NewDeterminizer(nfa)
	dfa := det.Run()

	mini := minimizer.NewMinimizer(dfa)
	minimizedDFA := mini.Minimize()

	tmpFileName := fmt.Sprintf("test_temp_%d.dot", os.Getpid())
	w := writer.NewWriter()
	err = w.WriteToFile(minimizedDFA, tmpFileName)
	if err != nil {
		t.Fatalf("Writer failed: %v", err)
	}
	defer os.Remove(tmpFileName)

	content, err := os.ReadFile(tmpFileName)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	actualOutput := string(content)

	assert.Equal(t, strings.TrimSpace(expectedOutput), strings.TrimSpace(actualOutput))
}

func Test01(t *testing.T) {
	const expectedResult01 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S0;
	node [shape = circle]; S1;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S1 [label = "a"];
	S0 -> S0 [label = "b"];
	S1 -> S0 [label = "a"];
	S1 -> S1 [label = "b"];
}`
	runTest(t, `(ab*a|b)*`, expectedResult01)
}

func Test02(t *testing.T) {
	const expectedResult02 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S0;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S0 [label = "a"];
	S0 -> S0 [label = "b"];
}`
	runTest(t, `(a*|b*)*`, expectedResult02)
}

func Test03(t *testing.T) {
	const expectedResult03 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1;
	node [shape = circle]; S0;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S0 [label = "a"];
	S0 -> S1 [label = "b"];
	S1 -> S0 [label = "a"];
	S1 -> S1 [label = "b"];
}`
	runTest(t, `(a*|b*|b)*b`, expectedResult03)
}

func Test04(t *testing.T) {
	const expectedResult04 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1;
	node [shape = circle]; S0;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S0 [label = "a"];
	S0 -> S1 [label = "b"];
	S0 -> S0 [label = "c"];
	S1 -> S1 [label = "a"];
	S1 -> S1 [label = "b"];
	S1 -> S1 [label = "c"];
}`
	runTest(t, `(a*c*a*)*b(a*b*c*)*`, expectedResult04)
}

func Test05(t *testing.T) {
	const expectedResult05 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1;
	node [shape = circle]; S0;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S0 [label = "r"];
	S0 -> S1 [label = "s"];
	S1 -> S0 [label = "t"];
	S1 -> S1 [label = "u"];
}`
	runTest(t, `(r*|su*t)*su*`, expectedResult05)
}

func Test06(t *testing.T) {
	const expectedResult06 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S4 S5 S6;
	node [shape = circle]; S0 S1 S2 S3 S7;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S1 [label = "a"];
	S1 -> S2 [label = "a"];
	S1 -> S1 [label = "b"];
	S1 -> S3 [label = "d"];
	S2 -> S3 [label = "d"];
	S3 -> S4 [label = "f"];
	S4 -> S6 [label = "a"];
	S4 -> S5 [label = "b"];
	S4 -> S3 [label = "d"];
	S5 -> S2 [label = "a"];
	S5 -> S7 [label = "b"];
	S5 -> S3 [label = "d"];
	S6 -> S6 [label = "a"];
	S6 -> S7 [label = "b"];
	S6 -> S3 [label = "d"];
	S7 -> S7 [label = "b"];
	S7 -> S3 [label = "d"];
}`
	runTest(t, `ab*((a|b*)df(b|a*))((a|b*)df(b|a*))*`, expectedResult06)
}

func Test07(t *testing.T) {
	const expectedResult07 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S0 S1 S2 S3 S6;
	node [shape = circle]; S4 S5 S7;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S7 [label = "a"];
	S0 -> S2 [label = "b"];
	S1 -> S4 [label = "a"];
	S1 -> S3 [label = "b"];
	S2 -> S4 [label = "a"];
	S2 -> S2 [label = "b"];
	S3 -> S5 [label = "a"];
	S3 -> S2 [label = "b"];
	S4 -> S2 [label = "a"];
	S4 -> S4 [label = "b"];
	S5 -> S2 [label = "a"];
	S5 -> S6 [label = "b"];
	S6 -> S1 [label = "a"];
	S6 -> S4 [label = "b"];
	S7 -> S2 [label = "a"];
	S7 -> S5 [label = "b"];
}`
	runTest(t, `(ab*a|b)(ab*a|b)*|abb(ab)*|ε`, expectedResult07)
}

func Test08(t *testing.T) {
	const expectedResult08 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1 S10 S3 S4 S5 S6 S8 S9;
	node [shape = circle]; S0 S2 S7;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S9 [label = "c"];
	S1 -> S3 [label = "a"];
	S1 -> S10 [label = "b"];
	S10 -> S10 [label = "b"];
	S2 -> S3 [label = "a"];
	S3 -> S2 [label = "b"];
	S4 -> S2 [label = "b"];
	S4 -> S5 [label = "c"];
	S5 -> S7 [label = "a"];
	S5 -> S1 [label = "b"];
	S5 -> S6 [label = "c"];
	S6 -> S2 [label = "b"];
	S6 -> S6 [label = "c"];
	S7 -> S8 [label = "c"];
	S8 -> S7 [label = "a"];
	S8 -> S10 [label = "b"];
	S9 -> S4 [label = "a"];
	S9 -> S10 [label = "b"];
}`
	runTest(t, `cac*(ba)*|(ca)*cb*`, expectedResult08)
}

func Test09(t *testing.T) {
	const expectedResult09 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S2 S3;
	node [shape = circle]; S0 S1 S4;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S4 [label = "a"];
	S1 -> S1 [label = "a"];
	S1 -> S3 [label = "b"];
	S2 -> S1 [label = "a"];
	S2 -> S2 [label = "b"];
	S4 -> S1 [label = "a"];
	S4 -> S2 [label = "b"];
}`
	runTest(t, `ab*b*a*b`, expectedResult09)
}

func Test10(t *testing.T) {
	const expectedResult10 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1;
	node [shape = circle]; S0;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S1 [label = "a"];
	S0 -> S1 [label = "b"];
	S0 -> S1 [label = "c"];
	S1 -> S1 [label = "a"];
	S1 -> S1 [label = "b"];
	S1 -> S1 [label = "c"];
}`
	runTest(t, `(a*cc*b*|b*aa*c*|a*bb*c*)(a*cc*b*|b*aa*c*|a*bb*c*)*`, expectedResult10)
}

func Test11(t *testing.T) {
	const expectedResult11 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1 S2 S3 S4 S5 S6 S7;
	node [shape = circle]; S0;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S1 [label = "a"];
	S0 -> S4 [label = "b"];
	S0 -> S3 [label = "c"];
	S1 -> S1 [label = "a"];
	S1 -> S7 [label = "b"];
	S1 -> S3 [label = "c"];
	S2 -> S2 [label = "b"];
	S3 -> S2 [label = "b"];
	S3 -> S3 [label = "c"];
	S4 -> S5 [label = "a"];
	S4 -> S4 [label = "b"];
	S4 -> S6 [label = "c"];
	S5 -> S5 [label = "a"];
	S5 -> S6 [label = "c"];
	S6 -> S6 [label = "c"];
	S7 -> S7 [label = "b"];
	S7 -> S6 [label = "c"];
}`
	runTest(t, `a*cc*b*|b*aa*c*|a*bb*c*`, expectedResult11)
}

func Test12(t *testing.T) {
	const expectedResult12 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S1 S3 S4;
	node [shape = circle]; S0 S2;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S1 [label = "a"];
	S0 -> S3 [label = "b"];
	S0 -> S2 [label = "c"];
	S1 -> S1 [label = "a"];
	S1 -> S3 [label = "b"];
	S1 -> S2 [label = "c"];
	S2 -> S2 [label = "a"];
	S2 -> S4 [label = "b"];
	S2 -> S2 [label = "c"];
	S3 -> S3 [label = "a"];
	S3 -> S3 [label = "b"];
}`
	runTest(t, `b|(a*c*a*)*b|aa*|((a|b)*bb*|aa*)((a|b)*bb*|aa*)*`, expectedResult12)
}

func Test13(t *testing.T) {
	const expectedResult13 = `digraph FiniteStateMachine {
	rankdir=LR;
	node [shape = doublecircle]; S0 S1 S2 S3 S6 S7 S8;
	node [shape = circle]; S4 S5 S9;
	start [shape=point, style=invis];
	start -> S0;
	S0 -> S1 [label = "a"];
	S0 -> S3 [label = "b"];
	S0 -> S6 [label = "c"];
	S1 -> S1 [label = "a"];
	S1 -> S4 [label = "b"];
	S1 -> S8 [label = "c"];
	S2 -> S2 [label = "a"];
	S2 -> S4 [label = "b"];
	S3 -> S2 [label = "a"];
	S3 -> S3 [label = "b"];
	S3 -> S5 [label = "c"];
	S4 -> S2 [label = "a"];
	S4 -> S4 [label = "b"];
	S5 -> S7 [label = "b"];
	S5 -> S5 [label = "c"];
	S6 -> S9 [label = "a"];
	S6 -> S7 [label = "b"];
	S6 -> S6 [label = "c"];
	S7 -> S7 [label = "b"];
	S7 -> S5 [label = "c"];
	S8 -> S9 [label = "a"];
	S8 -> S8 [label = "c"];
	S9 -> S9 [label = "a"];
	S9 -> S8 [label = "c"];
}`
	runTest(t, `a*(a|b)*a|b*|(c|b)*b|c*(c|a)*c`, expectedResult13)
}
