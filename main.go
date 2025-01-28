package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
	return fmt.Sprintf(
		`
    Welcome to Easy-ftp Client!
    Please enter the Hostname of the ftp server you wish to connect to

    %s
    %s
    %s

    %s`,
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
		"(esc to quit)",
	) + "\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			// Increment to the next input if we are not at the last input
			m.nextInput()

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyDown:
			m.nextInput()
		case tea.KeyUp:
			m.prevInput()
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Hostname"
	inputs[0].CharLimit = 200
	inputs[0].Prompt = ""
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "User"
	inputs[1].Prompt = ""

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Pass"
	inputs[2].Prompt = ""

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
