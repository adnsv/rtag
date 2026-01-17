package tui

import (
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/adnsv/go-utils/git"
	gitpkg "github.com/adnsv/rtag/internal/git"
	"github.com/adnsv/rtag/internal/version"
	tea "github.com/charmbracelet/bubbletea"
)

// loadRepoStats loads repository statistics
func loadRepoStats() tea.Msg {
	wd, stats, err := gitpkg.GetStats()
	return msgRepoStats{
		workDir: wd,
		stats:   stats,
		err:     err,
	}
}

// parseVersions parses the current version and generates available actions
func parseVersions(stats *git.Stats, prefix string) tea.Cmd {
	return func() tea.Msg {
		if stats == nil || stats.Description.Tag == "" {
			return msgVersionsParsed{
				versions: []version.Action{},
				prefix:   prefix,
			}
		}
		
		vi, err := git.ParseVersion(stats.Description)
		if err != nil {
			// Handle parse error - for now just return empty
			return msgVersionsParsed{
				versions: []version.Action{},
				prefix:   prefix,
			}
		}
		
		versions := version.CollectActions(vi.Semantic)
		return msgVersionsParsed{
			versions: versions,
			prefix:   prefix,
		}
	}
}

// executeGitCommand executes a git command and captures output
func executeGitCommand(args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("git", args...)

		// Get pipes for stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return msgCommandComplete{err: err}
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return msgCommandComplete{err: err}
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return msgCommandComplete{err: err}
		}

		// Read both outputs concurrently with proper synchronization
		var wg sync.WaitGroup
		var stdoutBytes, stderrBytes []byte
		var stdoutErr, stderrErr error

		wg.Add(2)
		go func() {
			defer wg.Done()
			stdoutBytes, stdoutErr = io.ReadAll(stdout)
		}()
		go func() {
			defer wg.Done()
			stderrBytes, stderrErr = io.ReadAll(stderr)
		}()

		// Wait for readers to finish before calling Wait
		wg.Wait()

		// Wait for command to complete
		err = cmd.Wait()

		// Combine output
		output := string(stdoutBytes) + string(stderrBytes)

		// Check for read errors
		if stdoutErr != nil && err == nil {
			err = stdoutErr
		}
		if stderrErr != nil && err == nil {
			err = stderrErr
		}

		return msgCommandComplete{
			output: output,
			err:    err,
		}
	}
}

// gitTag creates a new tag
func gitTag(tag, comment string) tea.Cmd {
	return executeGitCommand("tag", "-a", tag, "-m", comment)
}

// gitPushTag pushes a tag to origin
func gitPushTag(tag string) tea.Cmd {
	return executeGitCommand("push", "origin", tag)
}

// gitDeleteLocalTag deletes a local tag
func gitDeleteLocalTag(tag string) tea.Cmd {
	return executeGitCommand("tag", "-d", tag)
}

// gitDeleteRemoteTag deletes a remote tag
func gitDeleteRemoteTag(tag string) tea.Cmd {
	return executeGitCommand("push", "--delete", "origin", tag)
}

// generateTagComment generates a comment for the tag
func generateTagComment(tag string) string {
	return fmt.Sprintf("tagging as %s", tag)
}

// gitAddAndCommit stages all changes and commits
func gitAddAndCommit(message string) tea.Cmd {
	return func() tea.Msg {
		// First, add all files
		cmd := exec.Command("git", "add", "-A")
		if err := cmd.Run(); err != nil {
			return msgCommandComplete{err: fmt.Errorf("git add failed: %w", err)}
		}
		
		// Then commit
		return executeGitCommand("commit", "-m", message)()
	}
}

// gitStash stashes changes with optional message
func gitStash(message string) tea.Cmd {
	if message == "" {
		return executeGitCommand("stash")
	}
	return executeGitCommand("stash", "save", message)
}