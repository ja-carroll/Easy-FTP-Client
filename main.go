package main

import (
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type State int

const (
	connectionForm State = iota
	connecting
	welcoming
	uploadButtonActive
	downloadButtonActive
)

var (
	accentColor         = lipgloss.Color("99")
	yellowColor         = lipgloss.Color("#ECFD66")
	activeLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	uploadButtonStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Background(lipgloss.Color("#F25D94")).Padding(0, 3).MarginTop(1).MarginRight(2).Underline(true)
	downloadButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Background(lipgloss.Color("#F25D94")).Padding(0, 3).MarginTop(1)

	// dialog
	dialogBoxStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("#874BFD")).Padding(1, 0).BorderTop(true).BorderBottom(true).BorderLeft(true).BorderRight(true)
	// This is essentially the container page
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type welcomeScreenMsg struct{}
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
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	if m.state == connecting {
		return "\n " + m.spinner.View() + "Connecting to host"
	} else if m.state == welcoming || m.state == downloadButtonActive || m.state == uploadButtonActive {
		uploadButton := uploadButtonStyle
		downloadButton := downloadButtonStyle

		if m.state == uploadButtonActive {
			uploadButton = uploadButtonStyle.Background(accentColor).Foreground(yellowColor)
		}

		if m.state == downloadButtonActive {
			downloadButton = downloadButtonStyle.Background(accentColor).Foreground(yellowColor)
		}

		bannerMsg := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Easy-FTP-Client")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, uploadButton.Render("Upload"), downloadButton.Render("Download"))
		ui := lipgloss.JoinVertical(lipgloss.Center, bannerMsg, buttons)

		dialog := lipgloss.Place(100, 10, lipgloss.Center, lipgloss.Center, dialogBoxStyle.Render(ui))
		if physicalWidth > 0 {
			docStyle = docStyle.MaxWidth(physicalWidth)
		}

		return docStyle.Render(dialog)
	}
	return m.form.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connectionAchievedMsg:
		m.state = welcoming
		return m, displayWelcomeScreen()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "esc", "q"))):
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			switch m.state {
			case welcoming:
				m.state = uploadButtonActive
			case uploadButtonActive:
				m.state = downloadButtonActive
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			switch m.state {
			case downloadButtonActive:
				m.state = uploadButtonActive
			}
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

func displayWelcomeScreen() tea.Cmd {
	return func() tea.Msg {
		return welcomeScreenMsg{}
	}
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
