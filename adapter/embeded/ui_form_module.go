package embeded

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

type FormModel struct {
	textInput textinput.Model
	err       error
	question  string
}

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))

// Function to initialize the model with any question and placeholder
func initialModel(question, placeholder string) FormModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.TextStyle = focusedStyle
	ti.CursorStyle = cursorStyle // Set the style of the cursor

	return FormModel{
		textInput: ti,
		err:       nil,
		question:  question,
	}
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// Handle errors if any
	case errMsg:
		m.err = msg
		return m, nil
	}

	// Update the text input model
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m FormModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n", // Display the question dynamically
		m.question,
		m.textInput.View(),
	) + "\n"
}

func RenderForm(question string, options []string, placeholder string) (string, error) {
	p := tea.NewProgram(initialModel(question, placeholder))
	model, err := p.Run()
	if err != nil {
		return "", err
	}
	if qModel, ok := model.(FormModel); ok {
		return qModel.textInput.Value(), nil
	}
	return "", fmt.Errorf("could not get the answer")
}

func form(L *lua.LState) int {
	question := L.ToString(1)
	luaTable := L.ToTable(2)
	options := []string{}
	luaTable.ForEach(func(_, value lua.LValue) {
		options = append(options, value.String())
	})
	placeholder := L.ToString(3)
	answer, err := RenderForm(question, options, placeholder)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LString(answer))
	}
	return 1
}

func UIFormLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"new": form,
	}
	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
