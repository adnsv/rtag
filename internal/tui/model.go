package tui

import (
	"github.com/adnsv/go-utils/git"
	"github.com/adnsv/rtag/internal/version"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the current state of the application
type AppState int

const (
	StateLoading AppState = iota
	StateDirtyRepoChoice
	StateCommitMessage
	StateStashMessage
	StateExecutingGitCommand
	StateSelectAction
	StateSelectPRType
	StateConfirmTag
	StateEditTagComment
	StateExecutingTag
	StateConfirmPush
	StateExecutingPush
	StateUndo
	StateUndoSelectScope
	StateConfirmUndo
	StateExecutingUndo
	StateDone
	StateError
)

// Model represents the application state
type Model struct {
	// Application state
	state       AppState
	undo        bool
	opts        Options
	
	// Repository information
	workDir     string
	stats       *git.Stats
	currentTag  string
	versions    []version.Action
	prefix      string
	autoPrefix  string
	
	// UI components
	actionList      list.Model
	prTypeList      list.Model
	dirtyChoiceList list.Model
	undoScopeList   list.Model
	confirmList     list.Model
	spinner         spinner.Model
	viewport        viewport.Model
	textInput       textinput.Model
	keys            KeyMap

	// User selections
	selectedAction   *version.Action
	selectedPRType   string
	undoScope        string // "local", "remote", "both"
	dirtyRepoChoice  string // "commit", "stash", "proceed"
	
	// Tag creation state
	newTag           string
	tagComment       string
	tagPushed        bool
	
	// Command execution
	lastCommand      string
	commandOutput    string
	
	// Error handling
	err              error
	width            int
	height           int
}

// Options represents command-line options
type Options struct {
	Prefix     string
	AllowDirty bool
	Undo       bool
}

// Message types
type msgRepoStats struct {
	workDir string
	stats   *git.Stats
	err     error
}

type msgVersionsParsed struct {
	versions []version.Action
	prefix   string
}

type msgCommandComplete struct {
	output string
	err    error
}

type msgCommandOutput struct {
	line string
}

// NewModel creates a new TUI model
func NewModel(opts Options) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		state:   StateLoading,
		opts:    opts,
		spinner: s,
		keys:    DefaultKeyMap(),
		undo:    opts.Undo,
		width:   80,  // Default width
		height:  24,  // Default height
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadRepoStats,
	)
}