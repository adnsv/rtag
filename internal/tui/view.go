package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")) // Bold white

	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	warningStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214"))

	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")) // Green

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	keyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(20)

	valueStyle = lipgloss.NewStyle()

	tagStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")) // Green

	commandStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

)

// newItemDelegate creates a default delegate with green selection colors
func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// Green selection colors
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("46")).
		BorderForeground(lipgloss.Color("46"))

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("243")).
		BorderForeground(lipgloss.Color("46"))

	return d
}

// View renders the current state
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}
	
	switch m.state {
	case StateLoading:
		return m.viewLoading()
	case StateDirtyRepoChoice:
		return m.viewDirtyRepoChoice()
	case StateCommitMessage:
		return m.viewCommitMessage()
	case StateStashMessage:
		return m.viewStashMessage()
	case StateExecutingGitCommand:
		return m.viewExecuting()
	case StateSelectAction:
		return m.viewSelectAction()
	case StateSelectPRType:
		return m.viewSelectPRType()
	case StateConfirmTag:
		return m.viewConfirmTag()
	case StateExecutingTag, StateExecutingPush, StateExecutingUndo:
		return m.viewExecuting()
	case StateConfirmPush:
		return m.viewConfirmPush()
	case StateDone:
		return m.viewDone()
	case StateError:
		return m.viewError()
	case StateUndoSelectScope:
		return m.viewUndoSelectScope()
	case StateConfirmUndo:
		return m.viewConfirmUndo()
	default:
		return "Unknown state"
	}
}

func (m Model) viewLoading() string {
	return fmt.Sprintf("\n %s Loading repository information...\n", m.spinner.View())
}

func (m Model) viewSelectAction() string {
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

		if m.stats.Dirty {
			s.WriteString(m.renderKeyValue("state", warningStyle.Render("dirty")))
		} else {
			s.WriteString(m.renderKeyValue("state", successStyle.Render("clean")))
		}

		if m.prefix != "" {
			s.WriteString(m.renderKeyValue("prefix", m.prefix))
		}
	} else {
		s.WriteString(infoStyle.Render("No existing tags found.") + "\n")
	}

	// Action selection
	s.WriteString("\n" + titleStyle.Render("Select version action") + "\n\n")
	s.WriteString(m.actionList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ListHelp()))

	return s.String()
}

func (m Model) viewSelectPRType() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Select Pre-release Type"),
		m.prTypeList.View(),
		helpStyle.Render(m.keys.ListHelpWithBack()),
	)
}

func (m Model) viewConfirmTag() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Ready to Tag") + "\n\n")
	s.WriteString(m.renderKeyValue("tag", tagStyle.Render(m.newTag)))
	s.WriteString(m.renderKeyValue("comment", infoStyle.Render("\""+m.tagComment+"\"")))

	s.WriteString("\n")
	s.WriteString(m.confirmList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ConfirmHelp()))

	return s.String()
}

func (m Model) viewConfirmPush() string {
	var s strings.Builder

	s.WriteString(successStyle.Render("✓ Tag created successfully!") + "\n\n")
	s.WriteString("Your local repository is now tagged as " + tagStyle.Render(m.newTag) + "\n\n")
	s.WriteString(titleStyle.Render("Push to remote?") + "\n\n")
	s.WriteString(m.confirmList.View())
	s.WriteString("\n" + helpStyle.Render(m.keys.ConfirmHelp()))

	return s.String()
}

func (m Model) viewExecuting() string {
	var s strings.Builder
	
	action := "Executing"
	if m.state == StateExecutingTag {
		action = "Creating tag"
	} else if m.state == StateExecutingPush {
		action = "Pushing to remote"
	} else if m.state == StateExecutingGitCommand {
		if strings.Contains(m.lastCommand, "commit") {
			action = "Committing changes"
		} else if strings.Contains(m.lastCommand, "stash") {
			action = "Stashing changes"
		}
	} else if m.state == StateExecutingUndo {
		action = "Deleting tag"
	}
	
	s.WriteString(fmt.Sprintf("\n %s %s...\n\n", m.spinner.View(), action))
	
	if m.lastCommand != "" {
		s.WriteString(commandStyle.Render(m.lastCommand) + "\n")
	}
	
	if m.commandOutput != "" {
		s.WriteString("\n" + m.commandOutput)
	}
	
	return s.String()
}

func (m Model) viewDone() string {
	var s strings.Builder

	s.WriteString("\n" + successStyle.Render("✓ Done!") + "\n\n")

	if m.newTag != "" {
		if m.tagPushed {
			s.WriteString("Successfully created and pushed tag: " + tagStyle.Render(m.newTag) + "\n")
		} else {
			s.WriteString("Successfully created tag: " + tagStyle.Render(m.newTag) + "\n")
			s.WriteString(infoStyle.Render("Tag was not pushed to remote.") + "\n")
		}
	}

	s.WriteString("\n" + helpStyle.Render(m.keys.DoneHelp()))

	return s.String()
}

func (m Model) viewError() string {
	var s strings.Builder
	
	s.WriteString("\n" + errorStyle.Render("Error: "+m.err.Error()) + "\n")
	
	if m.commandOutput != "" {
		s.WriteString("\nCommand output:\n" + m.commandOutput)
	}
	
	s.WriteString("\n" + helpStyle.Render(m.keys.DoneHelp()))

	return s.String()
}

func (m Model) renderKeyValue(key string, value string) string {
	return fmt.Sprintf("%s %s\n", keyStyle.Render(key+":"), valueStyle.Render(value))
}