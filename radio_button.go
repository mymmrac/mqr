package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
)

type RadioButton struct {
	label      string
	onSelected func()
}

func NewRadioButton(label string, onSelected func()) *RadioButton {
	return &RadioButton{
		label:      label,
		onSelected: onSelected,
	}
}

type RadioButtonGroup struct {
	label       string
	buttons     []*RadioButton
	selected    int
	focused     int
	prevFocused int
	cursor      cursor.Model
}

func NewRadioButtonGroup(label string, buttons ...*RadioButton) *RadioButtonGroup {
	if len(buttons) == 0 {
		panic("no buttons in radio group")
	}

	return &RadioButtonGroup{
		label:       label,
		buttons:     buttons,
		selected:    0,
		focused:     -1,
		prevFocused: -2,
		cursor:      cursor.New(),
	}
}

func (g *RadioButtonGroup) Focus() tea.Cmd {
	g.focused = 0
	return g.cursor.Focus()
}

func (g *RadioButtonGroup) Focused() bool {
	return g.focused != -1
}

func (g *RadioButtonGroup) Blur() {
	g.focused = -1
	g.cursor.Blur()
}

func (g *RadioButtonGroup) Select(i int) {
	if i <= 0 || i >= len(g.buttons) {
		panic("select out of range")
	}
	g.selected = i
}

func (g *RadioButtonGroup) Update(untypedMsg tea.Msg) (*RadioButtonGroup, tea.Cmd) {
	if !g.Focused() {
		return g, nil
	}

	switch msg := untypedMsg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			g.focused = max(g.focused-1, 0)
		case tea.KeyRight:
			g.focused = min(g.focused+1, len(g.buttons)-1)
		case tea.KeySpace:
			g.selected = g.focused
			if g.selected >= 0 && g.selected < len(g.buttons) {
				g.buttons[g.selected].onSelected()
			}
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	g.cursor, cmd = g.cursor.Update(untypedMsg)
	cmds = append(cmds, cmd)

	if g.prevFocused != g.focused && g.cursor.Mode() == cursor.CursorBlink {
		g.cursor.Blink = false
		g.prevFocused = g.focused
		cmds = append(cmds, g.cursor.BlinkCmd())
	}

	return g, tea.Batch(cmds...)
}

func (g *RadioButtonGroup) View() string {
	s := strings.Builder{}
	s.WriteString(g.label + ": ")
	for i := 0; i < len(g.buttons); i++ {
		value := " "
		switch {
		case g.focused == i && g.selected == i:
			g.cursor.SetChar("X")
			value = g.cursor.View()
		case g.selected == i:
			value = "X"
		case g.focused == i:
			g.cursor.SetChar(" ")
			value = g.cursor.View()
		}
		s.WriteString(g.buttons[i].label + " [" + value + "]")
		if i != len(g.buttons)-1 {
			s.WriteString(" ")
		}
	}
	return s.String()
}
