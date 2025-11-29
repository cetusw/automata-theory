package minimizer

import (
	"fmt"
	"sort"
	"strings"

	"regex/pkg/model"
)

type Minimizer struct {
	dfa           *model.DFA
	partitions    map[string]int
	numPartitions int
}

func NewMinimizer(dfa *model.DFA) *Minimizer {
	return &Minimizer{dfa: dfa}
}

func (m *Minimizer) Minimize() *model.DFA {
	if len(m.dfa.States) == 0 {
		return m.dfa
	}
	m.removeUnreachableStates()
	if len(m.dfa.States) <= 1 {
		return m.dfa
	}
	m.initializePartitions()
	m.refinePartitions()
	return m.buildMinimizedDFA()
}

func (m *Minimizer) initializePartitions() {
	m.partitions = make(map[string]int)
	hasAccepting := false
	hasNonAccepting := false

	for _, state := range m.dfa.States {
		if m.dfa.AcceptingStates[state] {
			m.partitions[state] = 1
			hasAccepting = true
		} else {
			m.partitions[state] = 0
			hasNonAccepting = true
		}
	}

	if hasAccepting && hasNonAccepting {
		m.numPartitions = 2
	} else {
		m.numPartitions = 1
	}
}

func (m *Minimizer) refinePartitions() {
	for {
		if m.numPartitions >= len(m.dfa.States) {
			break
		}
		newPartitions, newCount := m.splitCurrentPartitions()
		if newCount == m.numPartitions {
			break
		}
		m.partitions = newPartitions
		m.numPartitions = newCount
	}
}

func (m *Minimizer) splitCurrentPartitions() (map[string]int, int) {
	newPartitions := make(map[string]int)
	partitionMap := make(map[string]int)
	nextPartitionID := 0

	sortedStates := make([]string, len(m.dfa.States))
	copy(sortedStates, m.dfa.States)
	sort.Strings(sortedStates)

	for _, state := range sortedStates {
		signature := m.getStateSignature(state)
		groupKey := fmt.Sprintf("%d-%s", m.partitions[state], signature)

		if _, exists := partitionMap[groupKey]; !exists {
			partitionMap[groupKey] = nextPartitionID
			nextPartitionID++
		}
		newPartitions[state] = partitionMap[groupKey]
	}
	return newPartitions, nextPartitionID
}

func (m *Minimizer) getStateSignature(state string) string {
	var signature strings.Builder
	for _, symbol := range m.dfa.Alphabet {
		destState := ""
		if transitions, ok := m.dfa.Transitions[state]; ok {
			destState = transitions[symbol]
		}

		partitionID := -1
		if destState != "" {
			partitionID = m.partitions[destState]
		}
		signature.WriteString(fmt.Sprintf("%d-", partitionID))
	}
	return signature.String()
}

func (m *Minimizer) buildMinimizedDFA() *model.DFA {
	minDFA := model.NewDFA()
	minDFA.Alphabet = m.dfa.Alphabet

	stateMap := m.buildCanonicalStateMap()
	m.populateMinimizedDFA(minDFA, stateMap)

	sort.Strings(minDFA.States)
	return minDFA
}

func (m *Minimizer) buildCanonicalStateMap() map[string]string {
	invertedPartitions := make(map[int][]string)
	for oldState, partID := range m.partitions {
		invertedPartitions[partID] = append(invertedPartitions[partID], oldState)
	}

	canonicalPartitions := make([][]string, 0, len(invertedPartitions))
	for _, states := range invertedPartitions {
		sort.Strings(states)
		canonicalPartitions = append(canonicalPartitions, states)
	}

	sort.Slice(canonicalPartitions, func(i, j int) bool {
		isStartI := false
		for _, s := range canonicalPartitions[i] {
			if s == m.dfa.StartState {
				isStartI = true
				break
			}
		}
		if isStartI {
			return true
		}
		isStartJ := false
		for _, s := range canonicalPartitions[j] {
			if s == m.dfa.StartState {
				isStartJ = true
				break
			}
		}
		if isStartJ {
			return false
		}
		return canonicalPartitions[i][0] < canonicalPartitions[j][0]
	})

	finalStateMap := make(map[string]string)
	for i, partition := range canonicalPartitions {
		newStateName := fmt.Sprintf("S%d", i)
		for _, oldState := range partition {
			finalStateMap[oldState] = newStateName
		}
	}
	return finalStateMap
}

func (m *Minimizer) populateMinimizedDFA(
	minDFA *model.DFA,
	stateMap map[string]string,
) {
	processedNewStates := make(map[string]bool)
	oldStates := make([]string, 0, len(stateMap))
	for oldState := range stateMap {
		oldStates = append(oldStates, oldState)
	}
	sort.Strings(oldStates)

	for _, oldState := range oldStates {
		newState := stateMap[oldState]
		if processedNewStates[newState] {
			continue
		}
		minDFA.States = append(minDFA.States, newState)
		m.addMinimizedTransitions(minDFA, oldState, newState, stateMap)

		if m.dfa.AcceptingStates[oldState] {
			minDFA.AcceptingStates[newState] = true
		}
		processedNewStates[newState] = true
	}
	minDFA.StartState = stateMap[m.dfa.StartState]
}

func (m *Minimizer) addMinimizedTransitions(
	minDFA *model.DFA,
	oldState,
	newState string,
	stateMap map[string]string,
) {
	minDFA.Transitions[newState] = make(map[string]string)
	if oldTransitions, ok := m.dfa.Transitions[oldState]; ok {
		for _, symbol := range m.dfa.Alphabet {
			if oldDest, ok := oldTransitions[symbol]; ok {
				minDFA.Transitions[newState][symbol] = stateMap[oldDest]
			}
		}
	}
}

func (m *Minimizer) removeUnreachableStates() {
	reachable := make(map[string]bool)
	if m.dfa.StartState == "" {
		return
	}
	queue := []string{m.dfa.StartState}
	reachable[m.dfa.StartState] = true
	head := 0
	for head < len(queue) {
		current := queue[head]
		head++
		if transitions, ok := m.dfa.Transitions[current]; ok {
			for _, dest := range transitions {
				if dest != "" && !reachable[dest] {
					reachable[dest] = true
					queue = append(queue, dest)
				}
			}
		}
	}
	var reachableStates []string
	for _, state := range m.dfa.States {
		if reachable[state] {
			reachableStates = append(reachableStates, state)
		} else {
			delete(m.dfa.Transitions, state)
			delete(m.dfa.AcceptingStates, state)
		}
	}
	m.dfa.States = reachableStates
}
