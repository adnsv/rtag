package main

import (
	"fmt"
	"os"

	"github.com/adnsv/go-utils/git"
)

func get_stats() (wd string, stats *git.Stats, err error) {
	wd, err = os.Getwd()
	if err != nil {
		return
	}
	stats, err = git.Stat(wd)
	if stats != nil {
		fmt.Println("Repository Info:")
		print_keyval("branch", stats.Branch)
		print_keyval("author date", stats.AuthorDate)
		print_keyval("hash", stats.Hash)
	}
	return
}
