package tui

import (
	"fmt"
	"strings"
)

func (m Model) viewDirtyRepoChoice() string {
	var s strings.Builder

	// Repo info header
	s.WriteString(titleStyle.Render("Repository Info") + "\n")

	if m.stats != nil {
		s.WriteString(m.renderKeyValue("branch", m.stats.Branch))

		if m.stats.Description.Tag != "" {
			s.WriteString(m.renderKeyValue("last tag", tagStyle.Render(m.stats.Description.Tag)))
			if m.stats.Description.AdditionalCommits > 0 {
				s.WriteString(m.renderKeyValue("additional commits", fmt.Sprint(m.stats.Description.AdditionalCommits)))
			}
		}

		s.WriteString(m.renderKeyValue("state", warningStyle.Render("dirty")))

		if m.prefix != "" {
			s.WriteString(m.renderKeyValue("prefix", m.prefix))
		}
	}

	// Warning message
	s.WriteString("\n" + warningStyle.Render("⚠  Your repository has uncommitted changes.") + "\n")
	s.WriteString(infoStyle.Render("Tags should typically be created on clean repositories.") + "\n\n")

	s.WriteString(m.dirtyChoiceList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ListHelp()))

	return s.String()
}

func (m Model) viewUndoPreview() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Undo Tag") + "\n\n")

	// Show loading state if tag list isn't loaded yet
	if m.tagList == nil {
		s.WriteString(fmt.Sprintf(" %s Loading tag information...\n", m.spinner.View()))
		return s.String()
	}

	// Show recent tags with their status
	s.WriteString("Current tags (newest first):\n")
	s.WriteString(infoStyle.Render("─────────────────────────────") + "\n")

	maxTags := 5 // Show up to 5 recent tags
	for i, tag := range m.tagList {
		if i >= maxTags {
			break
		}

		// Build status string
		var status string
		if tag.IsLocal && tag.IsRemote {
			status = "local + remote"
		} else if tag.IsLocal {
			status = "local only"
		} else {
			status = "remote only"
		}

		// Mark the tag being deleted
		marker := "  "
		suffix := ""
		if tag.Name == m.currentTag {
			marker = "  "
			suffix = " ← will be deleted"
		}

		s.WriteString(fmt.Sprintf("%s%-12s  %s%s\n", marker, tagStyle.Render(tag.Name), infoStyle.Render(status), warningStyle.Render(suffix)))
	}

	// Show what the new latest tag will be after deletion
	s.WriteString("\n")
	if len(m.tagList) > 1 {
		// Find next tag after the one being deleted
		for i, tag := range m.tagList {
			if tag.Name == m.currentTag && i+1 < len(m.tagList) {
				s.WriteString(infoStyle.Render(fmt.Sprintf("After deletion, latest tag will be: %s", tagStyle.Render(m.tagList[i+1].Name))) + "\n")
				break
			}
		}
	}

	s.WriteString("\n" + titleStyle.Render("Select deletion scope") + "\n\n")
	s.WriteString(m.undoScopeList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ListHelpWithBack()))

	return s.String()
}

func (m Model) viewConfirmUndo() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Confirm Tag Deletion") + "\n\n")

	s.WriteString(m.renderKeyValue("tag", tagStyle.Render(m.currentTag)))
	s.WriteString(m.renderKeyValue("scope", infoStyle.Render(m.undoScope)))

	s.WriteString("\n")
	s.WriteString(m.confirmList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ConfirmHelp()))

	return s.String()
}

func (m Model) viewCommitMessage() string {
	var s strings.Builder
	
	s.WriteString(titleStyle.Render("Commit Changes") + "\n\n")
	s.WriteString("Enter commit message:\n\n")
	s.WriteString(m.textInput.View() + "\n\n")
	s.WriteString(helpStyle.Render(m.keys.InputHelp()))
	
	return s.String()
}

func (m Model) viewStashMessage() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Stash Changes") + "\n\n")
	s.WriteString("Enter stash message (optional):\n\n")
	s.WriteString(m.textInput.View() + "\n\n")
	s.WriteString(helpStyle.Render(m.keys.InputHelp()))

	return s.String()
}

func (m Model) viewCustomTag() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Create Custom Tag") + "\n\n")
	s.WriteString("Enter semantic version (e.g., 1.0.0, 0.5.0-beta.1):\n\n")

	// Show prefix
	s.WriteString(m.prefix)
	s.WriteString(m.textInput.View() + "\n")

	// Show validation error if any
	if m.customTagError != "" {
		s.WriteString("\n" + errorStyle.Render("  ↑ "+m.customTagError) + "\n")
	}

	s.WriteString("\n" + helpStyle.Render(m.keys.InputHelp()))

	return s.String()
}