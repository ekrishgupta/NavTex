package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/latex"
)

type DiffTargetType int

const (
	TargetLastCommit DiffTargetType = iota
	TargetTag
	TargetFile
)

type DiffModal struct {
	visible     bool
	cursor      int
	targets     []diffTarget
	selectedTex string
	tags        []string
	files       []string
	state       diffModalState
}

type diffModalState int

const (
	stateTargetSelect diffModalState = iota
	stateTagSelect
	stateFileSelect
)

type diffTarget struct {
	name string
	typ  DiffTargetType
}

func NewDiffModal() DiffModal {
	return DiffModal{
		targets: []diffTarget{
			{"Last Git Commit", TargetLastCommit},
			{"Git Tag...", TargetTag},
			{"Other File...", TargetFile},
		},
	}
}

func (dm *DiffModal) Show(selectedTex string, tags []string, allFiles []string) {
	dm.visible = true
	dm.selectedTex = selectedTex
	dm.tags = tags
	dm.files = nil
	// Only include .tex files for comparison, excluding the selected one
	for _, f := range allFiles {
		if strings.HasSuffix(f, ".tex") && f != selectedTex {
			dm.files = append(dm.files, f)
		}
	}
	dm.cursor = 0
	dm.state = stateTargetSelect
}

func (dm *DiffModal) Hide() {
	dm.visible = false
}

func (dm *DiffModal) IsVisible() bool {
	return dm.visible
}

type RunDiffMsg struct {
	OldContent string
	OldPath    string
	NewPath    string
}

func (dm *DiffModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if !dm.visible {
		return nil
	}

	switch msg.Type {
	case tea.KeyEscape:
		if dm.state != stateTargetSelect {
			dm.state = stateTargetSelect
			dm.cursor = 0
			return nil
		}
		dm.Hide()
		return nil

	case tea.KeyUp, tea.KeyLeft:
		dm.cursor--
		if dm.cursor < 0 {
			dm.cursor = dm.maxCursor()
		}

	case tea.KeyDown, tea.KeyRight:
		dm.cursor++
		if dm.cursor > dm.maxCursor() {
			dm.cursor = 0
		}

	case tea.KeyEnter:
		return dm.selectCurrent()
	}

	return nil
}

func (dm *DiffModal) maxCursor() int {
	switch dm.state {
	case stateTargetSelect:
		return len(dm.targets) - 1
	case stateTagSelect:
		return len(dm.tags) - 1
	case stateFileSelect:
		return len(dm.files) - 1
	}
	return 0
}

func (dm *DiffModal) selectCurrent() tea.Cmd {
	switch dm.state {
	case stateTargetSelect:
		target := dm.targets[dm.cursor]
		switch target.typ {
		case TargetLastCommit:
			dm.Hide()
			return func() tea.Msg {
				content, err := latex.GetGitLastCommitContent(dm.selectedTex)
				if err != nil {
					return ErrorMsg{Err: err}
				}
				return RunDiffMsg{OldContent: content, NewPath: dm.selectedTex}
			}
		case TargetTag:
			if len(dm.tags) == 0 {
				dm.Hide()
				return func() tea.Msg { return ErrorMsg{Err: fmt.Errorf("no tags found")} }
			}
			dm.state = stateTagSelect
			dm.cursor = 0
		case TargetFile:
			if len(dm.files) == 0 {
				dm.Hide()
				return func() tea.Msg { return ErrorMsg{Err: fmt.Errorf("no other tex files found")} }
			}
			dm.state = stateFileSelect
			dm.cursor = 0
		}

	case stateTagSelect:
		tag := dm.tags[dm.cursor]
		dm.Hide()
		return func() tea.Msg {
			content, err := latex.GetGitVersionContent(dm.selectedTex, tag)
			if err != nil {
				return ErrorMsg{Err: err}
			}
			return RunDiffMsg{OldContent: content, NewPath: dm.selectedTex}
		}

	case stateFileSelect:
		otherFile := dm.files[dm.cursor]
		dm.Hide()
		return func() tea.Msg {
			// Read other file content
			content, err := latex.GetGitVersionContent(otherFile, "HEAD") // Fallback to current if not in git
			if err != nil {
				// Try reading from disk if git fails
				// But latexdiff can take a file path directly?
				// Actually our Compiler.Diff expects oldContent string for now.
				// Let's just use the path if we want to support non-git files.
				return RunDiffMsg{OldPath: otherFile, NewPath: dm.selectedTex}
			}
			return RunDiffMsg{OldContent: content, NewPath: dm.selectedTex}
		}
	}
	return nil
}

func (dm DiffModal) View(termWidth, termHeight int) string {
	if !dm.visible {
		return ""
	}

	var title string
	var items []string

	switch dm.state {
	case stateTargetSelect:
		title = "Compare against..."
		for i, t := range dm.targets {
			if i == dm.cursor {
				items = append(items, FileItemSelected.Render("> "+t.name))
			} else {
				items = append(items, FileItem.Render("  "+t.name))
			}
		}
	case stateTagSelect:
		title = "Select Tag"
		for i, t := range dm.tags {
			if i == dm.cursor {
				items = append(items, FileItemSelected.Render("> "+t))
			} else {
				items = append(items, FileItem.Render("  "+t))
			}
		}
	case stateFileSelect:
		title = "Select File"
		for i, f := range dm.files {
			if i == dm.cursor {
				items = append(items, FileItemSelected.Render("> "+f))
			} else {
				items = append(items, FileItem.Render("  "+f))
			}
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		ModalTitle.Render(title),
		"",
		strings.Join(items, "\n"),
		"",
		FileItemDim.Render("Press Esc to go back/close"),
	)

	modal := ModalBox.Width(40).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}
