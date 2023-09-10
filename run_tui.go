package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
)

func runTUI(app *cli.Context) error {
	program := tea.NewProgram(newModel(), tea.WithContext(app.Context))
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

type model struct {
	input textinput.Model

	output string

	err error

	width, height int

	exit bool
}

func newModel() *model {
	input := textinput.New()
	input.Placeholder = "..."
	input.Prompt = "Data: "
	input.Focus()

	return &model{
		input: input,
	}
}

func (m *model) Init() tea.Cmd {
	return m.input.Cursor.BlinkCmd()
}

func (m *model) Update(untypedMsg tea.Msg) (tea.Model, tea.Cmd) {
	if m.exit {
		return m, tea.Quit
	}

	keyUpdate := false
	switch msg := untypedMsg.(type) {
	case tea.KeyMsg:
		if msg.Type != tea.KeyLeft && msg.Type != tea.KeyRight {
			keyUpdate = true
		}

		switch {
		case key.Matches(msg, keys.ForceQuit):
			m.exit = true
			return m, nil
		case key.Matches(msg, keys.Quit):
			if m.input.Value() == "" {
				m.exit = true
				return m, nil
			}
			m.input.SetValue("")
		case key.Matches(msg, keys.Done):
			if m.output != "" {
				m.exit = true
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(untypedMsg)

	if keyUpdate {
		data := m.input.Value()
		if data == "" {
			m.err = nil
			m.output = ""
		} else {
			code, err := qrCodeFromData(m.input.Value(), 0, qrcode.Medium)
			if err != nil {
				m.err = err
				m.output = ""
			} else {
				m.err = nil
				m.output = code.ToSmallString(false)
			}
		}
	}

	return m, inputCmd
}

func (m *model) View() string {
	if m.exit {
		if m.output == "" {
			return "Bye!\n"
		}
		return m.output + "\n"
	}

	s := strings.Builder{}

	output := m.output
	if m.err != nil {
		output = "Error: " + m.err.Error()
		outputStyle = outputStyle.PaddingBottom(1)
	} else {
		if output == "" {
			output = "Start typing and your QR code will appear here."
			outputStyle = outputStyle.PaddingBottom(1)
		} else {
			outputStyle = outputStyle.PaddingBottom(0)
		}
	}
	s.WriteString(outputStyle.Render(strings.TrimSuffix(output, "\n")))
	s.WriteString("\n")

	inputBox := lipgloss.JoinHorizontal(lipgloss.Top,
		logoStyle.Render(logo), inputStyle.Render(m.input.View()),
	)
	s.WriteString(mainBorder.Width(m.width - 2).Render(inputBox))

	return s.String()
}
