package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application
type KeyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Tab    key.Binding
	ShiftTab key.Binding

	// Actions
	Quit       key.Binding
	Help       key.Binding
	Save       key.Binding
	Reload     key.Binding
	Noctalia   key.Binding
	RestartNiri key.Binding

	// Editing
	Toggle key.Binding
	Delete key.Binding
	Add    key.Binding
	Edit   key.Binding

	// Filtering
	Filter key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next section"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev section"),
		),

		// Actions
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Save: key.NewBinding(
			key.WithKeys("s", "ctrl+s"),
			key.WithHelp("s", "save"),
		),
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload config"),
		),
		Noctalia: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "noctalia"),
		),
		RestartNiri: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "restart niri"),
		),

		// Editing
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "delete"),
			key.WithHelp("d", "delete"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),

		// Filtering
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
	}
}

// ShortHelp returns a short help string for the footer
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit, k.Help}
}

// FullHelp returns the full help for the help overlay
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Tab, k.ShiftTab},
		{k.Save, k.Reload, k.Noctalia, k.RestartNiri},
		{k.Toggle, k.Add, k.Edit, k.Delete},
		{k.Filter, k.Help, k.Quit},
	}
}

// HelpLine generates a compact help line for the footer
func HelpLine(keys ...key.Binding) string {
	var parts []string
	for _, k := range keys {
		if k.Enabled() {
			help := k.Help()
			parts = append(parts, help.Key+" "+help.Desc)
		}
	}
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " • "
		}
		result += part
	}
	return result
}
