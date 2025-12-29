# Nirimatic

A Niri desktop environment installer and TUI configuration manager for Arch Linux/CachyOS.

![Eldritch Theme](https://github.com/eldritch-theme/eldritch)

## Features

- **TUI Configuration Manager**: Beautiful terminal interface for managing Niri settings
- **Eldritch Theme**: Gorgeous cosmic horror color scheme throughout
- **Service Dashboard**: Real-time status of Niri, Noctalia, and Stasis services
- **Smart Installer**: Detects existing packages and only installs what's missing
- **Config Backup**: Export and import your configuration with a single command

## Screenshots

*Coming soon*

## Installation

### From Source

```bash
git clone https://github.com/edellingham/nirimatic
cd nirimatic
make build
sudo make install
```

### From AUR (Coming Soon)

```bash
yay -S nirimatic
```

## Usage

```bash
# Run the TUI manager
nirimatic

# Run the installer (fresh install or update)
./installer/install.sh
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/k` | Navigate up |
| `↓/j` | Navigate down |
| `Enter` | Select item |
| `r` | Reload Niri config |
| `n` | Open Noctalia settings |
| `q` | Quit |
| `?` | Show help |

## Configuration

Nirimatic manages the following configuration files:

- `~/.config/niri/config.kdl` - Niri compositor settings
- `~/.config/wezterm/wezterm.lua` - Wezterm terminal
- `~/.config/gtk-3.0/` - GTK3 theming
- `~/.config/gtk-4.0/` - GTK4 theming
- `~/.config/qt5ct/` - Qt5 theming
- `~/.config/qt6ct/` - Qt6 theming
- `~/.config/btop/themes/eldritch.theme` - btop theme
- `~/.config/stasis/config.toml` - Stasis wallpaper manager

## Eldritch Color Palette

| Color | Hex | Name |
|-------|-----|------|
| Background | `#212337` | Sunken Depths Grey |
| Current Line | `#323449` | Shallow Depths Grey |
| Foreground | `#ebfafa` | Lighthouse White |
| Cyan | `#04d1f9` | Watery Tomb Blue |
| Green | `#37f499` | Great Old One Green |
| Purple | `#a48cf2` | Lovecraft Purple |
| Pink | `#f265b5` | Pustule Pink |
| Red | `#f16c75` | R'lyeh Red |
| Yellow | `#f1fc79` | Gold of Yuggoth |
| Orange | `#f7c67f` | Dreaming Orange |

## Development

```bash
# Run in development mode
make dev

# Build release binary
make release

# Run tests
make test

# Run linter
make lint
```

## Dependencies

- Go 1.22+
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components

## License

MIT License - See [LICENSE](LICENSE) for details.

## Acknowledgments

- [Eldritch Theme](https://github.com/eldritch-theme/eldritch) for the beautiful color palette
- [Charm](https://charm.sh) for the excellent TUI libraries
- [Niri](https://github.com/YaLTeR/niri) for the scrolling window manager
- [Noctalia](https://github.com/noctalia/noctalia-shell) for the shell integration
