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
	inputs []textinput.Model
	err    error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
	return fmt.Sprintf(
		"Welcome to Easy-ftp Client!\n\nPlease enter the Hostname of the ftp server you wish to connect to\n\n%s\n\n%s",
		m.inputs[0].View(),
		"(esc to quit)",
	) + "\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.inputs[0], cmd = m.inputs[0].Update(msg)

	return m, cmd
}

func main() {
	ti := textinput.New()
	ti.Placeholder = "Hostname"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 20

	mod := model{
		inputs: []textinput.Model{ti},
		err:    nil,
	}
	p := tea.NewProgram(mod)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
