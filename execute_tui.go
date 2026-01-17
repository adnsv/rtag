package main

import (
	"github.com/adnsv/rtag/internal/tui"
)

// executeTUI runs the interactive TUI version
func executeTUI(opts *exec_options) error {
	tuiOpts := tui.Options{
		Prefix:     opts.prefix,
		AllowDirty: opts.allow_dirty,
		Undo:       false,
	}
	
	return tui.Run(tuiOpts)
}

// cmdUndoTUI runs the undo flow in TUI mode
func cmdUndoTUI() error {
	tuiOpts := tui.Options{
		Undo: true,
	}
	
	return tui.Run(tuiOpts)
}