package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// BuildStatus represents the current compiler status.
type BuildStatus int

const (
	StatusIDLE BuildStatus = iota
	StatusBUILDING
	StatusSUCCESS
	StatusFAILED
)

// ActionBar is the bottom status/shortcut bar.
type ActionBar struct {
	width        int
	status       BuildStatus
	lastBuild    time.Duration
	errorCount   int
	projectRoot  string
	shortcutsStr string // Cached rendered shortcuts
}

// NewActionBar creates a new action bar.
func NewActionBar() ActionBar {
	return ActionBar{
		status: StatusIDLE,
	}
}

// SetWidth sets the action bar width and pre-renders shortcuts.
func (ab *ActionBar) SetWidth(w int) {
	ab.width = w
	ab.rebuildShortcuts()
}

func (ab *ActionBar) rebuildShortcuts() {
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"F5", "Compile"},
		{"F6", "Clean"},
		{"F7", "PDF"},
		{"h", "Shadow"},
		{"n", "New"},
		{"y", "Yank"},
		{"?", "Help"},
		{"q", "Quit"},
	}

	var parts []string
	for _, s := range shortcuts {
		parts = append(parts, ActionKey.Render(s.key)+" "+ActionDesc.Render(s.desc))
	}

	ab.shortcutsStr = strings.Join(parts, ActionSep.String())
}

// SetBuildStatus updates the build status display.
func (ab *ActionBar) SetBuildStatus(s BuildStatus, duration time.Duration, errors int) {
	ab.status = s
	ab.lastBuild = duration
	ab.errorCount = errors
}

// SetProjectRoot sets the displayed project path.
func (ab *ActionBar) SetProjectRoot(root string) {
	ab.projectRoot = root
}

// View renders the action bar.
func (ab ActionBar) View() string {
	left := ab.shortcutsStr

	// Status indicator
	var statusStr string
	switch ab.status {
	case StatusIDLE:
		statusStr = StatusIdle.Render("● Idle")
	case StatusBUILDING:
		statusStr = StatusBuilding.Render("◉ Building…")
	case StatusSUCCESS:
		statusStr = StatusSuccess.Render(fmt.Sprintf("✓ Built (%.1fs)", ab.lastBuild.Seconds()))
	case StatusFAILED:
		statusStr = StatusFailed.Render(fmt.Sprintf("✗ Failed (%d errors)", ab.errorCount))
	}

	// Layout: shortcuts on left, status on right
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(statusStr)
	gap := ab.width - leftWidth - rightWidth - 4

	if gap < 0 {
		gap = 1
	}

	bar := left + strings.Repeat(" ", gap) + statusStr

	return ActionBarBg.Width(ab.width).Render(bar)
}
