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

	// Screen models (will be added as we implement them)
	dashboard *DashboardModel
	// niriSettings  *NiriSettingsModel
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

	// Initialize dashboard
	dashboard := NewDashboardModel()

	return &App{
		currentScreen: ScreenDashboard,
		sidebar:       sidebar,
		keys:          DefaultKeyMap(),
		configPath:    configPath,
		dashboard:     dashboard,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.dashboard.Init(),
	)
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
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

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

		// Update sidebar height
		a.sidebar.SetHeight(a.height - 6) // Account for header/footer

		// Update dashboard dimensions
		contentWidth := a.width - 28
		a.dashboard.SetSize(contentWidth, a.height-6)
	}

	// Update sidebar
	var cmd tea.Cmd
	a.sidebar, cmd = a.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	// Check for screen change
	if item, ok := a.sidebar.SelectedItem().(sidebarItem); ok {
		a.currentScreen = item.screen
	}

	// Update current screen
	switch a.currentScreen {
	case ScreenDashboard:
		var dashCmd tea.Cmd
		a.dashboard, dashCmd = a.dashboard.Update(msg)
		cmds = append(cmds, dashCmd)
	// TODO: Add other screens
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

	// Render sidebar
	sidebarContent := SidebarTitleStyle.Render("nirimatic") + "\n\n" + a.sidebar.View()
	sidebar := SidebarStyle.Height(a.height - 6).Render(sidebarContent)

	// Render current screen content
	var content string
	switch a.currentScreen {
	case ScreenDashboard:
		content = a.dashboard.View()
	case ScreenNiriSettings:
		content = "Niri Settings - Coming Soon"
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
	contentBox := ContentStyle.Width(contentWidth).Height(a.height - 6).Render(content)

	// Layout main area
	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, contentBox)

	// Render footer help
	help := HelpStyle.Render(HelpLine(
		a.keys.Up, a.keys.Down, a.keys.Enter,
		a.keys.Noctalia, a.keys.Reload, a.keys.Quit, a.keys.Help,
	))

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
