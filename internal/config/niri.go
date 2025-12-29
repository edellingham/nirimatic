package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// NiriConfig holds the parsed Niri configuration
type NiriConfig struct {
	Path string

	// Layout settings
	Gaps         int
	BorderWidth  int
	FocusRingWidth int
	CornerRadius int

	// Shadow settings
	ShadowEnabled  bool
	ShadowSoftness int
	ShadowSpread   int

	// Behavior settings
	FocusFollowsMouse       bool
	WorkspaceAutoBackAndForth bool
}

// DefaultNiriConfig returns a config with default values
func DefaultNiriConfig() *NiriConfig {
	return &NiriConfig{
		Gaps:                      10,
		BorderWidth:               2,
		FocusRingWidth:            0,
		CornerRadius:              16,
		ShadowEnabled:             true,
		ShadowSoftness:            60,
		ShadowSpread:              10,
		FocusFollowsMouse:         true,
		WorkspaceAutoBackAndForth: true,
	}
}

// GetConfigPath returns the default Niri config path
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "niri", "config.kdl")
}

// LoadNiriConfig loads and parses the Niri configuration
func LoadNiriConfig(path string) (*NiriConfig, error) {
	config := DefaultNiriConfig()
	config.Path = path

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inLayout := false
	inShadow := false
	inInput := false
	inBorder := false
	inFocusRing := false
	braceDepth := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Track brace depth
		braceDepth += strings.Count(line, "{") - strings.Count(line, "}")

		// Track sections
		if strings.HasPrefix(line, "layout") && strings.Contains(line, "{") {
			inLayout = true
			continue
		}
		if strings.HasPrefix(line, "input") && strings.Contains(line, "{") {
			inInput = true
			continue
		}
		if inLayout && strings.HasPrefix(line, "shadow") && strings.Contains(line, "{") {
			inShadow = true
			continue
		}
		if inLayout && strings.HasPrefix(line, "border") && strings.Contains(line, "{") {
			inBorder = true
			continue
		}
		if inLayout && strings.HasPrefix(line, "focus-ring") && strings.Contains(line, "{") {
			inFocusRing = true
			continue
		}

		// Exit sections on closing brace
		if line == "}" {
			if inShadow {
				inShadow = false
			} else if inBorder {
				inBorder = false
			} else if inFocusRing {
				inFocusRing = false
			} else if inLayout && braceDepth == 0 {
				inLayout = false
			} else if inInput && braceDepth == 0 {
				inInput = false
			}
			continue
		}

		// Parse input section
		if inInput {
			if line == "focus-follows-mouse" {
				config.FocusFollowsMouse = true
			}
			if line == "workspace-auto-back-and-forth" {
				config.WorkspaceAutoBackAndForth = true
			}
		}

		// Parse layout section
		if inLayout && !inShadow && !inBorder && !inFocusRing {
			if val := parseIntValue(line, "gaps"); val != nil {
				config.Gaps = *val
			}
		}

		// Parse shadow section
		if inShadow {
			if line == "on" {
				config.ShadowEnabled = true
			}
			if line == "off" {
				config.ShadowEnabled = false
			}
			if val := parseIntValue(line, "softness"); val != nil {
				config.ShadowSoftness = *val
			}
			if val := parseIntValue(line, "spread"); val != nil {
				config.ShadowSpread = *val
			}
		}

		// Parse border section
		if inBorder {
			if val := parseIntValue(line, "width"); val != nil {
				config.BorderWidth = *val
			}
		}

		// Parse focus-ring section
		if inFocusRing {
			if val := parseIntValue(line, "width"); val != nil {
				config.FocusRingWidth = *val
			}
		}

		// Parse window-rule for corner-radius
		if strings.Contains(line, "geometry-corner-radius") {
			if val := parseIntValue(line, "geometry-corner-radius"); val != nil {
				config.CornerRadius = *val
			}
		}
	}

	return config, scanner.Err()
}

// parseIntValue extracts an integer value from a KDL property line
func parseIntValue(line, key string) *int {
	// Match patterns like "gaps 10" or "width 2"
	pattern := regexp.MustCompile(fmt.Sprintf(`^\s*%s\s+(\d+)`, regexp.QuoteMeta(key)))
	matches := pattern.FindStringSubmatch(line)
	if len(matches) >= 2 {
		val, err := strconv.Atoi(matches[1])
		if err == nil {
			return &val
		}
	}
	return nil
}

// SaveNiriConfig saves the configuration back to the file
// This performs surgical updates to preserve comments and formatting
func SaveNiriConfig(config *NiriConfig) error {
	content, err := os.ReadFile(config.Path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	inLayout := false
	inShadow := false
	inBorder := false
	inFocusRing := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track sections
		if strings.HasPrefix(trimmed, "layout") && strings.Contains(trimmed, "{") {
			inLayout = true
			continue
		}
		if inLayout && strings.HasPrefix(trimmed, "shadow") && strings.Contains(trimmed, "{") {
			inShadow = true
			continue
		}
		if inLayout && strings.HasPrefix(trimmed, "border") && strings.Contains(trimmed, "{") {
			inBorder = true
			continue
		}
		if inLayout && strings.HasPrefix(trimmed, "focus-ring") && strings.Contains(trimmed, "{") {
			inFocusRing = true
			continue
		}

		// Handle closing braces
		if trimmed == "}" {
			if inShadow {
				inShadow = false
			} else if inBorder {
				inBorder = false
			} else if inFocusRing {
				inFocusRing = false
			} else if inLayout {
				inLayout = false
			}
			continue
		}

		// Update values
		if inLayout && !inShadow && !inBorder && !inFocusRing {
			lines[i] = updateIntValue(line, "gaps", config.Gaps)
		}
		if inShadow {
			lines[i] = updateIntValue(line, "softness", config.ShadowSoftness)
			lines[i] = updateIntValue(lines[i], "spread", config.ShadowSpread)
			// Handle on/off
			if trimmed == "on" || trimmed == "off" {
				indent := line[:len(line)-len(trimmed)]
				if config.ShadowEnabled {
					lines[i] = indent + "on"
				} else {
					lines[i] = indent + "off"
				}
			}
		}
		if inBorder {
			lines[i] = updateIntValue(line, "width", config.BorderWidth)
		}
		if inFocusRing {
			lines[i] = updateIntValue(line, "width", config.FocusRingWidth)
		}

		// Update corner radius in window-rule
		if strings.Contains(trimmed, "geometry-corner-radius") {
			lines[i] = updateIntValue(line, "geometry-corner-radius", config.CornerRadius)
		}
	}

	return os.WriteFile(config.Path, []byte(strings.Join(lines, "\n")), 0644)
}

// updateIntValue updates an integer value in a line while preserving formatting
func updateIntValue(line, key string, value int) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(\s*%s\s+)\d+`, regexp.QuoteMeta(key)))
	return pattern.ReplaceAllString(line, fmt.Sprintf("${1}%d", value))
}
