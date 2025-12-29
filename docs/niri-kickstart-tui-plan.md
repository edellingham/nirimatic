# Nirimatic - Project Planning Document (TUI Version)

## Project Name: `nirimatic`

A complete Niri desktop environment installer and TUI configuration manager for Arch Linux.

---

## Architecture: Charm-Based TUI

### Why TUI over GUI?

| Aspect | TUI (Bubbletea) | GUI (Tauri) |
|--------|-----------------|-------------|
| Binary size | ~5-10MB | ~15-30MB |
| Dependencies | None (static Go binary) | WebKit runtime |
| Distribution | Single binary, AUR | AppImage |
| Feels native to | Terminal users, WM users | Desktop users |
| Development speed | Fast (Go) | Medium (Rust + JS) |
| Theming | Lipgloss (Eldritch!) | CSS |

**TUI wins** because:
- Niri users are terminal-comfortable
- Single static binary - just copy and run
- AUR package is trivial
- Charm tools look gorgeous
- Can use Eldritch colors natively in Lipgloss

### Tool Selection

| Tool | Use Case |
|------|----------|
| **Bubbletea** | Main TUI application framework |
| **Lipgloss** | Styling, borders, colors (Eldritch theme!) |
| **Bubbles** | Pre-built components (lists, text inputs, spinners) |
| **Gum** | Installer script interactions |

**Why not Gum for everything?**
Gum is excellent for shell script interactions but not ideal for a persistent, stateful config manager. Bubbletea gives us:
- Persistent state across screens
- Real-time file watching
- Complex navigation (tabs, nested menus)
- Live preview of changes

**Gum is perfect for the installer** - beautiful prompts without writing Go.

---

## Project Structure

```
nirimatic/
├── README.md
├── LICENSE
├── Makefile                      # Build targets
├── go.mod
├── go.sum
│
├── cmd/
│   └── nirimatic/
│       └── main.go               # Entry point
│
├── internal/
│   ├── tui/
│   │   ├── app.go                # Main Bubbletea app
│   │   ├── styles.go             # Lipgloss Eldritch theme
│   │   ├── keys.go               # Keybindings
│   │   │
│   │   ├── screens/
│   │   │   ├── dashboard.go      # Home screen
│   │   │   ├── niri.go           # Niri config editor
│   │   │   ├── animations.go     # Animation tuner
│   │   │   ├── keybinds.go       # Keybind editor
│   │   │   ├── startup.go        # Startup apps
│   │   │   └── backup.go         # Export/import
│   │   │
│   │   └── components/
│   │       ├── header.go         # App header
│   │       ├── sidebar.go        # Navigation
│   │       ├── editor.go         # KDL value editor
│   │       ├── toggle.go         # Boolean toggle
│   │       └── slider.go         # Numeric slider
│   │
│   ├── config/
│   │   ├── niri.go               # Niri config parsing
│   │   ├── kdl.go                # KDL read/write
│   │   ├── export.go             # Backup export
│   │   └── import.go             # Backup import
│   │
│   └── system/
│       ├── services.go           # systemd control
│       ├── packages.go           # pacman/yay queries
│       └── niri_ipc.go           # Niri socket communication
│
├── installer/
│   ├── install.sh                # Main installer (uses Gum)
│   ├── packages.sh
│   ├── aur.sh
│   ├── services.sh
│   ├── sddm.sh
│   ├── eldritch.sh
│   └── post-install.sh
│
├── configs/                      # Config templates
│   ├── niri/
│   │   ├── config.kdl
│   │   └── noctalia.kdl
│   ├── wezterm/
│   │   └── wezterm.lua
│   ├── gtk-3.0/
│   │   ├── settings.ini
│   │   └── gtk.css
│   ├── gtk-4.0/
│   │   ├── settings.ini
│   │   └── gtk.css
│   ├── qt5ct/
│   │   └── qt5ct.conf
│   ├── qt6ct/
│   │   └── qt6ct.conf
│   ├── btop/
│   │   └── themes/
│   │       └── eldritch.theme
│   ├── stasis/
│   │   └── config.toml
│   ├── environment.d/
│   │   └── wayland.conf
│   └── sddm/
│       └── sddm.conf
│
├── scripts/
│   ├── screenshot-gradia.sh
│   ├── export-config.sh
│   └── import-config.sh
│
└── docs/
    ├── INSTALLATION.md
    └── CONFIGURATION.md
```

---

## Lipgloss Eldritch Theme

```go
// internal/tui/styles.go
package tui

import "github.com/charmbracelet/lipgloss"

// Eldritch color palette
var (
    ColorBackground   = lipgloss.Color("#212337") // Sunken Depths Grey
    ColorCurrentLine  = lipgloss.Color("#323449") // Shallow Depths Grey
    ColorForeground   = lipgloss.Color("#ebfafa") // Lighthouse White
    ColorComment      = lipgloss.Color("#7081d0") // The Old One Purple
    ColorCyan         = lipgloss.Color("#04d1f9") // Watery Tomb Blue
    ColorGreen        = lipgloss.Color("#37f499") // Great Old One Green
    ColorOrange       = lipgloss.Color("#f7c67f") // Dreaming Orange
    ColorPink         = lipgloss.Color("#f265b5") // Pustule Pink
    ColorPurple       = lipgloss.Color("#a48cf2") // Lovecraft Purple
    ColorRed          = lipgloss.Color("#f16c75") // R'lyeh Red
    ColorYellow       = lipgloss.Color("#f1fc79") // Gold of Yuggoth
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
        PaddingLeft(1).
        SetString("▸ ")

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

    // Status indicators
    StatusOK = lipgloss.NewStyle().
        SetString("●").
        Foreground(ColorGreen)

    StatusError = lipgloss.NewStyle().
        SetString("●").
        Foreground(ColorRed)

    StatusWarning = lipgloss.NewStyle().
        SetString("●").
        Foreground(ColorOrange)

    // Gradient-style title (like your niri borders!)
    GradientTitle = lipgloss.NewStyle().
        Bold(true).
        Foreground(ColorPink). // Start color
        // Note: Lipgloss doesn't do gradients, but we can alternate
)

// Helper function for gradient-like effect on text
func GradientText(text string) string {
    colors := []lipgloss.Color{ColorPink, ColorPurple, ColorCyan}
    result := ""
    for i, char := range text {
        color := colors[i%len(colors)]
        result += lipgloss.NewStyle().Foreground(color).Render(string(char))
    }
    return result
}
```

---

## TUI Screens

### Main Layout
```
┌─────────────────────────────────────────────────────────────────────┐
│  ▄▄ nirimatic                                              v1.0.0   │
├──────────────────────┬──────────────────────────────────────────────┤
│                      │                                              │
│   ◆ Dashboard       │   Dashboard                                  │
│     Niri Settings    │   ─────────────────────────────────────────  │
│     Animations       │                                              │
│     Keybinds         │   ● Niri         Running                     │
│     Startup Apps     │   ● Noctalia     Running                     │
│     Backup           │   ● Stasis       Running                     │
│                      │                                              │
│   ─────────────────  │   Quick Actions                              │
│   Noctalia Settings  │   ─────────────────────────────────────────  │
│   Reload Niri        │   [r] Reload Config  [n] Noctalia Settings   │
│   Quit               │   [s] Restart Niri   [q] Quit                │
│                      │                                              │
├──────────────────────┴──────────────────────────────────────────────┤
│  ↑↓ navigate  enter select  q quit  ? help                         │
└─────────────────────────────────────────────────────────────────────┘
```

### Screen Breakdown

#### 1. Dashboard
- Service status (Niri, Noctalia, Stasis)
- Quick actions with hotkeys
- Recent config changes (if we track them)
- Link to open Noctalia settings (spawns the panel)

#### 2. Niri Settings
```
┌─────────────────────────────────────────────────────────────────────┐
│  Niri Settings                                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Layout                                                             │
│  ─────────────────────────────────────────────────────────────────  │
│  Gaps                    [    10    ] px                            │
│  Border Width            [     2    ] px                            │
│  Focus Ring Width        [     0    ] px                            │
│  Corner Radius           [    16    ] px                            │
│                                                                     │
│  Shadows                 [✓] Enabled                                │
│  Shadow Softness         [────●─────] 60                            │
│  Shadow Spread           [──●───────] 10                            │
│                                                                     │
│  Behavior                                                           │
│  ─────────────────────────────────────────────────────────────────  │
│  Focus Follows Mouse     [✓] Enabled                                │
│  Workspace Auto Back     [✓] Enabled                                │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  tab section  ↑↓ navigate  enter edit  s save  esc cancel          │
└─────────────────────────────────────────────────────────────────────┘
```

#### 3. Animations
```
┌─────────────────────────────────────────────────────────────────────┐
│  Animations                                                         │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Presets                                                            │
│  ─────────────────────────────────────────────────────────────────  │
│  ( ) Default    ( ) Smooth    (●) Snappy    ( ) None                │
│                                                                     │
│  Fine Tuning                                                        │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  ▼ workspace-switch                                                 │
│    Type: Spring                                                     │
│    Damping Ratio    [────────●─] 1.0                                │
│    Stiffness        [───────●──] 1000                               │
│    Epsilon          [●─────────] 0.0001                             │
│                                                                     │
│  ▶ window-open                                                      │
│  ▶ window-close                                                     │
│  ▶ horizontal-view-movement                                         │
│  ▶ window-resize                                                    │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  space expand  ↑↓ navigate  ←→ adjust  s save  r reset              │
└─────────────────────────────────────────────────────────────────────┘
```

#### 4. Keybinds
```
┌─────────────────────────────────────────────────────────────────────┐
│  Keybinds                                                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Filter: [________________]                    [a] Add  [d] Delete  │
│                                                                     │
│  ┌────────────────────────┬─────────────────────────────────────┐   │
│  │ Keybind                │ Action                              │   │
│  ├────────────────────────┼─────────────────────────────────────┤   │
│  │ Mod+T                  │ spawn "wezterm"                     │   │
│  │ Mod+B                  │ spawn "helium-browser"              │   │
│  │ Mod+Space              │ qs -c noctalia-shell ipc call...    │   │
│  │ Mod+Q                  │ close-window                        │   │
│  │ Mod+W                  │ toggle-window-floating              │   │
│  │ Mod+M                  │ maximize-column                     │   │
│  │ Mod+Left               │ focus-column-left                   │   │
│  │ Mod+Right              │ focus-column-right                  │   │
│  │ ...                    │ ...                                 │   │
│  └────────────────────────┴─────────────────────────────────────┘   │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  / filter  enter edit  a add  d delete  s save                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### 5. Startup Apps
```
┌─────────────────────────────────────────────────────────────────────┐
│  Startup Applications                                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  spawn-at-startup entries                      [a] Add  [d] Delete  │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  [✓] qs -c noctalia-shell                      Shell                │
│  [✓] /usr/lib/polkit-kde-authentication-ag...  Polkit               │
│  [✓] xwayland-satellite                        XWayland             │
│  [✓] pumble-desktop                            Chat                 │
│  [ ] discord                                   (disabled)           │
│                                                                     │
│  systemd user services                                              │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  ● stasis.service                              Active               │
│  ○ niri-session-manager.service                Inactive             │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  space toggle  a add  d delete  e enable service  s save            │
└─────────────────────────────────────────────────────────────────────┘
```

#### 6. Backup
```
┌─────────────────────────────────────────────────────────────────────┐
│  Backup & Restore                                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Export Configuration                                               │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  Include:                                                           │
│  [✓] Niri config                                                    │
│  [✓] Noctalia settings                                              │
│  [✓] Wezterm config                                                 │
│  [✓] GTK settings                                                   │
│  [✓] Qt settings                                                    │
│  [✓] Stasis config                                                  │
│  [ ] Monitor configuration  ← Machine-specific, off by default      │
│                                                                     │
│  [        Export to ~/niri-backup.tar.gz        ]                   │
│                                                                     │
│  ─────────────────────────────────────────────────────────────────  │
│  Import Configuration                                               │
│                                                                     │
│  [        Select backup file...                 ]                   │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  space toggle  enter action  esc back                               │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Installer with Gum

The installer uses Gum for beautiful prompts while remaining a bash script:

```bash
#!/bin/bash
# installer/install.sh

set -euo pipefail

# Check for gum, install if missing
if ! command -v gum &> /dev/null; then
    echo "Installing gum for pretty prompts..."
    sudo pacman -S --noconfirm gum
fi

# Header
gum style \
    --foreground "#04d1f9" \
    --border-foreground "#a48cf2" \
    --border double \
    --align center \
    --width 50 \
    --margin "1 2" \
    --padding "1 2" \
    '▄▄ nirimatic' \
    'Niri Desktop Environment Installer'

echo ""

# Confirm installation
gum confirm "Ready to install Nirimatic?" || exit 0

echo ""
gum style --foreground "#7081d0" "This will install:"
echo ""
gum style --foreground "#ebfafa" "  • Niri compositor"
gum style --foreground "#ebfafa" "  • Noctalia shell"
gum style --foreground "#ebfafa" "  • Eldritch theme configs"
gum style --foreground "#ebfafa" "  • Supporting tools (stasis, xwayland-satellite, etc.)"
gum style --foreground "#ebfafa" "  • SDDM with Astronaut theme"
gum style --foreground "#ebfafa" "  • Default applications (Wezterm, Helium, Gradia, etc.)"
echo ""

# Component selection
gum style --foreground "#f265b5" "Select components to install:"
COMPONENTS=$(gum choose --no-limit --selected="Core,Shell,Apps,Theme,SDDM" \
    "Core" \
    "Shell" \
    "Apps" \
    "Theme" \
    "SDDM" \
    "TUI Manager")

echo ""

# AUR helper check
if ! command -v yay &> /dev/null; then
    gum style --foreground "#f7c67f" "⚠ yay not found. Installing..."
    # Install yay
    gum spin --spinner dot --title "Installing yay..." -- \
        bash -c 'git clone https://aur.archlinux.org/yay-bin.git /tmp/yay-bin && \
                 cd /tmp/yay-bin && makepkg -si --noconfirm && rm -rf /tmp/yay-bin'
fi

# Install packages with spinner
gum style --foreground "#37f499" "Installing packages..."
echo ""

# Pacman packages
gum spin --spinner dot --title "Installing core packages..." -- \
    sudo pacman -S --noconfirm --needed \
    niri xwayland-satellite xdg-desktop-portal-gnome xdg-desktop-portal-gtk \
    pipewire wireplumber playerctl pamixer \
    polkit-gnome gnome-keyring seahorse \
    wl-clipboard cliphist grim slurp \
    qt5ct qt6ct kvantum nwg-look \
    brightnessctl wezterm nautilus \
    ttf-jetbrains-mono-nerd fzf bat fastfetch btop

# AUR packages
gum spin --spinner dot --title "Installing AUR packages..." -- \
    yay -S --noconfirm --needed \
    quickshell-git noctalia-shell-git stasis-git nirimation-git \
    niri-scratchpad-rs-git helium-browser-bin gradia

# SDDM if selected
if [[ "$COMPONENTS" == *"SDDM"* ]]; then
    gum spin --spinner dot --title "Setting up SDDM..." -- \
        bash -c 'sudo pacman -S --noconfirm --needed sddm qt6-svg qt6-virtualkeyboard qt6-multimedia-ffmpeg && \
                 yay -S --noconfirm sddm-astronaut-theme'
fi

# Deploy configs
gum style --foreground "#37f499" "Deploying configuration files..."
gum spin --spinner dot --title "Copying configs..." -- \
    bash ./installer/post-install.sh

# Success!
echo ""
gum style \
    --foreground "#212337" \
    --background "#37f499" \
    --bold \
    --padding "1 2" \
    "✓ Installation complete!"

echo ""
gum style --foreground "#ebfafa" "Next steps:"
gum style --foreground "#04d1f9" "  1. Log out"
gum style --foreground "#04d1f9" "  2. Select 'Niri' from SDDM"
gum style --foreground "#04d1f9" "  3. Run 'nirimatic' to configure"
echo ""

# Offer to install TUI manager
if [[ "$COMPONENTS" == *"TUI Manager"* ]]; then
    gum style --foreground "#a48cf2" "The TUI manager will be available as 'nirimatic' command"
fi
```

---

## Go Application Structure

### Main Entry Point

```go
// cmd/nirimatic/main.go
package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/nirimatic/internal/tui"
)

func main() {
    p := tea.NewProgram(
        tui.NewApp(),
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v", err)
        os.Exit(1)
    }
}
```

### App Model

```go
// internal/tui/app.go
package tui

import (
    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type screen int

const (
    screenDashboard screen = iota
    screenNiri
    screenAnimations
    screenKeybinds
    screenStartup
    screenBackup
)

type App struct {
    currentScreen screen
    sidebar       list.Model
    width         int
    height        int
    
    // Screen models
    dashboard     DashboardModel
    niri          NiriModel
    animations    AnimationsModel
    keybinds      KeybindsModel
    startup       StartupModel
    backup        BackupModel
    
    // Config state
    configPath    string
    configDirty   bool
}

func NewApp() App {
    // Initialize sidebar
    items := []list.Item{
        sidebarItem{title: "Dashboard", screen: screenDashboard},
        sidebarItem{title: "Niri Settings", screen: screenNiri},
        sidebarItem{title: "Animations", screen: screenAnimations},
        sidebarItem{title: "Keybinds", screen: screenKeybinds},
        sidebarItem{title: "Startup Apps", screen: screenStartup},
        sidebarItem{title: "Backup", screen: screenBackup},
    }
    
    sidebar := list.New(items, sidebarDelegate{}, 20, 14)
    sidebar.Title = "nirimatic"
    sidebar.SetShowStatusBar(false)
    sidebar.SetFilteringEnabled(false)
    sidebar.Styles.Title = SidebarTitleStyle
    
    return App{
        currentScreen: screenDashboard,
        sidebar:       sidebar,
        configPath:    os.ExpandEnv("$HOME/.config/niri/config.kdl"),
        dashboard:     NewDashboardModel(),
        niri:          NewNiriModel(),
        animations:    NewAnimationsModel(),
        keybinds:      NewKeybindsModel(),
        startup:       NewStartupModel(),
        backup:        NewBackupModel(),
    }
}

func (a App) Init() tea.Cmd {
    return tea.Batch(
        a.dashboard.Init(),
        loadConfig(a.configPath),
    )
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            if a.configDirty {
                // Prompt to save
                return a, nil
            }
            return a, tea.Quit
        case "n":
            // Open Noctalia settings
            return a, openNoctaliaSettings()
        case "r":
            // Reload niri config
            return a, reloadNiriConfig()
        }

    case tea.WindowSizeMsg:
        a.width = msg.Width
        a.height = msg.Height
        
    case configLoadedMsg:
        // Distribute config to all screens
        a.niri = a.niri.SetConfig(msg.config)
        a.animations = a.animations.SetConfig(msg.config)
        a.keybinds = a.keybinds.SetConfig(msg.config)
        a.startup = a.startup.SetConfig(msg.config)
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
    case screenDashboard:
        a.dashboard, cmd = a.dashboard.Update(msg)
    case screenNiri:
        a.niri, cmd = a.niri.Update(msg)
    case screenAnimations:
        a.animations, cmd = a.animations.Update(msg)
    case screenKeybinds:
        a.keybinds, cmd = a.keybinds.Update(msg)
    case screenStartup:
        a.startup, cmd = a.startup.Update(msg)
    case screenBackup:
        a.backup, cmd = a.backup.Update(msg)
    }
    cmds = append(cmds, cmd)

    return a, tea.Batch(cmds...)
}

func (a App) View() string {
    // Render sidebar
    sidebar := SidebarStyle.Render(a.sidebar.View())

    // Render current screen content
    var content string
    switch a.currentScreen {
    case screenDashboard:
        content = a.dashboard.View()
    case screenNiri:
        content = a.niri.View()
    case screenAnimations:
        content = a.animations.View()
    case screenKeybinds:
        content = a.keybinds.View()
    case screenStartup:
        content = a.startup.View()
    case screenBackup:
        content = a.backup.View()
    }

    contentWidth := a.width - 28 // Sidebar width + padding
    contentBox := ContentStyle.Width(contentWidth).Render(content)

    // Layout
    main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, contentBox)

    // Header
    header := HeaderStyle.Width(a.width - 4).Render(
        GradientText("▄▄ nirimatic") + "  v1.0.0",
    )

    // Footer/Help
    help := HelpStyle.Render("↑↓ navigate • enter select • n noctalia • r reload • q quit • ? help")

    return lipgloss.JoinVertical(lipgloss.Left, header, main, help)
}
```

---

## KDL Parsing

For KDL parsing in Go, we can use:
- https://github.com/sblinch/kdl-go
- Or parse manually (KDL is simple enough)

```go
// internal/config/kdl.go
package config

import (
    "os"
    "github.com/sblinch/kdl-go"
)

type NiriConfig struct {
    Input      InputConfig
    Layout     LayoutConfig
    Animations []AnimationConfig
    Binds      []BindConfig
    Spawns     []SpawnConfig
    // ...
}

func LoadNiriConfig(path string) (*NiriConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    doc, err := kdl.Parse(string(data))
    if err != nil {
        return nil, err
    }
    
    config := &NiriConfig{}
    
    for _, node := range doc.Nodes {
        switch node.Name {
        case "input":
            config.Input = parseInput(node)
        case "layout":
            config.Layout = parseLayout(node)
        case "animations":
            config.Animations = parseAnimations(node)
        case "binds":
            config.Binds = parseBinds(node)
        case "spawn-at-startup":
            config.Spawns = append(config.Spawns, parseSpawn(node))
        }
    }
    
    return config, nil
}

func SaveNiriConfig(path string, config *NiriConfig) error {
    // Regenerate KDL from config
    // Preserve comments where possible
    // ...
}
```

---

## Distribution

### Single Binary
```bash
# Build
go build -ldflags="-s -w" -o nirimatic ./cmd/nirimatic

# Or with Makefile
make build
```

### AUR Package (PKGBUILD)
```bash
# Maintainer: Your Name <email>
pkgname=nirimatic
pkgver=1.0.0
pkgrel=1
pkgdesc="Niri desktop environment installer and TUI configuration manager"
arch=('x86_64')
url="https://github.com/yourusername/nirimatic"
license=('MIT')
depends=('niri' 'gum')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "$pkgname-$pkgver"
    go build -ldflags="-s -w" -o "$pkgname" ./cmd/nirimatic
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
    
    # Install configs
    install -Dm644 configs/niri/config.kdl "$pkgdir/usr/share/$pkgname/configs/niri/config.kdl"
    # ... more config files
    
    # Install installer scripts
    install -Dm755 installer/install.sh "$pkgdir/usr/share/$pkgname/install.sh"
}
```

---

## Development Workflow

```bash
# Clone
git clone https://github.com/yourusername/nirimatic
cd nirimatic

# Install Go dependencies
go mod download

# Run in development
go run ./cmd/nirimatic

# Build
make build

# Run installer (uses gum)
./installer/install.sh

# Build AUR package locally
makepkg -si
```

---

## Phase Plan for Claude Code

### Phase 1: Project Setup & Installer
1. Initialize Go module
2. Create directory structure
3. Write installer shell scripts with Gum
4. Config file templates

### Phase 2: TUI Foundation
1. Lipgloss Eldritch theme (styles.go)
2. Main app model with sidebar navigation
3. Dashboard screen with service status

### Phase 3: Config Screens
1. KDL parser/writer
2. Niri settings screen
3. Animations screen with presets
4. Keybinds editor

### Phase 4: Utilities
1. Startup apps manager
2. Export/import with monitor toggle
3. Niri IPC integration

### Phase 5: Polish & Distribution
1. Help system
2. Error handling
3. AUR PKGBUILD
4. Documentation

---

## References

- Bubbletea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss
- Bubbles: https://github.com/charmbracelet/bubbles
- Gum: https://github.com/charmbracelet/gum
- KDL Go: https://github.com/sblinch/kdl-go
- Eldritch Theme: https://github.com/eldritch-theme/eldritch
