package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mcqModel struct {
	cursor   int
	choice   string
	question string
	options  []string
}

func (m mcqModel) Init() tea.Cmd {
	return nil
}

func (m mcqModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = m.options[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}
		}
	}

	return m, nil
}

func (m mcqModel) View() string {
	var (
		focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		unfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	)
	s := strings.Builder{}
	s.WriteString("\n\n" + m.question + "\n\n")

	for i := 0; i < len(m.options); i++ {
		if i == m.cursor {
			s.WriteString(focusedStyle.Render("> "))
			s.WriteString(focusedStyle.Render(m.options[i]))
		} else {
			s.WriteString(unfocusedStyle.Render("> "))
			s.WriteString(unfocusedStyle.Render(m.options[i]))
		}
		s.WriteString("\n")
	}
	return s.String()
}

func MCQ(question string, options []string, MCQChan chan string) {
	go func() {
		model := mcqModel{
			question: question,
			options:  options,
		}
		p := tea.NewProgram(model)
		m, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not run program: %v", err)
			os.Exit(1)
		}
		if m, ok := m.(mcqModel); ok && m.choice != "" {
			MCQChan <- m.choice
		} else {
			MCQChan <- "error"
		}
	}()

}
