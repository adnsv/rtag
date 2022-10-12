package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/adnsv/go-utils/prompt"
)

func cmd_undo() error {

	_, stats, err := get_stats()
	if err != nil {
		return err
	}

	fmt.Printf("deleting tag: %s\n", stats.Description.Tag)
	fmt.Println()
	response := prompt.Enum("choose which tag to delete", "local", "remote", "both")
	loc := response == "local" || response == "both"
	rem := response == "remote" || response == "both"
	fmt.Println()
	fmt.Printf("ready to execute the following commands:\n")
	if loc {
		fmt.Printf("- 'git tag -d %s'\n", stats.Description.Tag)
	}
	if rem {
		fmt.Printf("- 'git push --delete origin %s'\n", stats.Description.Tag)
	}
	if !loc && !rem {
		return errUserCancelled
	}

	fmt.Println()
	if !prompt.YN("proceed [y/n]?") {
		return errUserCancelled
	}

	fmt.Println()
	if loc {
		fmt.Printf("deleting local tag %s\n", stats.Description.Tag)
		cmd := exec.Command("git", "tag", "-d", stats.Description.Tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", cmd.String(), err)
		}
	}
	if rem {
		fmt.Printf("deleting remote origin tag %s\n", stats.Description.Tag)
		cmd := exec.Command("git", "push", "--delete", "origin", stats.Description.Tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", cmd.String(), err)
		}
	}

	return nil
}
