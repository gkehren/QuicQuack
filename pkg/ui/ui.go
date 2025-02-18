package ui

import (
	"QuicQuack/internal/benchmark"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateProtocolSelection State = iota
	StateBenchmarkSelection
	StateOptions
	StateRunning
	StateCompleted
)

type model struct {
	state              State
	cursor             int
	protocols          []string
	selectedProtocols  map[int]struct{}
	benchmarks         []string
	selectedBenchmarks map[int]struct{}
	options            map[string]string
	results            string
}

type benchmarkResultMsg struct {
	results string
}

func initialModel() model {
	return model{
		state:              StateProtocolSelection,
		protocols:          []string{"TCP", "UDP", "QUIC"},
		selectedProtocols:  make(map[int]struct{}),
		benchmarks:         []string{"Latency"},
		selectedBenchmarks: make(map[int]struct{}),
		options:            make(map[string]string),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateProtocolSelection:
			return updateProtocolSelection(m, msg)
		case StateBenchmarkSelection:
			return updateBenchmarkSelection(m, msg)
		case StateOptions:
			return updateOptions(m, msg)
		case StateRunning:
			return updateRunning(m, msg)
		}
	case benchmarkResultMsg:
		m.results = msg.results
		m.state = StateCompleted
		return m, tea.Quit
	}
	return m, nil
}

func updateProtocolSelection(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.protocols)-1 {
			m.cursor++
		}
	case " ":
		_, ok := m.selectedProtocols[m.cursor]
		if ok {
			delete(m.selectedProtocols, m.cursor)
		} else {
			m.selectedProtocols[m.cursor] = struct{}{}
		}
	case "enter":
		if len(m.selectedProtocols) > 0 {
			m.state = StateBenchmarkSelection
			m.cursor = 0
		}
	}
	return m, nil
}

func updateBenchmarkSelection(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.benchmarks)-1 {
			m.cursor++
		}
	case " ":
		_, ok := m.selectedBenchmarks[m.cursor]
		if ok {
			delete(m.selectedBenchmarks, m.cursor)
		} else {
			m.selectedBenchmarks[m.cursor] = struct{}{}
		}
	case "enter":
		if len(m.selectedBenchmarks) > 0 {
			m.state = StateOptions
			m.cursor = 0
		}
	}
	return m, nil
}

func updateOptions(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		m.state = StateRunning
		m.cursor = 0
		return m, runBenchmarks(m)
	}
	return m, nil
}

func updateRunning(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case StateProtocolSelection:
		return viewProtocolSelection(m)
	case StateBenchmarkSelection:
		return viewBenchmarkSelection(m)
	case StateOptions:
		return viewOptions(m)
	case StateRunning:
		return viewRunning(m)
	case StateCompleted:
		return viewCompleted(m)
	}
	return ""
}

func viewProtocolSelection(m model) string {
	s := "Choose protocols for benchmark:\n\n"
	for i, choice := range m.protocols {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selectedProtocols[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}

func viewBenchmarkSelection(m model) string {
	s := "Choose benchmarks to run:\n\n"
	for i, choice := range m.benchmarks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selectedBenchmarks[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}

func viewOptions(m model) string {
	return "Set options for benchmarks (not implemented yet)\nPress enter to run benchmarks."
}

func viewRunning(m model) string {
	return "Running benchmarks...\n" + m.results
}

func viewCompleted(m model) string {
	return fmt.Sprintf("Results:\n%s\n", m.results)
}

func StartUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func runBenchmarks(m model) tea.Cmd {
	return func() tea.Msg {
		var results strings.Builder

		for idx := range m.selectedProtocols {
			protocol := strings.ToLower(m.protocols[idx])

			lb := benchmark.NewLatencyBenchmark(protocol, "localhost:8080", 10)
			if err := lb.Run(); err != nil {
				return benchmarkResultMsg{results: fmt.Sprintf("Error running %s benchmark: %v", protocol, err)}
			}

			results.WriteString(fmt.Sprintf("\n%s Results:\n", m.protocols[idx]))
			results.WriteString(lb.Results())
			results.WriteString("\n")
		}

		return benchmarkResultMsg{results: results.String()}
	}
}
