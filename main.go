package main

import (
	"fmt"
	"os"

	"github.com/adnsv/rtag/internal/tui"
	"github.com/adnsv/rtag/internal/ui"
	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("rtag", "rtag is a git tag management utility that helps making consistent release tags")

	app.Version("version", app_version())

	app.Spec = "[--undo] [--prefix=<ver-prefix>] [--allow-dirty]"

	undo := false
	prefix := "AUTO"
	allowDirty := false

	app.BoolOptPtr(&undo, "u undo", false, "undo last tag locally and remotely")
	app.StringOptPtr(&prefix, "p prefix", "AUTO", "prefix for new tags")
	app.BoolOptPtr(&allowDirty, "d allow-dirty", false, "allow tagging of repos that contain uncommited changes")

	app.Action = func() {
		termstate := ui.ConfigureOutput(os.Stdout)
		defer termstate.Restore()

		opts := tui.Options{
			Prefix:     prefix,
			AllowDirty: allowDirty,
			Undo:       undo,
		}

		err := tui.Run(opts)

		termstate.Restore()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
