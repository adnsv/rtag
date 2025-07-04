package git

import (
	"fmt"
	"os"

	"github.com/adnsv/go-utils/git"
	"github.com/adnsv/rtag/internal/ui"
)

func GetStats() (wd string, stats *git.Stats, err error) {
	wd, err = os.Getwd()
	if err != nil {
		return
	}
	stats, err = git.Stat(wd)
	if stats != nil {
		fmt.Println("Repository Info:")
		ui.PrintKeyval("branch", stats.Branch)
		ui.PrintKeyval("author date", stats.AuthorDate)
		ui.PrintKeyval("hash", stats.Hash)
	}
	return
}
