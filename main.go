package main

import (
	"errors"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

var errUserCancelled = errors.New("exiting")

func main() {
	app := cli.App("rtag", "rtag is a git tag management utility that helps making consistent release tags")

	app.Version("version", app_version())

	app.Spec = "[--undo] [--prefix=<ver-prefix>] [--allow-dirty]"

	undo := false
	opts := exec_options{}

	app.BoolOptPtr(&undo, "u undo", false, "undo last tag locally and remotely")
	opts.bind_cli(app)

	app.Action = func() {

		if version {
			show_app_version()
			return
		}

		termstate := configure_output(os.Stdout)
		defer termstate.Restore()

		var err error

		if undo {
			err = cmd_undo()
		} else {
			err = execute(&opts)
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
