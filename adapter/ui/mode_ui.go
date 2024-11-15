package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var modes = []string{"Client", "Server"}

type ModeModel struct {
	cursor int
	choice string
}

func (m ModeModel) Init() tea.Cmd {
	return nil
}

func (m ModeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = modes[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(modes) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(modes) - 1
			}
		}
	}

	return m, nil
}

func (m ModeModel) View() string {
	var (
		focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		unfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	)

	s := strings.Builder{}
	s.WriteString("Which mode is this system on?\n\n")

	for i := 0; i < len(modes); i++ {
		if m.cursor == i {
			s.WriteString(focusedStyle.Render("> "))
			s.WriteString(focusedStyle.Render(modes[i]))
		} else {
			s.WriteString(unfocusedStyle.Render("> "))
			s.WriteString(unfocusedStyle.Render(modes[i]))
		}
		s.WriteString("\n")
	}

	return s.String()
}

func RenderModeForm(ModeChan chan string) {
	go func() {
		p := tea.NewProgram(ModeModel{})

		m, err := p.Run()
		if err != nil {
			fmt.Println("ERORR from Mode Chooser:", err)
			os.Exit(1)
		}

		if m, ok := m.(ModeModel); ok && m.choice != "" {
			ModeChan <- m.choice
		} else {
			ModeChan <- "error"
		}
	}()
}
