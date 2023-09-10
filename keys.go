package main

import "github.com/charmbracelet/bubbles/key"

type keybindings struct {
	ForceQuit key.Binding
	Quit      key.Binding
	Done      key.Binding
}

var keys = keybindings{
	ForceQuit: key.NewBinding(key.WithKeys("ctrl+c", "ctrl+d")),
	Quit:      key.NewBinding(key.WithKeys("esc")),
	Done:      key.NewBinding(key.WithKeys("enter")),
}
