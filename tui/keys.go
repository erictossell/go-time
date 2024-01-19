package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	start       key.Binding
	stop        key.Binding
	up          key.Binding
	down        key.Binding
	list        key.Binding
	quit        key.Binding
	selectTimer key.Binding
}
