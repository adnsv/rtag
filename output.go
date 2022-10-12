package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/adnsv/go-utils/ansi"
)

func configure_output(output *os.File) *ansi.OutputState {
	termstate := ansi.SetupOutput(output)
	if termstate.Supported() {
		fmt_tag = func(s string) string {
			return ansi.Underline + s + ansi.Reset
		}
		fmt_bold = func(s string) string {
			return ansi.Bold + s + ansi.Reset
		}
		fmt_dim = func(s string) string {
			return ansi.Dim + s + ansi.Reset
		}
		begin_dim = func() {
			fmt.Printf(ansi.Dim)
		}
		end_dim = func() {
			fmt.Printf(ansi.Reset)
		}
	}
	return termstate
}

var fmt_tag = func(s string) string {
	return fmt.Sprintf("'%s'", s)
}

var fmt_bold = func(s string) string {
	return s
}
var fmt_dim = func(s string) string {
	return s
}
var fmt_keyval = func(k string, v any) string {
	sp := 18 - len(k)
	if sp < 0 {
		sp = 0
	}
	return fmt.Sprintf("- %s: %s%v", k, strings.Repeat(" ", sp), v)
}

func print_keyval(k string, v any) {
	fmt.Println(fmt_keyval(k, v))
}

func print_dim(s string) {
	fmt.Println(fmt_dim(s))
}

var begin_dim = func() {}

var end_dim = func() {}
