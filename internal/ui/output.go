package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/adnsv/go-utils/ansi"
)

func ConfigureOutput(output *os.File) *ansi.OutputState {
	termstate := ansi.SetupOutput(output)
	if termstate.Supported() {
		FmtTag = func(s string) string {
			return ansi.Underline + s + ansi.Reset
		}
		FmtBold = func(s string) string {
			return ansi.Bold + s + ansi.Reset
		}
		FmtDim = func(s string) string {
			return ansi.Dim + s + ansi.Reset
		}
		BeginDim = func() {
			fmt.Printf(ansi.Dim)
		}
		EndDim = func() {
			fmt.Printf(ansi.Reset)
		}
	}
	return termstate
}

var FmtTag = func(s string) string {
	return fmt.Sprintf("'%s'", s)
}

var FmtBold = func(s string) string {
	return s
}
var FmtDim = func(s string) string {
	return s
}
var FmtKeyval = func(k string, v any) string {
	sp := 18 - len(k)
	if sp < 0 {
		sp = 0
	}
	return fmt.Sprintf("- %s: %s%v", k, strings.Repeat(" ", sp), v)
}

func PrintKeyval(k string, v any) {
	fmt.Println(FmtKeyval(k, v))
}

func PrintDim(s string) {
	fmt.Println(FmtDim(s))
}

var BeginDim = func() {}

var EndDim = func() {}
