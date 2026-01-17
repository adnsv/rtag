package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/adnsv/rtag/internal/ui"
	cli "github.com/jawher/mow.cli"
)

var errUserCancelled = errors.New("exiting")

func main() {
	app := cli.App("rtag", "rtag is a git tag management utility that helps making consistent release tags")

	app.Version("version", app_version())

	app.Spec = "[--undo] [--prefix=<ver-prefix>] [--allow-dirty] [--classic]"

	undo := false
	classic := false
	opts := exec_options{}

	app.BoolOptPtr(&undo, "u undo", false, "undo last tag locally and remotely")
	app.BoolOptPtr(&classic, "classic", false, "use classic prompt-based interface")
	opts.bind_cli(app)

	app.Action = func() {
		termstate := ui.ConfigureOutput(os.Stdout)
		defer termstate.Restore()

		var err error

		if undo {
			if classic {
				err = cmd_undo()
			} else {
				err = cmdUndoTUI()
			}
		} else {
			if classic {
				err = execute(&opts)
			} else {
				err = executeTUI(&opts)
			}
		}

		termstate.Restore()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Println()
			fmt.Println("mission accomplished")
		}
	}

	app.Run(os.Args)
}
