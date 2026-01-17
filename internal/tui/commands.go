package tui

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/adnsv/go-utils/git"
	gitpkg "github.com/adnsv/rtag/internal/git"
	"github.com/adnsv/rtag/internal/version"
	tea "github.com/charmbracelet/bubbletea"
)

// TagInfo represents a tag with its local/remote status
type TagInfo struct {
	Name     string
	IsLocal  bool
	IsRemote bool
}

// msgTagList is sent when tag list has been loaded
type msgTagList struct {
	tags []TagInfo
	err  error
}

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

// loadTagList loads the list of tags with their local/remote status
func loadTagList() tea.Cmd {
	return func() tea.Msg {
		// Get local tags (sorted by date, newest first)
		localCmd := exec.Command("git", "tag", "--sort=-creatordate")
		localOutput, err := localCmd.Output()
		if err != nil {
			return msgTagList{err: fmt.Errorf("failed to get local tags: %w", err)}
		}

		localTags := make(map[string]bool)
		for _, line := range strings.Split(strings.TrimSpace(string(localOutput)), "\n") {
			if line != "" {
				localTags[line] = true
			}
		}

		// Get remote tags
		remoteCmd := exec.Command("git", "ls-remote", "--tags", "origin")
		remoteOutput, _ := remoteCmd.Output() // Ignore error - remote may not exist

		remoteTags := make(map[string]bool)
		for _, line := range strings.Split(string(remoteOutput), "\n") {
			// Format: <sha>\trefs/tags/<tagname>
			// Skip lines ending in ^{} (these are dereferenced tags)
			if strings.Contains(line, "refs/tags/") && !strings.HasSuffix(line, "^{}") {
				parts := strings.Split(line, "refs/tags/")
				if len(parts) == 2 {
					tagName := strings.TrimSpace(parts[1])
					remoteTags[tagName] = true
				}
			}
		}

		// Merge into unified list
		allTags := make(map[string]bool)
		for tag := range localTags {
			allTags[tag] = true
		}
		for tag := range remoteTags {
			allTags[tag] = true
		}

		// Convert to sorted slice (maintain local tag order for newest first)
		var tags []TagInfo
		// First add local tags in order
		for _, line := range strings.Split(strings.TrimSpace(string(localOutput)), "\n") {
			if line != "" {
				tags = append(tags, TagInfo{
					Name:     line,
					IsLocal:  true,
					IsRemote: remoteTags[line],
				})
				delete(allTags, line)
			}
		}
		// Then add any remote-only tags
		for tag := range allTags {
			tags = append(tags, TagInfo{
				Name:     tag,
				IsLocal:  false,
				IsRemote: true,
			})
		}

		return msgTagList{tags: tags}
	}
}