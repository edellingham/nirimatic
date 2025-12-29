package tui

// Re-export from shared styles package for backwards compatibility
import (
	"github.com/charmbracelet/lipgloss"
	"github.com/edellingham/nirimatic/internal/styles"
)

// Re-export colors
var (
	ColorBackground  = styles.ColorBackground
	ColorCurrentLine = styles.ColorCurrentLine
	ColorForeground  = styles.ColorForeground
	ColorComment     = styles.ColorComment
	ColorCyan        = styles.ColorCyan
	ColorGreen       = styles.ColorGreen
	ColorOrange      = styles.ColorOrange
	ColorPink        = styles.ColorPink
	ColorPurple      = styles.ColorPurple
	ColorRed         = styles.ColorRed
	ColorYellow      = styles.ColorYellow
)

// Re-export base styles
var (
	BaseStyle       = styles.BaseStyle
	TitleStyle      = styles.TitleStyle
	SubtitleStyle   = styles.SubtitleStyle
	SelectedStyle   = styles.SelectedStyle
	NormalItemStyle = styles.NormalItemStyle
	DimmedStyle     = styles.DimmedStyle
	ErrorStyle      = styles.ErrorStyle
	SuccessStyle    = styles.SuccessStyle
	WarningStyle    = styles.WarningStyle
	AccentStyle     = styles.AccentStyle
	LinkStyle       = styles.LinkStyle
)

// Re-export component styles
var (
	SidebarStyle       = styles.SidebarStyle
	SidebarTitleStyle  = styles.SidebarTitleStyle
	SidebarItemStyle   = styles.SidebarItemStyle
	SidebarActiveStyle = styles.SidebarActiveStyle
	ContentStyle       = styles.ContentStyle
	HeaderStyle        = styles.HeaderStyle
	HelpStyle          = styles.HelpStyle
	InputStyle         = styles.InputStyle
	InputFocusedStyle  = styles.InputFocusedStyle
	ButtonStyle        = styles.ButtonStyle
	ButtonActiveStyle  = styles.ButtonActiveStyle
	CardStyle          = styles.CardStyle
	CardTitleStyle     = styles.CardTitleStyle
	SectionStyle       = styles.SectionStyle
	LabelStyle         = styles.LabelStyle
	ValueStyle         = styles.ValueStyle
	ToggleOnStyle      = styles.ToggleOnStyle
	ToggleOffStyle     = styles.ToggleOffStyle
)

// Re-export status styles
var (
	StatusOK       = styles.StatusOK
	StatusError    = styles.StatusError
	StatusWarning  = styles.StatusWarning
	StatusInactive = styles.StatusInactive
)

// Re-export symbols
const (
	SymbolOK       = styles.SymbolOK
	SymbolError    = styles.SymbolError
	SymbolWarning  = styles.SymbolWarning
	SymbolInactive = styles.SymbolInactive
	SymbolArrow    = styles.SymbolArrow
	SymbolCheck    = styles.SymbolCheck
	SymbolCross    = styles.SymbolCross
)

// Re-export functions
var GradientText = styles.GradientText
var RenderStatus = styles.RenderStatus
var RenderToggle = styles.RenderToggle

// Additional helper - keeping lipgloss import used
var _ = lipgloss.Color("")
