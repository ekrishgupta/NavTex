package ui

import "github.com/charmbracelet/lipgloss"

// ── Color Palette ──
// A nap-inspired subdued dark theme with muted, restful colors.
var (
	ColorBg             = lipgloss.Color("0")       // True black background
	ColorFg             = lipgloss.Color("7")       // Terminal white
	ColorPrimary        = lipgloss.Color("#AFBEE1")  // Subdued bright blue
	ColorPrimarySubdued = lipgloss.Color("#64708D")  // Muted blue
	ColorGreen          = lipgloss.Color("#527251")  // Dark green
	ColorBrightGreen    = lipgloss.Color("#BCE1AF")  // Bright green
	ColorRed            = lipgloss.Color("#A46060")  // Dark red
	ColorBrightRed      = lipgloss.Color("#E49393")  // Bright red
	ColorYellow         = lipgloss.Color("#E0AF68")  // Warm amber
	ColorBlack          = lipgloss.Color("#373B41")  // Off-black for borders
	ColorGray           = lipgloss.Color("240")      // Neutral gray
	ColorWhite          = lipgloss.Color("#FFFFFF")  // Pure white for contrast
	ColorDimBg          = lipgloss.Color("#24283B")  // Slightly elevated bg
	ColorSurfaceBg      = lipgloss.Color("#1E2030")  // Surface for title bars (blurred)
)

// ══════════════════════════════════════════════════════════════
// Pane Style System — Focused & Blurred variants per pane
// Inspired by github.com/maaslalani/nap
// ══════════════════════════════════════════════════════════════

// BrowserBaseStyle holds the visual properties for the file browser pane.
type BrowserBaseStyle struct {
	Base             lipgloss.Style
	TitleBar         lipgloss.Style
	TitleText        lipgloss.Style
	CategoryHeader   lipgloss.Style
	SelectedItem     lipgloss.Style
	UnselectedItem   lipgloss.Style
	DimItem          lipgloss.Style
	SelectedPrefix   string
	UnselectedPrefix string
}

// BrowserStyle holds focused and blurred browser styles.
type BrowserStyle struct {
	Focused BrowserBaseStyle
	Blurred BrowserBaseStyle
}

// InspectorBaseStyle holds the visual properties for the inspector pane.
type InspectorBaseStyle struct {
	Base          lipgloss.Style
	TitleBar      lipgloss.Style
	TitleText     lipgloss.Style
	SectionTitle  lipgloss.Style
	MetaLabel     lipgloss.Style
	MetaValue     lipgloss.Style
	Separator     lipgloss.Style
	SelectedRow   lipgloss.Style
	UnselectedRow lipgloss.Style
	PackageTag    lipgloss.Style
	KeywordTag    lipgloss.Style
	ErrorText     lipgloss.Style
	WarningText   lipgloss.Style
	SuccessText   lipgloss.Style
}

// InspectorStyle holds focused and blurred inspector styles.
type InspectorStyle struct {
	Focused InspectorBaseStyle
	Blurred InspectorBaseStyle
}

// Styles is the root container for all application styles.
type Styles struct {
	Browser   BrowserStyle
	Inspector InspectorStyle
}

// DefaultStyles configures the application's visual system.
func DefaultStyles() Styles {
	return Styles{
		Browser: BrowserStyle{
			Focused: BrowserBaseStyle{
				Base: lipgloss.NewStyle(),
				TitleBar: lipgloss.NewStyle().
					Background(ColorPrimarySubdued).
					Foreground(ColorWhite).
					Padding(0, 1).
					Bold(true),
				TitleText: lipgloss.NewStyle().
					Foreground(ColorWhite).
					Bold(true),
				CategoryHeader: lipgloss.NewStyle().
					Foreground(ColorPrimary).
					Bold(true).
					PaddingLeft(1).
					MarginTop(1),
				SelectedItem: lipgloss.NewStyle().
					Foreground(ColorPrimary).
					Bold(true).
					PaddingLeft(1),
				UnselectedItem: lipgloss.NewStyle().
					Foreground(ColorFg).
					PaddingLeft(3),
				DimItem: lipgloss.NewStyle().
					Foreground(ColorGray).
					PaddingLeft(3),
				SelectedPrefix:   " ▸ ",
				UnselectedPrefix: "   ",
			},
			Blurred: BrowserBaseStyle{
				Base: lipgloss.NewStyle(),
				TitleBar: lipgloss.NewStyle().
					Background(ColorSurfaceBg).
					Foreground(ColorGray).
					Padding(0, 1),
				TitleText: lipgloss.NewStyle().
					Foreground(ColorGray),
				CategoryHeader: lipgloss.NewStyle().
					Foreground(ColorPrimarySubdued).
					Bold(true).
					PaddingLeft(1).
					MarginTop(1),
				SelectedItem: lipgloss.NewStyle().
					Foreground(ColorPrimary).
					PaddingLeft(1),
				UnselectedItem: lipgloss.NewStyle().
					Foreground(ColorGray).
					PaddingLeft(3),
				DimItem: lipgloss.NewStyle().
					Foreground(lipgloss.Color("237")).
					PaddingLeft(3),
				SelectedPrefix:   " ▸ ",
				UnselectedPrefix: "   ",
			},
		},
		Inspector: InspectorStyle{
			Focused: InspectorBaseStyle{
				Base: lipgloss.NewStyle(),
				TitleBar: lipgloss.NewStyle().
					Background(ColorPrimarySubdued).
					Foreground(ColorWhite).
					Padding(0, 1).
					Bold(true),
				TitleText: lipgloss.NewStyle().
					Foreground(ColorWhite).
					Bold(true),
				SectionTitle: lipgloss.NewStyle().
					Foreground(ColorPrimary).
					Bold(true).
					PaddingLeft(1).
					MarginTop(1),
				MetaLabel: lipgloss.NewStyle().
					Foreground(ColorPrimarySubdued).
					Bold(true).
					Width(16).
					PaddingLeft(2),
				MetaValue: lipgloss.NewStyle().
					Foreground(ColorFg),
				Separator: lipgloss.NewStyle().
					Foreground(ColorBlack).
					MarginTop(0).
					MarginBottom(0),
				SelectedRow: lipgloss.NewStyle().
					Foreground(ColorPrimary).
					Bold(true),
				UnselectedRow: lipgloss.NewStyle().
					Foreground(ColorFg),
				PackageTag: lipgloss.NewStyle().
					Foreground(ColorBrightGreen).
					PaddingRight(1),
				KeywordTag: lipgloss.NewStyle().
					Foreground(ColorYellow).
					Italic(true).
					PaddingRight(1),
				ErrorText: lipgloss.NewStyle().
					Foreground(ColorBrightRed).
					Bold(true),
				WarningText: lipgloss.NewStyle().
					Foreground(ColorYellow),
				SuccessText: lipgloss.NewStyle().
					Foreground(ColorBrightGreen).
					Bold(true),
			},
			Blurred: InspectorBaseStyle{
				Base: lipgloss.NewStyle(),
				TitleBar: lipgloss.NewStyle().
					Background(ColorSurfaceBg).
					Foreground(ColorGray).
					Padding(0, 1),
				TitleText: lipgloss.NewStyle().
					Foreground(ColorGray),
				SectionTitle: lipgloss.NewStyle().
					Foreground(ColorPrimarySubdued).
					Bold(true).
					PaddingLeft(1).
					MarginTop(1),
				MetaLabel: lipgloss.NewStyle().
					Foreground(ColorPrimarySubdued).
					Width(16).
					PaddingLeft(2),
				MetaValue: lipgloss.NewStyle().
					Foreground(ColorGray),
				Separator: lipgloss.NewStyle().
					Foreground(ColorBlack),
				SelectedRow: lipgloss.NewStyle().
					Foreground(ColorPrimary),
				UnselectedRow: lipgloss.NewStyle().
					Foreground(ColorGray),
				PackageTag: lipgloss.NewStyle().
					Foreground(ColorGreen).
					PaddingRight(1),
				KeywordTag: lipgloss.NewStyle().
					Foreground(ColorYellow).
					Italic(true).
					PaddingRight(1),
				ErrorText: lipgloss.NewStyle().
					Foreground(ColorBrightRed).
					Bold(true),
				WarningText: lipgloss.NewStyle().
					Foreground(ColorYellow),
				SuccessText: lipgloss.NewStyle().
					Foreground(ColorBrightGreen).
					Bold(true),
			},
		},
	}
}

// ══════════════════════════════════════════════════════════════
// Global Utility Styles — shared across components
// ══════════════════════════════════════════════════════════════

var (
	// DimText is for hints and secondary information.
	DimText = lipgloss.NewStyle().Foreground(ColorGray)

	// ModalFrame is the shared border for all modal dialogs.
	ModalFrame = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimarySubdued).
			Padding(1, 2).
			Background(ColorBg)

	// ModalTitleBar is a colored title bar inside modals.
	ModalTitleBar = lipgloss.NewStyle().
			Background(ColorPrimarySubdued).
			Foreground(ColorWhite).
			Padding(0, 1).
			Bold(true)

	// ModalHint is the dim instructions at the bottom of modals.
	ModalHint = lipgloss.NewStyle().
			Foreground(ColorGray).
			PaddingLeft(1).
			MarginTop(1)
)

// ── Action Bar Styles ──
var (
	ActionBarBg = lipgloss.NewStyle().
			Background(ColorDimBg).
			Padding(0, 1)

	ActionKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(ColorDimBg)

	ActionDesc = lipgloss.NewStyle().
			Foreground(ColorGray).
			Background(ColorDimBg)

	ActionSep = lipgloss.NewStyle().
			Foreground(ColorBlack).
			Background(ColorDimBg).
			SetString(" │ ")

	StatusIdle = lipgloss.NewStyle().
			Foreground(ColorBrightGreen).
			Background(ColorDimBg).
			Bold(true)

	StatusBuilding = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Background(ColorDimBg).
			Bold(true)

	StatusFailed = lipgloss.NewStyle().
			Foreground(ColorBrightRed).
			Background(ColorDimBg).
			Bold(true)

	StatusSuccess = lipgloss.NewStyle().
			Foreground(ColorBrightGreen).
			Background(ColorDimBg).
			Bold(true)
)

// ── Input Styles (for new project modal etc.) ──
var (
	InputLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimarySubdued).
			PaddingRight(1)

	InputField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBlack).
			Padding(0, 1)

	InputFieldActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)
)

// ── Logo ──
const Logo = `╔╗╔┌─┐┬  ┬╔╦╗┌─┐─┐ ┬
║║║├─┤└┐┌┘ ║ ├┤ ┌┘│
╝╚╝┴ ┴ └┘  ╩ └─┘└─┘`

var LogoStyle = lipgloss.NewStyle().
	Foreground(ColorPrimary).
	Bold(true).
	Align(lipgloss.Center)

// ── Helpers ──

// SeparatorLine returns a thin line separator for the given width.
func SeparatorLine(width int) string {
	if width <= 0 {
		return ""
	}
	sep := lipgloss.NewStyle().Foreground(ColorBlack)
	return sep.Render(repeatChar('─', width))
}

func repeatChar(ch rune, n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
