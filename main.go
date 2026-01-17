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

	app.Spec = "[--prefix=<ver-prefix>]"

	prefix := "AUTO"

	app.StringOptPtr(&prefix, "p prefix", "AUTO", "prefix for new tags")

	app.Action = func() {
		termstate := ui.ConfigureOutput(os.Stdout)
		defer termstate.Restore()

		opts := tui.Options{
			Prefix: prefix,
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
