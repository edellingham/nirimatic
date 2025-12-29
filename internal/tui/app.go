package tui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/edellingham/nirimatic/internal/tui/screens"
)

// Screen represents the different screens in the app
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenNiriSettings
	ScreenAnimations
	ScreenKeybinds
	ScreenStartup
	ScreenBackup
)

// App version
const Version = "0.1.0"

// App is the main application model
type App struct {
	currentScreen Screen
	sidebar       list.Model
	width         int
	height        int
	keys          KeyMap
	ready         bool
	focusContent  bool // true = content focused, false = sidebar focused

	// Screen models
	dashboard    *DashboardModel
	niriSettings *screens.NiriSettingsModel
	// animations    *AnimationsModel
	// keybinds      *KeybindsModel
	// startup       *StartupModel
	// backup        *BackupModel

	// Config state
	configPath  string
	configDirty bool
}

// sidebarItem represents an item in the sidebar
type sidebarItem struct {
	title  string
	screen Screen
}

func (i sidebarItem) Title() string       { return i.title }
func (i sidebarItem) Description() string { return "" }
func (i sidebarItem) FilterValue() string { return i.title }

// sidebarDelegate is a custom delegate for the sidebar list
type sidebarDelegate struct{}

func (d sidebarDelegate) Height() int                             { return 1 }
func (d sidebarDelegate) Spacing() int                            { return 0 }
func (d sidebarDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d sidebarDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(sidebarItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	var line string
	if isSelected {
		line = SidebarActiveStyle.Render(SymbolArrow + " " + i.title)
	} else {
		line = SidebarItemStyle.Render("  " + i.title)
	}
	fmt.Fprint(w, line)
}

// NewApp creates a new application instance
func NewApp() *App {
	// Get config path
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "niri", "config.kdl")

	// Initialize sidebar
	items := []list.Item{
		sidebarItem{title: "Dashboard", screen: ScreenDashboard},
		sidebarItem{title: "Niri Settings", screen: ScreenNiriSettings},
		sidebarItem{title: "Animations", screen: ScreenAnimations},
		sidebarItem{title: "Keybinds", screen: ScreenKeybinds},
		sidebarItem{title: "Startup Apps", screen: ScreenStartup},
		sidebarItem{title: "Backup", screen: ScreenBackup},
	}

	sidebar := list.New(items, sidebarDelegate{}, 20, 14)
	sidebar.Title = "nirimatic"
	sidebar.SetShowStatusBar(false)
	sidebar.SetFilteringEnabled(false)
	sidebar.SetShowHelp(false)
	sidebar.SetShowTitle(false)
	sidebar.Styles.Title = SidebarTitleStyle

	// Initialize screen models
	dashboard := NewDashboardModel()
	niriSettings := screens.NewNiriSettingsModel()

	return &App{
		currentScreen: ScreenDashboard,
		sidebar:       sidebar,
		keys:          DefaultKeyMap(),
		configPath:    configPath,
		dashboard:     dashboard,
		niriSettings:  niriSettings,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.dashboard.Init(),
		a.niriSettings.Init(),
	)
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keys work regardless of focus
		switch {
		case key.Matches(msg, a.keys.Quit):
			if a.configDirty {
				// TODO: Prompt to save
				return a, tea.Quit
			}
			return a, tea.Quit

		case key.Matches(msg, a.keys.Noctalia):
			return a, openNoctaliaSettings()

		case key.Matches(msg, a.keys.Reload):
			return a, reloadNiriConfig()
		}

		// Focus switching
		if a.focusContent {
			// Escape returns to sidebar
			if msg.String() == "esc" {
				a.focusContent = false
				return a, nil
			}
			// Route to content screen
			switch a.currentScreen {
			case ScreenDashboard:
				var dashCmd tea.Cmd
				a.dashboard, dashCmd = a.dashboard.Update(msg)
				cmds = append(cmds, dashCmd)
			case ScreenNiriSettings:
				var settingsCmd tea.Cmd
				a.niriSettings, settingsCmd = a.niriSettings.Update(msg)
				cmds = append(cmds, settingsCmd)
			}
			return a, tea.Batch(cmds...)
		} else {
			// Sidebar focused: Enter switches to content
			if key.Matches(msg, a.keys.Enter) {
				a.focusContent = true
				return a, nil
			}
			// Update sidebar navigation
			var cmd tea.Cmd
			a.sidebar, cmd = a.sidebar.Update(msg)
			cmds = append(cmds, cmd)

			// Check for screen change
			if item, ok := a.sidebar.SelectedItem().(sidebarItem); ok {
				a.currentScreen = item.screen
			}
			return a, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

		// Update sidebar height
		a.sidebar.SetHeight(a.height - 6) // Account for header/footer

		// Update screen dimensions
		contentWidth := a.width - 28
		a.dashboard.SetSize(contentWidth, a.height-6)
		a.niriSettings.SetSize(contentWidth, a.height-6)
	}

	// Pass non-key messages to current screen
	switch a.currentScreen {
	case ScreenDashboard:
		var dashCmd tea.Cmd
		a.dashboard, dashCmd = a.dashboard.Update(msg)
		cmds = append(cmds, dashCmd)
	case ScreenNiriSettings:
		var settingsCmd tea.Cmd
		a.niriSettings, settingsCmd = a.niriSettings.Update(msg)
		cmds = append(cmds, settingsCmd)
	}

	return a, tea.Batch(cmds...)
}

// View renders the application
func (a *App) View() string {
	if !a.ready {
		return "Loading..."
	}

	// Render header
	headerText := GradientText("▄▄ nirimatic") + "  v" + Version
	header := HeaderStyle.Width(a.width - 4).Render(headerText)

	// Render sidebar with focus indicator
	sidebarContent := SidebarTitleStyle.Render("nirimatic") + "\n\n" + a.sidebar.View()
	sidebarStyle := SidebarStyle.Height(a.height - 6)
	if !a.focusContent {
		sidebarStyle = sidebarStyle.BorderForeground(ColorCyan) // Highlight when focused
	}
	sidebar := sidebarStyle.Render(sidebarContent)

	// Render current screen content
	var content string
	switch a.currentScreen {
	case ScreenDashboard:
		content = a.dashboard.View()
	case ScreenNiriSettings:
		content = a.niriSettings.View()
	case ScreenAnimations:
		content = "Animations - Coming Soon"
	case ScreenKeybinds:
		content = "Keybinds - Coming Soon"
	case ScreenStartup:
		content = "Startup Apps - Coming Soon"
	case ScreenBackup:
		content = "Backup - Coming Soon"
	}

	contentWidth := a.width - 28
	contentStyle := ContentStyle.Width(contentWidth).Height(a.height - 6)
	if a.focusContent {
		contentStyle = contentStyle.BorderForeground(ColorCyan) // Highlight when focused
	}
	contentBox := contentStyle.Render(content)

	// Layout main area
	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, contentBox)

	// Render footer help based on focus
	var helpText string
	if a.focusContent {
		helpText = DimmedStyle.Render("esc") + " back  " + HelpLine(
			a.keys.Up, a.keys.Down, a.keys.Enter,
			a.keys.Noctalia, a.keys.Reload, a.keys.Quit, a.keys.Help,
		)
	} else {
		helpText = HelpLine(
			a.keys.Up, a.keys.Down, a.keys.Enter,
			a.keys.Noctalia, a.keys.Reload, a.keys.Quit, a.keys.Help,
		)
	}
	help := HelpStyle.Render(helpText)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, help)
}

// Commands

// openNoctaliaSettings opens the Noctalia settings panel
func openNoctaliaSettings() tea.Cmd {
	return func() tea.Msg {
		cmd := newExecCmd("qs", "-c", "noctalia-shell", "ipc", "call", "settings", "toggle")
		_ = cmd.Run()
		return nil
	}
}

// reloadNiriConfig reloads the Niri configuration
func reloadNiriConfig() tea.Cmd {
	return func() tea.Msg {
		cmd := newExecCmd("niri", "msg", "action", "reload-config")
		_ = cmd.Run()
		return nil
	}
}

// newExecCmd creates an exec.Cmd for running shell commands
func newExecCmd(name string, args ...string) *execCmd {
	return &execCmd{name: name, args: args}
}

type execCmd struct {
	name string
	args []string
}

func (c *execCmd) Run() error {
	cmd := exec.Command(c.name, c.args...)
	return cmd.Run()
}
