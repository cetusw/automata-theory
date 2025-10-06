package minimizer

import (
	"fmt"
	"sort"
	"strings"

	"minimisation/pkg/model"
)

const (
	deadTransition = "X-"
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
	m.removeUnreachableStates()
	m.initializePartitions()
	m.refinePartitions()
	return m.buildMinimizedDFA()
}

func (m *Minimizer) initializePartitions() {
	m.partitions = make(map[string]int)
	for _, state := range m.dfa.States {
		if m.dfa.AcceptingStates[state] {
			m.partitions[state] = 1
		} else {
			m.partitions[state] = 0
		}
	}
	m.numPartitions = 2
}

func (m *Minimizer) refinePartitions() {
	for {
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

	for _, state := range m.dfa.States {
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
		destState := m.dfa.Transitions[state][symbol]
		if destState == "" {
			signature.WriteString(deadTransition)
			continue
		}
		signature.WriteString(fmt.Sprintf("%d-", m.partitions[destState]))
	}
	return signature.String()
}

func (m *Minimizer) buildMinimizedDFA() *model.DFA {
	minDFA := model.NewDFA()
	minDFA.Alphabet = m.dfa.Alphabet

	partitionNames := m.createPartitionNames()
	stateMap := m.buildStateMap(partitionNames)

	m.populateMinimizedDFA(minDFA, stateMap)

	sort.Strings(minDFA.States)
	return minDFA
}

func (m *Minimizer) createPartitionNames() map[int]string {
	names := make(map[int]string)
	for i := 0; i < m.numPartitions; i++ {
		names[i] = fmt.Sprintf("S%d", i)
	}
	return names
}

func (m *Minimizer) buildStateMap(names map[int]string) map[string]string {
	stateMap := make(map[string]string)
	for oldState, partID := range m.partitions {
		stateMap[oldState] = names[partID]
	}
	return stateMap
}

func (m *Minimizer) populateMinimizedDFA(
	minDFA *model.DFA,
	stateMap map[string]string,
) {
	processedNewStates := make(map[string]bool)
	for oldState, newState := range stateMap {
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
	for _, symbol := range m.dfa.Alphabet {
		oldDest := m.dfa.Transitions[oldState][symbol]
		if oldDest != "" {
			minDFA.Transitions[newState][symbol] = stateMap[oldDest]
		}
	}
}

func (m *Minimizer) removeUnreachableStates() {
	reachable := make(map[string]bool)
	queue := []string{m.dfa.StartState}
	reachable[m.dfa.StartState] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, symbol := range m.dfa.Alphabet {
			if dest, ok := m.dfa.Transitions[current][symbol]; ok && !reachable[dest] {
				reachable[dest] = true
				queue = append(queue, dest)
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
