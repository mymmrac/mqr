package main

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
)

func runTUI(app *cli.Context, data string) error {
	program := tea.NewProgram(newModel(data), tea.WithContext(app.Context))
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

type model struct {
	input              textinput.Model
	sizeInput          *RadioButtonGroup
	invertedInput      *RadioButtonGroup
	recoveryLevelInput *RadioButtonGroup

	big           bool
	inverted      bool
	recoveryLevel qrcode.RecoveryLevel
	output        string

	err error

	width, height int

	exit         bool
	updateNeeded bool
}

func newModel(data string) *model {
	input := textinput.New()
	input.Placeholder = "..."
	input.Prompt = "Data: "
	input.SetValue(data)

	m := &model{
		input:         input,
		recoveryLevel: qrcode.Medium,
		updateNeeded:  true,
	}

	m.sizeInput = NewRadioButtonGroup("Size",
		NewRadioButton("small", func() {
			m.big = false
		}),
		NewRadioButton("big", func() {
			m.big = true
		}),
	)

	m.invertedInput = NewRadioButtonGroup("Color",
		NewRadioButton("regular", func() {
			m.inverted = false
		}),
		NewRadioButton("inverted", func() {
			m.inverted = true
		}),
	)

	m.recoveryLevelInput = NewRadioButtonGroup("Recovery level",
		NewRadioButton("low", func() {
			m.recoveryLevel = qrcode.Low
		}),
		NewRadioButton("medium", func() {
			m.recoveryLevel = qrcode.Medium
		}),
		NewRadioButton("high", func() {
			m.recoveryLevel = qrcode.High
		}),
		NewRadioButton("highest", func() {
			m.recoveryLevel = qrcode.Highest
		}),
	)
	m.recoveryLevelInput.Select(1)

	return m
}

func (m *model) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *model) Update(untypedMsg tea.Msg) (tea.Model, tea.Cmd) {
	if m.exit {
		return m, tea.Quit
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := untypedMsg.(type) {
	case tea.KeyMsg:
		if !slices.Contains([]tea.KeyType{tea.KeyRight, tea.KeyLeft, tea.KeyUp, tea.KeyDown}, msg.Type) {
			m.updateNeeded = true
		}

		switch {
		case key.Matches(msg, keys.ForceQuit):
			m.exit = true
			m.output = ""
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
		case msg.Type == tea.KeyUp:
			switch {
			case m.sizeInput.Focused():
				m.sizeInput.Blur()
				cmds = append(cmds, m.input.Focus())
			case m.invertedInput.Focused():
				m.invertedInput.Blur()
				cmds = append(cmds, m.sizeInput.Focus())
			case m.recoveryLevelInput.Focused():
				m.recoveryLevelInput.Blur()
				cmds = append(cmds, m.invertedInput.Focus())
			}
		case msg.Type == tea.KeyDown:
			switch {
			case m.input.Focused():
				m.input.Blur()
				cmds = append(cmds, m.sizeInput.Focus())
			case m.sizeInput.Focused():
				m.sizeInput.Blur()
				cmds = append(cmds, m.invertedInput.Focus())
			case m.invertedInput.Focused():
				m.invertedInput.Blur()
				cmds = append(cmds, m.recoveryLevelInput.Focus())
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.input.Width = m.width - logoWidth - 14
	}

	m.input, cmd = m.input.Update(untypedMsg)
	cmds = append(cmds, cmd)

	m.sizeInput, cmd = m.sizeInput.Update(untypedMsg)
	cmds = append(cmds, cmd)

	m.invertedInput, cmd = m.invertedInput.Update(untypedMsg)
	cmds = append(cmds, cmd)

	m.recoveryLevelInput, cmd = m.recoveryLevelInput.Update(untypedMsg)
	cmds = append(cmds, cmd)

	if m.updateNeeded {
		data := m.input.Value()
		if data == "" {
			m.err = nil
			m.output = ""
		} else {
			code, err := qrCodeFromData(m.input.Value(), 0, m.recoveryLevel)
			if err != nil {
				m.err = err
				m.output = ""
			} else {
				m.err = nil

				if m.big {
					m.output = code.ToString(m.inverted)
				} else {
					m.output = code.ToSmallString(m.inverted)
				}
			}
		}

		m.updateNeeded = false
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.exit {
		if m.output == "" {
			return "Bye!\n"
		}
		return "\n " + strings.ReplaceAll(m.output, "\n", "\n ") + "\n"
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
		logoStyle.Render(logo), inputStyle.Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.input.View(),
			"",
			m.sizeInput.View(),
			"",
			m.invertedInput.View(),
			"",
			m.recoveryLevelInput.View(),
		)),
	)
	s.WriteString(mainBorder.Width(m.width - 2).Render(inputBox))

	return s.String()
}
