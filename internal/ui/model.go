package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// Model represents the root application state.
type Model struct {
	width  int
	height int

	// State
	rootPath string
	engine   string
	focused  int // 0: Browser, 1: Inspector

	// Components
	browser   FileBrowser
	inspector Inspector
	actionBar ActionBar
	compiler  *core.Compiler

	// Modals
	errorModal      ErrorModal
	newProjectModal NewProjectModal
	helpModal       HelpModal

	// Shared data
	projectFiles *core.ProjectFiles
}

// NewModel creates a new root model.
func NewModel(root, engine string) Model {
	if root == "" {
		root, _ = os.Getwd()
	}
	if engine == "" {
		engine = "pdflatex"
	}

	return Model{
		rootPath:        root,
		engine:          engine,
		browser:         NewFileBrowser(),
		inspector:       NewInspector(),
		actionBar:       NewActionBar(),
		compiler:        core.NewCompiler(),
		errorModal:      NewErrorModal(),
		newProjectModal: NewNewProjectModal(),
		helpModal:       NewHelpModal(),
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return m.scanDirCmd(m.rootPath)
}

// Update handles messages and updates state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

	case tea.KeyMsg:
		// Modal handling
		if m.helpModal.IsVisible() {
			if msg.Type == tea.KeyEscape || msg.String() == "?" {
				m.helpModal.Hide()
			}
			return m, nil
		}

		if m.errorModal.IsVisible() {
			switch msg.Type {
			case tea.KeyEscape:
				m.errorModal.Hide()
			case tea.KeyUp, tea.KeyPgUp:
				m.errorModal.ScrollUp()
			case tea.KeyDown, tea.KeyPgDown:
				m.errorModal.ScrollDown()
			}
			return m, nil
		}

		if m.newProjectModal.IsVisible() {
			cmd := m.newProjectModal.HandleKey(msg)
			return m, cmd
		}

		// Global keys
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.focused = (m.focused + 1) % 2
			m.browser.SetFocused(m.focused == 0)
			m.inspector.SetFocused(m.focused == 1)

		case "h":
			m.browser.ToggleShadow()

		case "n":
			m.newProjectModal.Show(m.rootPath)

		case "?":
			m.helpModal.Show()

		case "F5":
			if !m.compiler.IsBusy() {
				path, cat := m.browser.SelectedFile()
				if cat == core.CategorySource && path != "" {
					m.actionBar.SetBuildStatus(StatusBUILDING, 0, 0)
					cmds = append(cmds, m.compileCmd(path))
				}
			}

		case "F6":
			cmds = append(cmds, m.cleanCmd())

		case "F7":
			cmds = append(cmds, m.openPdfCmd())

		case "y":
			// Yank citekey to clipboard
			key := m.inspector.SelectedBibKey()
			if key != "" {
				YankToClipboard(key)
			}

		case "up", "k":
			if m.focused == 0 {
				m.browser.MoveUp()
				path, cat := m.browser.SelectedFile()
				m.inspector.SetFile(path, cat)
			} else {
				m.inspector.MoveBibUp()
			}

		case "down", "j":
			if m.focused == 0 {
				m.browser.MoveDown()
				path, cat := m.browser.SelectedFile()
				m.inspector.SetFile(path, cat)
			} else {
				m.inspector.MoveBibDown()
			}
		}

	case ScannedMsg:
		m.projectFiles = msg.Files
		m.browser.SetFiles(msg.Files)
		m.actionBar.SetProjectRoot(msg.Files.Root)
		path, cat := m.browser.SelectedFile()
		m.inspector.SetFile(path, cat)

	case BuildFinishedMsg:
		if msg.Err != nil {
			m.actionBar.SetBuildStatus(StatusFAILED, msg.Result.Duration, 1)
			// Trigger log parse
			cmds = append(cmds, m.parseLogCmd(msg.Result.LogPath))
		} else if !msg.Result.Success {
			// Trigger log parse
			cmds = append(cmds, m.parseLogCmd(msg.Result.LogPath))
		} else {
			m.actionBar.SetBuildStatus(StatusSUCCESS, msg.Result.Duration, 0)
			// Re-scan to see new PDF/aux files
			cmds = append(cmds, m.scanDirCmd(m.rootPath))
		}

	case LogParsedMsg:
		if msg.Err == nil && len(msg.Entries) > 0 {
			m.actionBar.SetBuildStatus(StatusFAILED, 0, core.ErrorCount(msg.Entries))
			if core.ErrorCount(msg.Entries) > 0 {
				m.errorModal.Show(msg.Entries)
			}
		}

	case CleanedMsg:
		if msg.Err == nil {
			cmds = append(cmds, m.scanDirCmd(m.rootPath))
		}

	case ProjectCreatedMsg:
		if msg.Err == nil {
			m.rootPath = msg.Path
			m.newProjectModal.Hide()
			cmds = append(cmds, m.scanDirCmd(m.rootPath))
		}

	case ErrorMsg:
		// General error handling
	}

	return m, tea.Batch(cmds...)
}

// View renders the application.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Main panes
	browserView := m.browser.View()
	inspectorView := m.inspector.View()
	actionBarView := m.actionBar.View()

	main := lipgloss.JoinHorizontal(lipgloss.Top, browserView, inspectorView)
	app := lipgloss.JoinVertical(lipgloss.Left, main, actionBarView)

	// Layer modals
	if m.helpModal.IsVisible() {
		return m.helpModal.View(m.width, m.height)
	}
	if m.newProjectModal.IsVisible() {
		return m.newProjectModal.View(m.width, m.height)
	}
	if m.errorModal.IsVisible() {
		return m.errorModal.View(m.width, m.height)
	}

	return app
}

// ── Commands ──

func (m *Model) updateLayout() {
	footerHeight := 1
	mainHeight := m.height - footerHeight
	browserWidth := int(float64(m.width) * 0.35)
	inspectorWidth := m.width - browserWidth

	m.browser.SetSize(browserWidth, mainHeight)
	m.inspector.SetSize(inspectorWidth, mainHeight)
	m.actionBar.SetWidth(m.width)
}
