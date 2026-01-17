package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/adnsv/go-utils/git"
	"github.com/adnsv/rtag/internal/version"
	"github.com/blang/semver/v4"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) initDirtyChoiceList() {
	items := []list.Item{
		dirtyChoiceItem{choice: "commit", title: "Commit", description: "Stage and commit all current changes"},
		dirtyChoiceItem{choice: "stash", title: "Stash", description: "Stash changes to create tag on clean state"},
		dirtyChoiceItem{choice: "proceed", title: "Proceed anyway", description: "Create tag on dirty repository (not recommended)"},
	}
	m.dirtyChoiceList = list.New(items, newItemDelegate(), m.width-4, min(len(items)*3+5, m.height-10))
	m.dirtyChoiceList.SetShowTitle(false)
	m.dirtyChoiceList.SetShowStatusBar(false)
	m.dirtyChoiceList.SetFilteringEnabled(false)
	m.dirtyChoiceList.SetShowHelp(false)
}

func (m *Model) initUndoScopeList() {
	items := []list.Item{
		undoScopeItem{scope: "local", title: "Local only", description: "Delete tag from local repository"},
		undoScopeItem{scope: "remote", title: "Remote only", description: "Delete tag from remote origin"},
		undoScopeItem{scope: "both", title: "Both", description: "Delete tag from local and remote"},
	}
	m.undoScopeList = list.New(items, newItemDelegate(), m.width-4, min(len(items)*3+5, m.height-10))
	m.undoScopeList.SetShowTitle(false)
	m.undoScopeList.SetShowStatusBar(false)
	m.undoScopeList.SetFilteringEnabled(false)
	m.undoScopeList.SetShowHelp(false)
}

func (m *Model) initConfirmList(yesDesc, noDesc string) {
	items := []list.Item{
		confirmItem{choice: true, title: "Yes", description: yesDesc},
		confirmItem{choice: false, title: "No", description: noDesc},
	}
	m.confirmList = list.New(items, newItemDelegate(), m.width-4, min(len(items)*3+4, m.height-10))
	m.confirmList.SetShowTitle(false)
	m.confirmList.SetShowStatusBar(false)
	m.confirmList.SetFilteringEnabled(false)
	m.confirmList.SetShowHelp(false)
}

func (m *Model) initTagConfirmList() {
	items := []list.Item{
		tagConfirmItem{action: "yes", title: "Yes", description: "Create this tag"},
		tagConfirmItem{action: "edit", title: "Edit comment", description: "Customize the tag message"},
		tagConfirmItem{action: "no", title: "No", description: "Cancel"},
	}
	m.confirmList = list.New(items, newItemDelegate(), m.width-4, min(len(items)*3+4, m.height-10))
	m.confirmList.SetShowTitle(false)
	m.confirmList.SetShowStatusBar(false)
	m.confirmList.SetFilteringEnabled(false)
	m.confirmList.SetShowHelp(false)
}

func (m *Model) initActionList() {
	// Create action list with version options + undo
	items := make([]list.Item, 0, len(m.versions)+1)
	for _, v := range m.versions {
		items = append(items, actionItem{
			action: v,
			prefix: m.prefix,
		})
	}

	// Add undo option if there's a tag to delete
	if m.stats != nil && m.stats.Description.Tag != "" {
		items = append(items, actionItem{
			isUndo:  true,
			undoTag: m.stats.Description.Tag,
		})
	}

	// Leave room for repo info header (~8 lines) and help text
	listHeight := min(len(items)*3+4, m.height-12)
	m.actionList = list.New(items, newItemDelegate(), m.width-4, listHeight)
	m.actionList.SetShowTitle(false)
	m.actionList.SetShowStatusBar(false)
	m.actionList.SetFilteringEnabled(false)
	m.actionList.SetShowHelp(false)
}

func (m *Model) initFirstTagActionList() {
	// Create action list with first tag options: v0.1.0 and Custom
	items := []list.Item{
		actionItem{
			action: version.Action{
				Desc: "Initial version|recommended starting version",
				Ver:  semver.Version{Major: 0, Minor: 1, Patch: 0},
			},
			prefix: m.prefix,
		},
		actionItem{
			isCustom: true,
			prefix:   m.prefix,
		},
	}

	listHeight := min(len(items)*3+4, m.height-12)
	m.actionList = list.New(items, newItemDelegate(), m.width-4, listHeight)
	m.actionList.SetShowTitle(false)
	m.actionList.SetShowStatusBar(false)
	m.actionList.SetFilteringEnabled(false)
	m.actionList.SetShowHelp(false)
}

// Update handles all messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update list dimensions if they exist
		if m.state == StateSelectAction {
			m.actionList.SetWidth(m.width - 4)
			m.actionList.SetHeight(min(len(m.actionList.Items())*3+5, m.height-10))
		}
		if m.state == StateSelectPRType {
			m.prTypeList.SetWidth(m.width - 4)
			m.prTypeList.SetHeight(min(len(m.prTypeList.Items())*3+5, m.height-10))
		}
		return m, nil
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			// Go back to previous state
			return m.handleBack()
		}
	}
	
	// Handle state-specific updates
	switch m.state {
	case StateLoading:
		return m.updateLoading(msg)
	case StateDirtyRepoChoice:
		return m.updateDirtyRepoChoice(msg)
	case StateCommitMessage:
		return m.updateCommitMessage(msg)
	case StateStashMessage:
		return m.updateStashMessage(msg)
	case StateExecutingGitCommand:
		return m.updateExecutingGitCommand(msg)
	case StateSelectAction:
		return m.updateSelectAction(msg)
	case StateSelectPRType:
		return m.updateSelectPRType(msg)
	case StateConfirmTag:
		return m.updateConfirmTag(msg)
	case StateEditTagComment:
		return m.updateEditTagComment(msg)
	case StateExecutingTag:
		return m.updateExecuting(msg)
	case StateConfirmPush:
		return m.updateConfirmPush(msg)
	case StateExecutingPush:
		return m.updateExecuting(msg)
	case StateError:
		return m.updateError(msg)
	case StateUndoPreview:
		return m.updateUndoPreview(msg)
	case StateConfirmUndo:
		return m.updateConfirmUndo(msg)
	case StateExecutingUndo:
		return m.updateExecuting(msg)
	case StateCustomTag:
		return m.updateCustomTag(msg)
	case StateDone:
		return m.updateDone(msg)
	}

	return m, nil
}

func (m Model) updateDone(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateLoading(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgRepoStats:
		if msg.err != nil {
			if errors.Is(msg.err, git.ErrNoTags) {
				// Handle first tag scenario - show action selection with first tag options
				m.stats = nil
				m.workDir = msg.workDir
				m.autoPrefix = "v"
				if m.opts.Prefix == "" || m.opts.Prefix == "AUTO" {
					m.prefix = m.autoPrefix
				} else {
					m.prefix = m.opts.Prefix
				}
				// Create action list with first tag options
				m.initFirstTagActionList()
				m.state = StateSelectAction
				return m, nil
			} else {
				m.err = fmt.Errorf("failed to obtain git stats: %w", msg.err)
				m.state = StateError
			}
			return m, nil
		}

		m.workDir = msg.workDir
		m.stats = msg.stats

		// Determine auto prefix
		if m.stats.Description.Tag != "" {
			vi, err := git.ParseVersion(m.stats.Description)
			if err == nil && vi != nil {
				tagLen := len(vi.Semantic.String())
				fullLen := len(m.stats.Description.Tag)
				if tagLen < fullLen {
					m.autoPrefix = m.stats.Description.Tag[:fullLen-tagLen]
				}
			}
		}
		if m.autoPrefix == "" {
			m.autoPrefix = "v"
		}

		// Set prefix
		if m.opts.Prefix == "" || m.opts.Prefix == "AUTO" {
			m.prefix = m.autoPrefix
		} else {
			m.prefix = m.opts.Prefix
		}

		// Check for dirty repo first (before parsing versions)
		if m.stats.Dirty {
			// Parse versions in background, then go to dirty choice
			return m, parseVersions(m.stats, m.prefix)
		}

		// Normal flow - parse versions
		return m, parseVersions(m.stats, m.prefix)

	case msgVersionsParsed:
		m.versions = msg.versions

		// Now handle the transition based on current state
		// Check if we need dirty repo choice
		if m.stats != nil && m.stats.Dirty {
			m.initDirtyChoiceList()
			m.state = StateDirtyRepoChoice
			return m, nil
		}

		// Normal flow - go to action selection
		if len(m.versions) == 0 {
			// First tag scenario
			m.initFirstTagActionList()
			m.state = StateSelectAction
			return m, nil
		}

		// Create action list with version options + undo
		m.initActionList()
		m.state = StateSelectAction
		return m, nil

	default:
		// Update spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) updateSelectAction(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.actionList.SelectedItem().(actionItem)
			if !ok {
				return m, nil
			}

			// Handle undo action
			if selected.isUndo {
				m.currentTag = selected.undoTag
				m.state = StateUndoPreview
				return m, loadTagList()
			}

			// Handle custom tag action
			if selected.isCustom {
				ti := textinput.New()
				ti.Placeholder = "1.0.0"
				ti.Focus()
				ti.CharLimit = 50
				ti.Width = 40
				m.textInput = ti
				m.customTagError = ""
				m.state = StateCustomTag
				return m, textinput.Blink
			}

			m.selectedAction = &selected.action

			if selected.action.ShowPRChoice {
				// Need to select pre-release type
				items := []list.Item{
					prTypeItem{prType: "alpha", prefix: m.prefix, version: selected.action.Ver},
					prTypeItem{prType: "beta", prefix: m.prefix, version: selected.action.Ver},
					prTypeItem{prType: "rc", prefix: m.prefix, version: selected.action.Ver},
					prTypeItem{prType: "release", prefix: m.prefix, version: selected.action.Ver},
				}

				m.prTypeList = list.New(items, newItemDelegate(), m.width-4, min(len(items)*3+5, m.height-10))
				m.prTypeList.SetShowTitle(false)
				m.prTypeList.SetShowStatusBar(false)
				m.prTypeList.SetFilteringEnabled(false)
				m.prTypeList.SetShowHelp(false)
				m.state = StateSelectPRType
			} else {
				// Direct to confirmation
				m.newTag = m.prefix + selected.action.Ver.String()
				m.tagComment = generateTagComment(m.newTag)
				m.initTagConfirmList()
				m.state = StateConfirmTag
			}
			return m, nil
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.actionList, cmd = m.actionList.Update(msg)
	return m, cmd
}

func (m Model) updateSelectPRType(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.prTypeList.SelectedItem().(prTypeItem)
			if ok {
				m.selectedPRType = selected.prType
				
				// Create the new version based on selection
				if m.selectedAction != nil {
					var newVer = m.selectedAction.Ver
					if selected.prType != "release" {
						newVer.Pre = version.MakePR(selected.prType, 1)
					} else {
						newVer.Pre = newVer.Pre[:0]
					}
					
					m.newTag = m.prefix + newVer.String()
					m.tagComment = generateTagComment(m.newTag)
					m.initTagConfirmList()
					m.state = StateConfirmTag
				}
				return m, nil
			}
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.prTypeList, cmd = m.prTypeList.Update(msg)
	return m, cmd
}

func (m Model) updateConfirmTag(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.confirmList.SelectedItem().(tagConfirmItem)
			if ok {
				switch selected.action {
				case "yes":
					m.state = StateExecutingTag
					m.lastCommand = fmt.Sprintf("git tag -a %s -m \"%s\"", m.newTag, m.tagComment)
					return m, gitTag(m.newTag, m.tagComment)
				case "edit":
					ti := textinput.New()
					ti.SetValue(m.tagComment)
					ti.Focus()
					ti.CharLimit = 200
					ti.Width = 60
					m.textInput = ti
					m.state = StateEditTagComment
					return m, textinput.Blink
				case "no":
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.confirmList, cmd = m.confirmList.Update(msg)
	return m, cmd
}

func (m Model) updateEditTagComment(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.tagComment = m.textInput.Value()
			m.initTagConfirmList()
			m.state = StateConfirmTag
			return m, nil
		case tea.KeyEsc:
			m.initTagConfirmList()
			m.state = StateConfirmTag
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateConfirmPush(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.confirmList.SelectedItem().(confirmItem)
			if ok {
				if selected.choice {
					m.state = StateExecutingPush
					m.lastCommand = fmt.Sprintf("git push origin %s", m.newTag)
					return m, gitPushTag(m.newTag)
				}
				m.state = StateDone
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.confirmList, cmd = m.confirmList.Update(msg)
	return m, cmd
}

func (m Model) updateExecuting(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgCommandComplete:
		m.commandOutput = msg.output
		
		if msg.err != nil {
			m.err = msg.err
			m.state = StateError
			return m, nil
		}
		
		// Success - determine next state
		if m.state == StateExecutingTag {
			m.initConfirmList("Push to remote", "Keep local only")
			m.state = StateConfirmPush
		} else if m.state == StateExecutingPush {
			m.tagPushed = true
			m.state = StateDone
		} else if m.state == StateExecutingUndo {
			// Check if we need to delete remote too
			if m.undoScope == "both" && m.currentTag != "" && m.lastCommand != "" && !strings.Contains(m.lastCommand, "push") {
				// Just finished local delete, now do remote
				m.lastCommand = fmt.Sprintf("git push --delete origin %s", m.currentTag)
				return m, gitDeleteRemoteTag(m.currentTag)
			}
			m.state = StateDone
		}
		return m, nil
		
	default:
		// Update spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) updateError(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) handleBack() (Model, tea.Cmd) {
	switch m.state {
	case StateSelectAction:
		// No previous screen - quit
		return m, tea.Quit
	case StateSelectPRType:
		m.state = StateSelectAction
	case StateConfirmTag:
		if m.selectedAction == nil {
			// First tag or custom tag scenario - go back to action selection
			if m.stats == nil {
				m.initFirstTagActionList()
			} else {
				m.initActionList()
			}
			m.state = StateSelectAction
		} else if m.selectedAction.ShowPRChoice {
			m.state = StateSelectPRType
		} else {
			m.state = StateSelectAction
		}
	case StateEditTagComment:
		m.initTagConfirmList()
		m.state = StateConfirmTag
	case StateConfirmPush:
		// Can't go back from here - tag already created
	case StateError:
		return m, tea.Quit
	case StateDirtyRepoChoice:
		// No previous screen - quit
		return m, tea.Quit
	case StateCommitMessage, StateStashMessage:
		m.state = StateDirtyRepoChoice
	case StateUndoPreview:
		// Go back to action selection
		m.tagList = nil
		m.initActionList()
		m.state = StateSelectAction
	case StateConfirmUndo:
		m.state = StateUndoPreview
	case StateCustomTag:
		// Go back to action selection
		if m.stats == nil {
			m.initFirstTagActionList()
		} else {
			m.initActionList()
		}
		m.state = StateSelectAction
	}
	return m, nil
}

func (m Model) updateDirtyRepoChoice(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.dirtyChoiceList.SelectedItem().(dirtyChoiceItem)
			if ok {
				m.dirtyRepoChoice = selected.choice
				switch selected.choice {
				case "commit":
					ti := textinput.New()
					ti.Placeholder = "Enter commit message"
					ti.Focus()
					ti.CharLimit = 200
					ti.Width = 60
					m.textInput = ti
					m.state = StateCommitMessage
					return m, textinput.Blink
				case "stash":
					ti := textinput.New()
					ti.Placeholder = "Optional stash message (press Enter to skip)"
					ti.Focus()
					ti.CharLimit = 100
					ti.Width = 60
					m.textInput = ti
					m.state = StateStashMessage
					return m, textinput.Blink
				case "proceed":
					if len(m.versions) == 0 {
						// First tag scenario
						m.initFirstTagActionList()
					} else {
						m.initActionList()
					}
					m.state = StateSelectAction
					return m, nil
				}
			}
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.dirtyChoiceList, cmd = m.dirtyChoiceList.Update(msg)
	return m, cmd
}

func (m Model) updateCommitMessage(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() == "" {
				// Don't allow empty commit message
				return m, nil
			}
			// Execute git add and commit
			m.state = StateExecutingGitCommand
			m.lastCommand = "git add -A && git commit"
			return m, gitAddAndCommit(m.textInput.Value())
		case tea.KeyEsc:
			// Go back to dirty repo choice
			m.state = StateDirtyRepoChoice
			return m, nil
		}
	}
	
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateStashMessage(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Execute git stash
			m.state = StateExecutingGitCommand
			if m.textInput.Value() == "" {
				m.lastCommand = "git stash"
				return m, gitStash("")
			} else {
				m.lastCommand = "git stash save"
				return m, gitStash(m.textInput.Value())
			}
		case tea.KeyEsc:
			// Go back to dirty repo choice
			m.state = StateDirtyRepoChoice
			return m, nil
		}
	}
	
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) updateExecutingGitCommand(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgCommandComplete:
		m.commandOutput = msg.output
		
		if msg.err != nil {
			m.err = msg.err
			m.state = StateError
			return m, nil
		}
		
		// Success - reload stats and go back to show stats
		m.state = StateLoading
		return m, loadRepoStats
		
	default:
		// Update spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) updateUndoPreview(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgTagList:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateError
			return m, nil
		}
		m.tagList = msg.tags
		// Initialize the undo scope list
		m.initUndoScopeList()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.undoScopeList.SelectedItem().(undoScopeItem)
			if ok {
				m.undoScope = selected.scope
				m.initConfirmList("Delete tag", "Cancel")
				m.state = StateConfirmUndo
				return m, nil
			}
		}
	}

	// Update spinner while loading, or update the list
	if m.tagList == nil {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.undoScopeList, cmd = m.undoScopeList.Update(msg)
	return m, cmd
}

func (m Model) updateCustomTag(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Validate the input
			input := m.textInput.Value()
			if input == "" {
				m.customTagError = "Version cannot be empty"
				return m, nil
			}

			// Parse with semver
			_, err := semver.Parse(input)
			if err != nil {
				m.customTagError = fmt.Sprintf("Invalid: %s", err.Error())
				return m, nil
			}

			// Valid - proceed to confirm
			m.newTag = m.prefix + input
			m.tagComment = generateTagComment(m.newTag)
			m.initTagConfirmList()
			m.state = StateConfirmTag
			return m, nil
		case tea.KeyEsc:
			// Go back to action selection
			if m.stats == nil {
				// First tag scenario
				m.initFirstTagActionList()
			} else {
				m.initActionList()
			}
			m.state = StateSelectAction
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	// Clear error on any input change
	m.customTagError = ""
	return m, cmd
}

func (m Model) updateConfirmUndo(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.confirmList.SelectedItem().(confirmItem)
			if !ok {
				break
			}
			if !selected.choice {
				return m, tea.Quit
			}

			m.state = StateExecutingUndo

			// Use the tag stored in currentTag (set when undo was selected)
			if m.currentTag == "" {
				m.err = fmt.Errorf("no tag to delete")
				m.state = StateError
				return m, nil
			}

			// Execute based on scope
			switch m.undoScope {
			case "local":
				m.lastCommand = fmt.Sprintf("git tag -d %s", m.currentTag)
				return m, gitDeleteLocalTag(m.currentTag)
			case "remote":
				m.lastCommand = fmt.Sprintf("git push --delete origin %s", m.currentTag)
				return m, gitDeleteRemoteTag(m.currentTag)
			case "both":
				// For both, we'll do local first
				m.lastCommand = fmt.Sprintf("git tag -d %s", m.currentTag)
				return m, gitDeleteLocalTag(m.currentTag)
			}
		}
	}

	var cmd tea.Cmd
	m.confirmList, cmd = m.confirmList.Update(msg)
	return m, cmd
}

// List item types

// actionItem represents a version action in the list
type actionItem struct {
	action   version.Action
	prefix   string
	isUndo   bool   // if true, this is the undo action
	undoTag  string // tag to undo (only used if isUndo is true)
	isCustom bool   // if true, this is the custom version option
}

func (i actionItem) Title() string {
	if i.isUndo {
		return "Undo last tag"
	}
	if i.isCustom {
		return "Custom version"
	}
	desc := i.action.Desc
	// Split by pipe to get the main description
	if idx := strings.Index(desc, "|"); idx >= 0 {
		desc = desc[:idx]
	}
	return fmt.Sprintf("%s %s", desc, i.prefix+i.action.Ver.String())
}

func (i actionItem) Description() string {
	if i.isUndo {
		return "delete " + i.undoTag
	}
	if i.isCustom {
		return "enter your own semantic version"
	}
	desc := i.action.Desc
	if idx := strings.Index(desc, "|"); idx >= 0 {
		return desc[idx+1:]
	}
	return ""
}

func (i actionItem) FilterValue() string { return i.Title() }

type prTypeItem struct {
	prType  string
	prefix  string
	version semver.Version
}

func (i prTypeItem) Title() string {
	return strings.ToUpper(i.prType[:1]) + i.prType[1:]
}

func (i prTypeItem) Description() string {
	var displayVer semver.Version
	if i.prType == "release" {
		displayVer = i.version
		displayVer.Pre = nil
	} else {
		displayVer = i.version
		displayVer.Pre = version.MakePR(i.prType, 1)
	}
	return i.prefix + displayVer.String()
}

func (i prTypeItem) FilterValue() string { return i.prType }

type dirtyChoiceItem struct {
	choice      string
	title       string
	description string
}

func (i dirtyChoiceItem) Title() string       { return i.title }
func (i dirtyChoiceItem) Description() string { return i.description }
func (i dirtyChoiceItem) FilterValue() string { return i.choice }

type undoScopeItem struct {
	scope       string
	title       string
	description string
}

func (i undoScopeItem) Title() string       { return i.title }
func (i undoScopeItem) Description() string { return i.description }
func (i undoScopeItem) FilterValue() string { return i.scope }

type confirmItem struct {
	choice      bool
	title       string
	description string
}

func (i confirmItem) Title() string       { return i.title }
func (i confirmItem) Description() string { return i.description }
func (i confirmItem) FilterValue() string { return i.title }

type tagConfirmItem struct {
	action      string // "yes", "edit", "no"
	title       string
	description string
}

func (i tagConfirmItem) Title() string       { return i.title }
func (i tagConfirmItem) Description() string { return i.description }
func (i tagConfirmItem) FilterValue() string { return i.action }