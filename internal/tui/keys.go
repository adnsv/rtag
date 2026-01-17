package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all key bindings for the application
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns key bindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

// FullHelp returns key bindings for the full help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Back, k.Quit},
	}
}

// ListHelp returns help for list navigation screens
func (k KeyMap) ListHelp() string {
	return "↑/↓: navigate • enter: select • q: quit"
}

// ListHelpWithBack returns help for list screens that support going back
func (k KeyMap) ListHelpWithBack() string {
	return "↑/↓: navigate • enter: select • esc: back • q: quit"
}

// ConfirmHelp returns help for confirmation screens
func (k KeyMap) ConfirmHelp() string {
	return "↑/↓: navigate • enter: select"
}

// InputHelp returns help for text input screens
func (k KeyMap) InputHelp() string {
	return "enter: confirm • esc: cancel"
}

// DoneHelp returns help for done/error screens
func (k KeyMap) DoneHelp() string {
	return "enter/q: exit"
}
