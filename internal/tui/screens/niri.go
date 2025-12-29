package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/edellingham/nirimatic/internal/config"
	"github.com/edellingham/nirimatic/internal/styles"
)

// Field represents an editable field in the settings
type Field struct {
	Label    string
	Value    int
	Min      int
	Max      int
	Step     int
	Unit     string
	IsToggle bool
	Enabled  bool
}

// NiriSettingsModel is the model for the Niri settings screen
type NiriSettingsModel struct {
	config   *config.NiriConfig
	fields   []Field
	cursor   int
	width    int
	height   int
	dirty    bool
	err      error
	message  string
}

// configLoadedMsg is sent when config is loaded
type configLoadedMsg struct {
	config *config.NiriConfig
	err    error
}

// configSavedMsg is sent when config is saved
type configSavedMsg struct {
	err error
}

// NewNiriSettingsModel creates a new Niri settings model
func NewNiriSettingsModel() *NiriSettingsModel {
	return &NiriSettingsModel{
		fields: []Field{
			// Layout section
			{Label: "Gaps", Value: 10, Min: 0, Max: 50, Step: 1, Unit: "px"},
			{Label: "Border Width", Value: 2, Min: 0, Max: 10, Step: 1, Unit: "px"},
			{Label: "Focus Ring Width", Value: 0, Min: 0, Max: 10, Step: 1, Unit: "px"},
			{Label: "Corner Radius", Value: 16, Min: 0, Max: 32, Step: 1, Unit: "px"},
			// Shadow section
			{Label: "Shadows", IsToggle: true, Enabled: true},
			{Label: "Shadow Softness", Value: 60, Min: 0, Max: 100, Step: 5, Unit: ""},
			{Label: "Shadow Spread", Value: 10, Min: 0, Max: 50, Step: 1, Unit: ""},
			// Behavior section
			{Label: "Focus Follows Mouse", IsToggle: true, Enabled: true},
			{Label: "Workspace Auto Back", IsToggle: true, Enabled: true},
		},
	}
}

// Init initializes the model
func (m *NiriSettingsModel) Init() tea.Cmd {
	return m.loadConfig()
}

// SetSize sets the dimensions
func (m *NiriSettingsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages
func (m *NiriSettingsModel) Update(msg tea.Msg) (*NiriSettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case configLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.config = msg.config
		m.syncFromConfig()
		return m, nil

	case configSavedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error saving: %v", msg.err)
		} else {
			m.message = "Configuration saved!"
			m.dirty = false
		}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyUp):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keyDown):
			if m.cursor < len(m.fields)-1 {
				m.cursor++
			}
		case key.Matches(msg, keyLeft):
			m.decreaseValue()
		case key.Matches(msg, keyRight):
			m.increaseValue()
		case key.Matches(msg, keyToggle):
			m.toggleValue()
		case key.Matches(msg, keySave):
			return m, m.saveConfig()
		case key.Matches(msg, keyReset):
			return m, m.loadConfig()
		}
	}

	return m, nil
}

// View renders the settings screen
func (m *NiriSettingsModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Niri Settings"))
	b.WriteString("\n")
	b.WriteString(styles.SectionStyle.Render("─────────────────────────────────────────"))
	b.WriteString("\n\n")

	// Error display
	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Message display
	if m.message != "" {
		b.WriteString(styles.SuccessStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	// Layout section
	b.WriteString(styles.CardTitleStyle.Render("Layout"))
	b.WriteString("\n\n")

	for i := 0; i < 4; i++ {
		b.WriteString(m.renderField(i))
	}

	b.WriteString("\n")
	b.WriteString(styles.CardTitleStyle.Render("Shadows"))
	b.WriteString("\n\n")

	for i := 4; i < 7; i++ {
		b.WriteString(m.renderField(i))
	}

	b.WriteString("\n")
	b.WriteString(styles.CardTitleStyle.Render("Behavior"))
	b.WriteString("\n\n")

	for i := 7; i < len(m.fields); i++ {
		b.WriteString(m.renderField(i))
	}

	// Dirty indicator
	if m.dirty {
		b.WriteString("\n")
		b.WriteString(styles.WarningStyle.Render("* Unsaved changes"))
	}

	// Help line
	b.WriteString("\n\n")
	b.WriteString(styles.DimmedStyle.Render("↑↓ navigate • ←→ adjust • space toggle • s save • r reload"))

	return b.String()
}

// renderField renders a single field
func (m *NiriSettingsModel) renderField(index int) string {
	field := m.fields[index]
	isSelected := index == m.cursor

	// Label
	labelStyle := styles.LabelStyle
	if isSelected {
		labelStyle = labelStyle.Foreground(styles.ColorGreen)
	}
	label := labelStyle.Width(22).Render(field.Label)

	// Value
	var valueStr string
	if field.IsToggle {
		valueStr = m.renderToggle(field.Enabled, isSelected)
	} else {
		valueStr = m.renderSlider(field.Value, field.Min, field.Max, isSelected)
		if field.Unit != "" {
			valueStr += " " + field.Unit
		}
	}

	// Cursor indicator
	cursor := "  "
	if isSelected {
		cursor = styles.SuccessStyle.Render(styles.SymbolArrow + " ")
	}

	return fmt.Sprintf("%s%s %s\n", cursor, label, valueStr)
}

// renderSlider renders a slider with the current value
func (m *NiriSettingsModel) renderSlider(value, min, max int, selected bool) string {
	// Calculate position
	width := 20
	range_ := max - min
	if range_ == 0 {
		range_ = 1
	}
	pos := (value - min) * width / range_
	if pos > width {
		pos = width
	}
	if pos < 0 {
		pos = 0
	}

	// Build slider
	slider := strings.Repeat("─", pos) + "●" + strings.Repeat("─", width-pos)

	// Style
	style := styles.DimmedStyle
	if selected {
		style = lipgloss.NewStyle().Foreground(styles.ColorCyan)
	}

	valueStyle := styles.ValueStyle
	if selected {
		valueStyle = valueStyle.Foreground(styles.ColorGreen).Bold(true)
	}

	return fmt.Sprintf("[%s] %s", style.Render(slider), valueStyle.Render(fmt.Sprintf("%d", value)))
}

// renderToggle renders a toggle
func (m *NiriSettingsModel) renderToggle(enabled bool, selected bool) string {
	if enabled {
		style := styles.ToggleOnStyle
		if selected {
			style = style.Bold(true)
		}
		return style.Render("[✓] Enabled")
	}
	style := styles.ToggleOffStyle
	if selected {
		style = style.Foreground(styles.ColorComment)
	}
	return style.Render("[ ] Disabled")
}

// syncFromConfig syncs field values from the loaded config
func (m *NiriSettingsModel) syncFromConfig() {
	if m.config == nil {
		return
	}

	m.fields[0].Value = m.config.Gaps
	m.fields[1].Value = m.config.BorderWidth
	m.fields[2].Value = m.config.FocusRingWidth
	m.fields[3].Value = m.config.CornerRadius
	m.fields[4].Enabled = m.config.ShadowEnabled
	m.fields[5].Value = m.config.ShadowSoftness
	m.fields[6].Value = m.config.ShadowSpread
	m.fields[7].Enabled = m.config.FocusFollowsMouse
	m.fields[8].Enabled = m.config.WorkspaceAutoBackAndForth
}

// syncToConfig syncs field values to the config
func (m *NiriSettingsModel) syncToConfig() {
	if m.config == nil || len(m.fields) < 9 {
		return
	}

	m.config.Gaps = m.fields[0].Value
	m.config.BorderWidth = m.fields[1].Value
	m.config.FocusRingWidth = m.fields[2].Value
	m.config.CornerRadius = m.fields[3].Value
	m.config.ShadowEnabled = m.fields[4].Enabled
	m.config.ShadowSoftness = m.fields[5].Value
	m.config.ShadowSpread = m.fields[6].Value
	m.config.FocusFollowsMouse = m.fields[7].Enabled
	m.config.WorkspaceAutoBackAndForth = m.fields[8].Enabled
}

// increaseValue increases the current field's value
func (m *NiriSettingsModel) increaseValue() {
	field := &m.fields[m.cursor]
	if field.IsToggle {
		return
	}
	if field.Value < field.Max {
		field.Value += field.Step
		if field.Value > field.Max {
			field.Value = field.Max
		}
		m.dirty = true
	}
}

// decreaseValue decreases the current field's value
func (m *NiriSettingsModel) decreaseValue() {
	field := &m.fields[m.cursor]
	if field.IsToggle {
		return
	}
	if field.Value > field.Min {
		field.Value -= field.Step
		if field.Value < field.Min {
			field.Value = field.Min
		}
		m.dirty = true
	}
}

// toggleValue toggles the current field
func (m *NiriSettingsModel) toggleValue() {
	field := &m.fields[m.cursor]
	if field.IsToggle {
		field.Enabled = !field.Enabled
		m.dirty = true
	}
}

// loadConfig loads the config file
func (m *NiriSettingsModel) loadConfig() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.LoadNiriConfig(config.GetConfigPath())
		return configLoadedMsg{config: cfg, err: err}
	}
}

// saveConfig saves the config file
func (m *NiriSettingsModel) saveConfig() tea.Cmd {
	return func() (msg tea.Msg) {
		// Recover from any panics
		defer func() {
			if r := recover(); r != nil {
				msg = configSavedMsg{err: fmt.Errorf("save failed: %v", r)}
			}
		}()

		if m.config == nil {
			return configSavedMsg{err: fmt.Errorf("no config loaded")}
		}
		m.syncToConfig()
		err := config.SaveNiriConfig(m.config)
		return configSavedMsg{err: err}
	}
}

// Key bindings
var (
	keyUp = key.NewBinding(
		key.WithKeys("up", "k"),
	)
	keyDown = key.NewBinding(
		key.WithKeys("down", "j"),
	)
	keyLeft = key.NewBinding(
		key.WithKeys("left", "h"),
	)
	keyRight = key.NewBinding(
		key.WithKeys("right", "l"),
	)
	keyToggle = key.NewBinding(
		key.WithKeys(" "),
	)
	keySave = key.NewBinding(
		key.WithKeys("s"),
	)
	keyReset = key.NewBinding(
		key.WithKeys("r"),
	)
)
