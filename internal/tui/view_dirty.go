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

func (m Model) viewUndoSelectScope() string {
	var s strings.Builder

	// Repo info header
	s.WriteString(titleStyle.Render("Repository Info") + "\n")

	if m.stats != nil {
		s.WriteString(m.renderKeyValue("branch", m.stats.Branch))

		if m.stats.Description.Tag != "" {
			s.WriteString(m.renderKeyValue("tag to delete", tagStyle.Render(m.stats.Description.Tag)))
		}

		if m.stats.Dirty {
			s.WriteString(m.renderKeyValue("state", warningStyle.Render("dirty")))
		} else {
			s.WriteString(m.renderKeyValue("state", successStyle.Render("clean")))
		}
	}

	s.WriteString("\n" + titleStyle.Render("Select deletion scope") + "\n\n")
	s.WriteString(m.undoScopeList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ListHelp()))

	return s.String()
}

func (m Model) viewConfirmUndo() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Confirm Tag Deletion") + "\n\n")

	tag := ""
	if m.stats != nil && m.stats.Description.Tag != "" {
		tag = m.stats.Description.Tag
	}

	s.WriteString(m.renderKeyValue("tag", tagStyle.Render(tag)))
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