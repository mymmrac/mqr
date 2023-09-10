package main

import "github.com/charmbracelet/lipgloss"

var (
	mainBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder())

	logo = "" +
		"███╗   ███╗ ██████╗ ██████╗ " + "\n" +
		"████╗ ████║██╔═══██╗██╔══██╗" + "\n" +
		"██╔████╔██║██║   ██║██████╔╝" + "\n" +
		"██║╚██╔╝██║██║▄▄ ██║██╔══██╗" + "\n" +
		"██║ ╚═╝ ██║╚██████╔╝██║  ██║" + "\n" +
		"╚═╝     ╚═╝ ╚══▀▀═╝ ╚═╝  ╚═╝" + "\n"

	logoStyle = lipgloss.NewStyle().
			Padding(1).
			PaddingBottom(0)

	inputStyle = lipgloss.NewStyle().
			Height(lipgloss.Height(logo)+1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).Padding(0, 1)

	outputStyle = lipgloss.NewStyle().
			Padding(1)
)
