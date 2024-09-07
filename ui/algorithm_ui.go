package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var Algorithms = []string{
	"Advanced Encryption Standard(AES)",
	"ChaCha20",
	"TwoFish",
}

type AlgorithmModel struct {
	cursor int
	choice string
}

func (m AlgorithmModel) Init() tea.Cmd {
	return nil
}

func (m AlgorithmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = Algorithms[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(Algorithms) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(Algorithms) - 1
			}
		}
	}

	return m, nil
}

func (m AlgorithmModel) View() string {
	var (
		focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		unfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	)

	s := strings.Builder{}
	s.WriteString("\n\nWhich Cryptographic Algorithm should be used?\n\n")

	for i := 0; i < len(Algorithms); i++ {
		if m.cursor == i {
			s.WriteString(focusedStyle.Render("> "))
			s.WriteString(focusedStyle.Render(Algorithms[i]))
		} else {
			s.WriteString(unfocusedStyle.Render("> "))
			s.WriteString(unfocusedStyle.Render(Algorithms[i]))
		}
		s.WriteString("\n")
	}

	return s.String()
}

func RenderAlgoForm(AlgoChan chan string) {
	go func() {
		p := tea.NewProgram(AlgorithmModel{})
		m, err := p.Run()
		if err != nil {
			fmt.Println("ERROR from Algorithm Chooser:", err)
			os.Exit(1)
		}
		if model, ok := m.(AlgorithmModel); ok && model.choice != "" {
			AlgoChan <- model.choice
		} else {
			AlgoChan <- "error"
		}
	}()
}
