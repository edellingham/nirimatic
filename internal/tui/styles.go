package tui

import "github.com/charmbracelet/lipgloss"

// Eldritch color palette
// https://github.com/eldritch-theme/eldritch
var (
	ColorBackground  = lipgloss.Color("#212337") // Sunken Depths Grey
	ColorCurrentLine = lipgloss.Color("#323449") // Shallow Depths Grey
	ColorForeground  = lipgloss.Color("#ebfafa") // Lighthouse White
	ColorComment     = lipgloss.Color("#7081d0") // The Old One Purple
	ColorCyan        = lipgloss.Color("#04d1f9") // Watery Tomb Blue
	ColorGreen       = lipgloss.Color("#37f499") // Great Old One Green
	ColorOrange      = lipgloss.Color("#f7c67f") // Dreaming Orange
	ColorPink        = lipgloss.Color("#f265b5") // Pustule Pink
	ColorPurple      = lipgloss.Color("#a48cf2") // Lovecraft Purple
	ColorRed         = lipgloss.Color("#f16c75") // R'lyeh Red
	ColorYellow      = lipgloss.Color("#f1fc79") // Gold of Yuggoth
)

// Base styles
var (
	BaseStyle = lipgloss.NewStyle().
			Background(ColorBackground).
			Foreground(ColorForeground)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorComment).
			Italic(true)

	SelectedStyle = lipgloss.NewStyle().
			Background(ColorCurrentLine).
			Foreground(ColorGreen).
			Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(ColorForeground).
			PaddingLeft(2)

	DimmedStyle = lipgloss.NewStyle().
			Foreground(ColorComment)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorOrange)

	AccentStyle = lipgloss.NewStyle().
			Foreground(ColorPink)

	LinkStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Underline(true)
)

// Component styles
var (
	// Sidebar navigation
	SidebarStyle = lipgloss.NewStyle().
			Width(24).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple)

	SidebarTitleStyle = lipgloss.NewStyle().
				Foreground(ColorCyan).
				Bold(true).
				Padding(0, 1).
				MarginBottom(1)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(ColorForeground).
				PaddingLeft(1)

	SidebarActiveStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Bold(true).
				PaddingLeft(0)

	// Main content area
	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCurrentLine)

	// Header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBackground).
			Background(ColorCyan).
			Padding(0, 2).
			MarginBottom(1)

	// Footer / Help
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorComment).
			MarginTop(1)

	// Input fields
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorPurple).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorCyan).
				Padding(0, 1)

	// Buttons
	ButtonStyle = lipgloss.NewStyle().
			Foreground(ColorBackground).
			Background(ColorPurple).
			Padding(0, 2).
			MarginRight(1)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(ColorBackground).
				Background(ColorGreen).
				Padding(0, 2).
				MarginRight(1)

	// Cards/Boxes
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCurrentLine).
			Padding(1, 2).
			MarginBottom(1)

	CardTitleStyle = lipgloss.NewStyle().
			Foreground(ColorPink).
			Bold(true).
			MarginBottom(1)

	// Section headers
	SectionStyle = lipgloss.NewStyle().
			Foreground(ColorPurple).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	// Value display
	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorComment).
			Width(20)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorForeground)

	// Toggle/checkbox
	ToggleOnStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	ToggleOffStyle = lipgloss.NewStyle().
			Foreground(ColorComment)
)

// Status indicator styles
var (
	StatusOK = lipgloss.NewStyle().
			Foreground(ColorGreen)

	StatusError = lipgloss.NewStyle().
			Foreground(ColorRed)

	StatusWarning = lipgloss.NewStyle().
			Foreground(ColorOrange)

	StatusInactive = lipgloss.NewStyle().
			Foreground(ColorComment)
)

// Status symbols
const (
	SymbolOK       = "●"
	SymbolError    = "●"
	SymbolWarning  = "●"
	SymbolInactive = "○"
	SymbolArrow    = "▸"
	SymbolCheck    = "✓"
	SymbolCross    = "✗"
)

// GradientText creates a gradient-like effect on text using Eldritch colors
func GradientText(text string) string {
	colors := []lipgloss.Color{ColorPink, ColorPurple, ColorCyan}
	result := ""
	for i, char := range text {
		color := colors[i%len(colors)]
		result += lipgloss.NewStyle().Foreground(color).Render(string(char))
	}
	return result
}

// RenderStatus renders a status indicator with the appropriate style
func RenderStatus(status string) string {
	switch status {
	case "running", "active", "ok":
		return StatusOK.Render(SymbolOK)
	case "stopped", "inactive", "error":
		return StatusError.Render(SymbolError)
	case "warning":
		return StatusWarning.Render(SymbolWarning)
	default:
		return StatusInactive.Render(SymbolInactive)
	}
}

// RenderToggle renders a toggle with the appropriate style
func RenderToggle(enabled bool) string {
	if enabled {
		return ToggleOnStyle.Render("[✓]")
	}
	return ToggleOffStyle.Render("[ ]")
}
