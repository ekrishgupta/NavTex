package ui

import "github.com/charmbracelet/lipgloss"

// ── Color Palette ──
// A warm, academic-inspired dark theme.
var (
	ColorBg        = lipgloss.AdaptiveColor{Light: "#FAFAFA", Dark: "#1A1B26"}
	ColorFg        = lipgloss.AdaptiveColor{Light: "#24283B", Dark: "#C0CAF5"}
	ColorDim       = lipgloss.AdaptiveColor{Light: "#9699A3", Dark: "#565F89"}
	ColorAccent    = lipgloss.AdaptiveColor{Light: "#2E7DE9", Dark: "#7AA2F7"}
	ColorGreen     = lipgloss.AdaptiveColor{Light: "#587539", Dark: "#9ECE6A"}
	ColorYellow    = lipgloss.AdaptiveColor{Light: "#8C6C3E", Dark: "#E0AF68"}
	ColorRed       = lipgloss.AdaptiveColor{Light: "#F52A65", Dark: "#F7768E"}
	ColorCyan      = lipgloss.AdaptiveColor{Light: "#007197", Dark: "#7DCFFF"}
	ColorMagenta   = lipgloss.AdaptiveColor{Light: "#9854F1", Dark: "#BB9AF7"}
	ColorBorder    = lipgloss.AdaptiveColor{Light: "#DFE0E5", Dark: "#3B4261"}
	ColorHighlight = lipgloss.AdaptiveColor{Light: "#E1E2E7", Dark: "#283457"}
)

// ── Section Header Styles ──
var (
	SectionHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			PaddingLeft(1).
			MarginBottom(0)

	SectionHeaderActive = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorBg).
				Background(ColorAccent).
				PaddingLeft(1).
				PaddingRight(1)
)

// ── File Browser Styles ──
var (
	CategoryLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			PaddingLeft(1)

	FileItem = lipgloss.NewStyle().
			PaddingLeft(3).
			Foreground(ColorFg)

	FileItemSelected = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorBg).
				Background(ColorAccent).
				Bold(true)

	FileItemDim = lipgloss.NewStyle().
			PaddingLeft(3).
			Foreground(ColorDim)

	ShadowBinLabel = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			PaddingLeft(1)
)

// ── Inspector Styles ──
var (
	InspectorTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMagenta).
			MarginBottom(1).
			PaddingLeft(1)

	MetaLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan).
			Width(14).
			PaddingLeft(1)

	MetaValue = lipgloss.NewStyle().
			Foreground(ColorFg)

	PackageTag = lipgloss.NewStyle().
			Foreground(ColorGreen).
			PaddingRight(1)

	KeywordTag = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Italic(true).
			PaddingRight(1)

	BibTableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			PaddingRight(2)

	BibTableRow = lipgloss.NewStyle().
			Foreground(ColorFg).
			PaddingRight(2)

	BibTableRowSelected = lipgloss.NewStyle().
				Foreground(ColorBg).
				Background(ColorAccent).
				PaddingRight(2)

	ErrorText = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	WarningText = lipgloss.NewStyle().
			Foreground(ColorYellow)

	SuccessText = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)
)

// ── Pane Styles ──
var (
	PaneBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	PaneBorderActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorAccent).
				Padding(0, 1)

	PaneTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Padding(0, 1)
)

// ── Action Bar Styles ──
var (
	ActionBarStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Background(lipgloss.AdaptiveColor{Light: "#E1E2E7", Dark: "#24283B"}).
			Padding(0, 1)

	ActionKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Background(lipgloss.AdaptiveColor{Light: "#E1E2E7", Dark: "#24283B"})

	ActionDesc = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(lipgloss.AdaptiveColor{Light: "#E1E2E7", Dark: "#24283B"})

	ActionSep = lipgloss.NewStyle().
			Foreground(ColorBorder).
			Background(lipgloss.AdaptiveColor{Light: "#E1E2E7", Dark: "#24283B"}).
			SetString(" │ ")

	StatusIdle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	StatusBuilding = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Bold(true)

	StatusError = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)
)

// ── Modal Styles ──
var (
	ModalOverlay = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "#00000033", Dark: "#00000088"})

	ModalBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2).
			Background(ColorBg)

	ModalTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	InputLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan).
			PaddingRight(1)

	InputField = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)
)

// ── Logo ──
const Logo = `╔╗╔┌─┐┬  ┬╔╦╗┌─┐─┐ ┬
║║║├─┤└┐┌┘ ║ ├┤ ┌┘│
╝╚╝┴ ┴ └┘  ╩ └─┘└─┘`

var LogoStyle = lipgloss.NewStyle().
	Foreground(ColorAccent).
	Bold(true).
	Align(lipgloss.Center)
