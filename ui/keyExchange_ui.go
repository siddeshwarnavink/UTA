package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var KeyExchangeProtocol = []string{
	"Diffie Hellman Key Exchange",
}

type KeyExchangeProtocolModel struct {
	cursor int
	choice string
}

func (m KeyExchangeProtocolModel) Init() tea.Cmd {
	return nil
}

func (m KeyExchangeProtocolModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = KeyExchangeProtocol[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(KeyExchangeProtocol) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(KeyExchangeProtocol) - 1
			}
		}
	}

	return m, nil
}

func (m KeyExchangeProtocolModel) View() string {
	var (
		focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		unfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	)
	s := strings.Builder{}
	s.WriteString("\n\nWhich Key Exchange Algorithm should be used?\n\n")

	for i := 0; i < len(KeyExchangeProtocol); i++ {
		if m.cursor == i {
			s.WriteString(focusedStyle.Render(focusedStyle.Render("> ")))
			s.WriteString(focusedStyle.Render(KeyExchangeProtocol[i]))
		} else {
			s.WriteString(unfocusedStyle.Render("> "))
			s.WriteString(unfocusedStyle.Render(KeyExchangeProtocol[i]))
		}
		s.WriteString("\n")
	}

	return s.String()
}

func RenderKeyProtoForm(keyProtoChan chan string) {
	go func() {
		p := tea.NewProgram(KeyExchangeProtocolModel{})
		m, err := p.Run()
		if err != nil {
			fmt.Println("ERROR from Key Exchange Chooser:", err)
			os.Exit(1)
		}
		if model, ok := m.(KeyExchangeProtocolModel); ok && model.choice != "" {
			keyProtoChan <- model.choice
		} else {
			keyProtoChan <- "error"
		}
	}()
}
