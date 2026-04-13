package service

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var exportRe = regexp.MustCompile(`^\s*export\s+([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)

// CollectEnvVars scans the user's shell config files for exported environment variables.
func CollectEnvVars() []string {
	shell := detectShell()
	var configFile string

	switch {
	case strings.Contains(shell, "zsh"):
		configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
	case strings.Contains(shell, "bash"):
		configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
	default:
		// Try both
		home := os.Getenv("HOME")
		if _, err := os.Stat(filepath.Join(home, ".zshrc")); err == nil {
			configFile = filepath.Join(home, ".zshrc")
		} else if _, err := os.Stat(filepath.Join(home, ".bashrc")); err == nil {
			configFile = filepath.Join(home, ".bashrc")
		} else {
			return nil
		}
	}

	return parseExports(configFile)
}

func detectShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}
	if runtime.GOOS == "darwin" {
		return "/bin/zsh"
	}
	return "/bin/bash"
}

func parseExports(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var envVars []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		matches := exportRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		name := matches[1]
		value := matches[2]

		// Remove surrounding quotes
		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Skip values that reference other variables (like $PATH) — these won't expand correctly
		if strings.Contains(value, "$") {
			continue
		}

		envVars = append(envVars, fmt.Sprintf("%s=%s", name, value))
	}

	return envVars
}
