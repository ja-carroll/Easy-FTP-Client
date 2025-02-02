package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	connectionForm State = iota
	connecting
)

var (
	activeLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
)

type connectionAchievedMsg struct{}

type (
	errMsg error
)

type model struct {
	state   State
	inputs  []textinput.Model
	focused int
	err     error
	spinner spinner.Model
	form    *huh.Form // huh.Form is just a tea.Model
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) View() string {
	// return fmt.Sprintf(
	//		`
	//    Welcome to Easy-ftp Client!
	//    Please enter the Hostname of the ftp server you wish to connect to
	//
	//    %s
	//    %s
	//    %s
	//
	//    %s`,
	//		m.inputs[0].View(),
	//		m.inputs[1].View(),
	//		m.inputs[2].View(),
	//		"(esc to quit)",
	//	) + "\n"

	if m.state == connecting {
		return "\n " + m.spinner.View() + "Connecting to host"
	}
	return m.form.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	//
	//	switch msg := msg.(type) {
	//	case tea.KeyMsg:
	//		switch msg.Type {
	//
	//		case tea.KeyEnter:
	//			if m.focused == len(m.inputs)-1 {
	//				return m, tea.Quit
	//			}
	//			// Increment to the next input if we are not at the last input
	//			m.nextInput()
	//
	//		case tea.KeyCtrlC, tea.KeyEsc:
	//			return m, tea.Quit
	//
	//		case tea.KeyDown:
	//			m.nextInput()
	//		case tea.KeyUp:
	//			m.prevInput()
	//		}
	//
	//		for i := range m.inputs {
	//			m.inputs[i].Blur()
	//		}
	//		m.inputs[m.focused].Focus()
	//
	//	case errMsg:
	//		m.err = msg
	//		return m, nil
	//	}
	//
	//	for i := range m.inputs {
	//		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	//	}
	//
	//	return m, tea.Batch(cmds...)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.state {
	case connectionForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}
		if m.form.State == huh.StateCompleted {
			m.state = connecting
			return m, tea.Batch(m.spinner.Tick, m.MakeConnection())
		}
	case connecting:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
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

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Hostname").Key("host"),
			huh.NewInput().Title("User").Key("key"),
			huh.NewInput().Title("Password").Password(true).Key("pass"),
		),
	)

	loadingSpinner := spinner.New()
	loadingSpinner.Style = activeLabelStyle
	loadingSpinner.Spinner = spinner.Dot

	return model{
		state:   connectionForm,
		inputs:  inputs,
		focused: 0,
		err:     nil,
		form:    form,
		spinner: loadingSpinner,
	}
}
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) MakeConnection() tea.Cmd {
	return func() tea.Msg {
		time.After(5 * time.Second)
		return connectionAchievedMsg{}
	}
}
