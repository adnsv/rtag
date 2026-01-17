package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the TUI application
func Run(opts Options) error {
	model := NewModel(opts)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}
	
	// Check if we exited with an error
	if m, ok := finalModel.(Model); ok && m.err != nil {
		return m.err
	}
	
	return nil
}