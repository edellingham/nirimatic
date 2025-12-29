package tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Service represents a monitored service
type Service struct {
	Name   string
	Status string // "running", "stopped", "unknown"
}

// DashboardModel is the model for the dashboard screen
type DashboardModel struct {
	services    []Service
	width       int
	height      int
	lastRefresh time.Time
}

// serviceStatusMsg is sent when service status is updated
type serviceStatusMsg struct {
	services []Service
}

// tickMsg is sent periodically to refresh status
type tickMsg time.Time

// NewDashboardModel creates a new dashboard model
func NewDashboardModel() *DashboardModel {
	return &DashboardModel{
		services: []Service{
			{Name: "niri", Status: "unknown"},
			{Name: "noctalia-shell", Status: "unknown"},
			{Name: "stasis", Status: "unknown"},
		},
	}
}

// Init initializes the dashboard
func (m *DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.checkServices(),
		m.tick(),
	)
}

// SetSize sets the dashboard dimensions
func (m *DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages
func (m *DashboardModel) Update(msg tea.Msg) (*DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case serviceStatusMsg:
		m.services = msg.services
		m.lastRefresh = time.Now()
		return m, nil

	case tickMsg:
		return m, tea.Batch(m.checkServices(), m.tick())
	}

	return m, nil
}

// View renders the dashboard
func (m *DashboardModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render("Dashboard"))
	b.WriteString("\n")
	b.WriteString(SectionStyle.Render("─────────────────────────────────────────"))
	b.WriteString("\n\n")

	// Service status section
	b.WriteString(CardTitleStyle.Render("Service Status"))
	b.WriteString("\n\n")

	for _, svc := range m.services {
		status := RenderStatus(svc.Status)
		name := lipgloss.NewStyle().Width(20).Render(svc.Name)
		statusText := m.getStatusText(svc.Status)
		b.WriteString(fmt.Sprintf("  %s %s  %s\n", status, name, statusText))
	}

	b.WriteString("\n")

	// Quick Actions section
	b.WriteString(SectionStyle.Render("─────────────────────────────────────────"))
	b.WriteString("\n\n")
	b.WriteString(CardTitleStyle.Render("Quick Actions"))
	b.WriteString("\n\n")

	actions := []struct {
		key  string
		desc string
	}{
		{"r", "Reload Niri Config"},
		{"n", "Noctalia Settings"},
		{"R", "Restart Niri"},
		{"q", "Quit"},
	}

	for _, action := range actions {
		keyStyle := ButtonStyle.Render(action.key)
		b.WriteString(fmt.Sprintf("  %s %s\n", keyStyle, action.desc))
	}

	b.WriteString("\n")

	// Last refresh time
	if !m.lastRefresh.IsZero() {
		refreshText := DimmedStyle.Render(fmt.Sprintf("Last updated: %s", m.lastRefresh.Format("15:04:05")))
		b.WriteString(refreshText)
	}

	return b.String()
}

// getStatusText returns a human-readable status text
func (m *DashboardModel) getStatusText(status string) string {
	switch status {
	case "running":
		return SuccessStyle.Render("Running")
	case "stopped":
		return ErrorStyle.Render("Stopped")
	default:
		return DimmedStyle.Render("Unknown")
	}
}

// checkServices checks the status of all services
func (m *DashboardModel) checkServices() tea.Cmd {
	return func() tea.Msg {
		services := make([]Service, len(m.services))
		for i, svc := range m.services {
			services[i] = Service{
				Name:   svc.Name,
				Status: checkServiceStatus(svc.Name),
			}
		}
		return serviceStatusMsg{services: services}
	}
}

// tick returns a command that sends a tick message after a delay
func (m *DashboardModel) tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// checkServiceStatus checks if a service is running
func checkServiceStatus(name string) string {
	// For niri, check if the compositor is running via niri msg
	if name == "niri" {
		cmd := exec.Command("niri", "msg", "version")
		if err := cmd.Run(); err == nil {
			return "running"
		}
		return "stopped"
	}

	// For other services, check systemd user services
	cmd := exec.Command("systemctl", "--user", "is-active", name)
	output, _ := cmd.Output()
	status := strings.TrimSpace(string(output))

	switch status {
	case "active":
		return "running"
	case "inactive", "failed":
		return "stopped"
	default:
		// Also check if process is running directly
		cmd := exec.Command("pgrep", "-x", name)
		if err := cmd.Run(); err == nil {
			return "running"
		}
		return "unknown"
	}
}
