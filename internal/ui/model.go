package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// Model represents the root application state.
type Model struct {
	width  int
	height int

	// State
	rootPath  string
	engine    string
	focused   int // 0: Browser, 1: Inspector
	filtering bool

	// Components
	browser     FileBrowser
	inspector   Inspector
	actionBar   ActionBar
	compiler    *core.Compiler
	filterInput textinput.Model
	watcher     *core.Watcher

	// Modals
	errorModal      ErrorModal
	newProjectModal NewProjectModal
	helpModal       HelpModal
	searchModal     SearchModal
	diffModal       DiffModal

	// Shared data
	projectFiles     *core.ProjectFiles
	globalBibEntries []core.BibEntry
}

// NewModel creates a new root model.
func NewModel(root, engine string) Model {
	if root == "" {
		root, _ = os.Getwd()
	}
	if engine == "" {
		engine = "pdflatex"
	}

	ti := textinput.New()
	ti.Placeholder = "Filter files... (Esc to cancel, Enter to accept)"
	ti.Prompt = " / "
	ti.CharLimit = 50

	w, _ := core.NewWatcher(root)

	return Model{
		rootPath:        root,
		engine:          engine,
		browser:         NewFileBrowser(),
		inspector:       NewInspector(),
		actionBar:       NewActionBar(),
		compiler:        core.NewCompiler(),
		filterInput:     ti,
		watcher:         w,
		errorModal:      NewErrorModal(),
		newProjectModal: NewNewProjectModal(),
		helpModal:       NewHelpModal(),
		diffModal:       NewDiffModal(),
		searchModal:     NewSearchModal(),
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.scanDirCmd(m.rootPath),
		m.listenForFileEventCmd(),
		m.loadGlobalBibCmd(),
	)
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
				m.errorModal.MoveUp()
			case tea.KeyDown, tea.KeyPgDown:
				m.errorModal.MoveDown()
			case tea.KeyEnter:
				// Jump to line
				if entry := m.errorModal.SelectedEntry(); entry != nil && entry.File != "" {
					return m, m.openEditorCmd(entry.File, entry.Line)
				}
			}
			return m, nil
		}

		if m.newProjectModal.IsVisible() {
			cmd := m.newProjectModal.HandleKey(msg)
			return m, cmd
		}

		if m.searchModal.IsVisible() {
			cmd := m.searchModal.HandleKey(msg)
			return m, cmd
		}

		if m.filtering {
			switch msg.Type {
			case tea.KeyEscape, tea.KeyEnter:
				m.filtering = false
				m.filterInput.Blur()
				return m, nil
			}

			var cmd tea.Cmd
			m.filterInput, cmd = m.filterInput.Update(msg)

			// Let browser know filter changed
			m.browser.SetFilter(m.filterInput.Value())

			// Also update inspector to match new browser selection
			path, cat := m.browser.SelectedFile()
			m.inspector.SetFile(path, cat)

			return m, cmd
		}

		// Global keys
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "/":
			if m.focused == 0 && !m.filtering {
				m.filtering = true
				m.filterInput.Focus()
				return m, textinput.Blink
			}

		case "tab":
			m.focused = (m.focused + 1) % 2
			m.browser.SetFocused(m.focused == 0)
			m.inspector.SetFocused(m.focused == 1)

		case "h":
			m.browser.ToggleShadow()

		case "n":
			m.newProjectModal.Show(m.rootPath)

		case "d":
			path, cat := m.browser.SelectedFile()
			if cat == core.CategorySource && path != "" {
				// Prepare list of files and tags
				var allTex []string
				if m.projectFiles != nil {
					for _, f := range m.projectFiles.Source {
						allTex = append(allTex, f.Path)
					}
				}
				cmds = append(cmds, m.listTagsCmd(path, allTex))
			}

		case "?":

			m.helpModal.Show()

		case "s":
			if !m.filtering {
				m.searchModal.Show(m.globalBibEntries)
				return m, nil
			}

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

		case "enter":
			if m.focused == 0 {
				path, cat := m.browser.SelectedFile()
				if path != "" && cat != core.CategoryOutput && cat != core.CategoryAssets {
					// Open in editor (only source, data, aux)
					return m, m.openEditorCmd(path, 0)
				}
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

		if m.diffModal.IsVisible() {
			cmd := m.diffModal.HandleKey(msg)
			if cmd != nil {
				return m, cmd
			}
			return m, nil
		}

	case ScannedMsg:
		m.projectFiles = msg.Files
		m.browser.SetFiles(msg.Files)
		m.actionBar.SetProjectRoot(msg.Files.Root)
		path, cat := m.browser.SelectedFile()
		m.inspector.SetFile(path, cat)

	case RunDiffMsg:
		m.actionBar.SetBuildStatus(StatusBUILDING, 0, 0)
		cmds = append(cmds, m.diffCmd(msg))

	case TagsListedMsg:
		m.diffModal.Show(msg.SelectedPath, msg.Tags, msg.AllFiles)

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

	case EditorClosedMsg:
		if msg.Err != nil {
			// Could show an error indicator, but usually it's fine
		}

	case FileEventMsg:
		// Queue a rescan and immediately re-listen
		cmds = append(cmds, m.scanDirCmd(m.rootPath), m.listenForFileEventCmd())

	case GlobalBibLoadedMsg:
		if msg.Err == nil {
			m.globalBibEntries = msg.Entries
		}

	case ErrorMsg:
		// General error handling
	}

	if m.filtering {
		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
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

	// Action/Filter Bar
	var bottomBarView string
	if m.filtering {
		inputView := m.filterInput.View()
		bottomBarView = lipgloss.NewStyle().
			Width(m.width).
			Padding(1, 2).
			Background(lipgloss.Color("0")).
			Foreground(lipgloss.Color("15")).
			Render(inputView)
	} else {
		bottomBarView = m.actionBar.View()
	}

	main := lipgloss.JoinHorizontal(lipgloss.Top, browserView, inspectorView)
	app := lipgloss.JoinVertical(lipgloss.Left, main, bottomBarView)

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
	if m.searchModal.IsVisible() {
		return m.searchModal.View(m.width, m.height)
	}
	if m.diffModal.IsVisible() {
		return m.diffModal.View(m.width, m.height)
	}

	return app
}

// GlobalBibLoadedMsg is emitted when the global bibliography is loaded.
type GlobalBibLoadedMsg struct {
	Entries []core.BibEntry
	Err     error
}

// ── Commands ──

func (m Model) loadGlobalBibCmd() tea.Cmd {
	return func() tea.Msg {
		config := core.LoadGlobalConfig()
		if config.GlobalBibPath == "" {
			return GlobalBibLoadedMsg{Err: fmt.Errorf("no global bib path defined")}
		}

		entries, err := core.BibMetadata(config.GlobalBibPath)
		return GlobalBibLoadedMsg{Entries: entries, Err: err}
	}
}

func (m Model) listTagsCmd(selectedPath string, allFiles []string) tea.Cmd {
	return func() tea.Msg {
		tags, _ := core.ListGitTags() // ignore error, just show empty if not git
		return TagsListedMsg{
			SelectedPath: selectedPath,
			Tags:         tags,
			AllFiles:     allFiles,
		}
	}
}

func (m Model) diffCmd(msg RunDiffMsg) tea.Cmd {
	return func() tea.Msg {
		res, err := m.compiler.Diff(msg.OldPath, msg.OldContent, msg.NewPath, m.rootPath, m.engine)
		return BuildFinishedMsg{Result: res, Err: err}
	}
}

func (m *Model) updateLayout() {
	footerHeight := 1
	mainHeight := m.height - footerHeight
	browserWidth := int(float64(m.width) * 0.35)
	inspectorWidth := m.width - browserWidth

	m.browser.SetSize(browserWidth, mainHeight)
	m.inspector.SetSize(inspectorWidth, mainHeight)
	m.actionBar.SetWidth(m.width)
}
