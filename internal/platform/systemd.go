//go:build linux

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/simpossible/mini_launch/internal/service"
)

type systemd struct{}

func newPlatform() Platform {
	return &systemd{}
}

// unitFilePath returns the path where the unit file is stored (in service directory).
func unitFilePath(svc *service.Service) string {
	return filepath.Join(svc.Dir, unitFilename(svc))
}

// unitLinkPath returns the symlink path in the systemd user directory.
func unitLinkPath(svc *service.Service) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "systemd", "user", unitFilename(svc))
}

func unitFilename(svc *service.Service) string {
	return fmt.Sprintf("mini-launch-%s.service", svc.Name)
}

func unitName(svc *service.Service) string {
	return fmt.Sprintf("mini-launch-%s.service", svc.Name)
}

func (s *systemd) Generate(svc *service.Service) error {
	filePath := unitFilePath(svc)
	linkPath := unitLinkPath(svc)

	// Write unit file to service directory
	unit := buildUnit(svc)
	if err := os.WriteFile(filePath, []byte(unit), 0644); err != nil {
		return fmt.Errorf("cannot write service file: %w", err)
	}

	// Ensure systemd user directory exists
	linkDir := filepath.Dir(linkPath)
	if err := os.MkdirAll(linkDir, 0755); err != nil {
		return fmt.Errorf("cannot create systemd user directory: %w", err)
	}

	// Remove existing symlink or file at link path
	os.Remove(linkPath)

	// Create symlink: systemd user dir -> service directory
	if err := os.Symlink(filePath, linkPath); err != nil {
		return fmt.Errorf("cannot create symlink in systemd directory: %w", err)
	}

	// Reload systemd daemon
	exec.Command("systemctl", "--user", "daemon-reload").Run()

	fmt.Printf("Generated: %s\n", filePath)
	fmt.Printf("Linked:    %s -> %s\n", linkPath, filePath)
	return nil
}

func (s *systemd) Start(svc *service.Service) error {
	if !s.IsConfigured(svc) {
		return fmt.Errorf("service not configured. Run 'mini_launch initial' first")
	}

	cmd := exec.Command("systemctl", "--user", "start", unitName(svc))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start service: %s\n%s", err, string(out))
	}

	fmt.Printf("Service '%s' started\n", svc.Name)
	return nil
}

func (s *systemd) Stop(svc *service.Service) error {
	cmd := exec.Command("systemctl", "--user", "stop", unitName(svc))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop service: %s\n%s", err, string(out))
	}

	fmt.Printf("Service '%s' stopped\n", svc.Name)
	return nil
}

func (s *systemd) Restart(svc *service.Service) error {
	cmd := exec.Command("systemctl", "--user", "restart", unitName(svc))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart service: %s\n%s", err, string(out))
	}

	fmt.Printf("Service '%s' restarted\n", svc.Name)
	return nil
}

func (s *systemd) Status(svc *service.Service) (string, error) {
	cmd := exec.Command("systemctl", "--user", "is-active", unitName(svc))
	out, _ := cmd.CombinedOutput()
	status := strings.TrimSpace(string(out))

	switch status {
	case "active":
		return "running", nil
	case "inactive", "stopped":
		return "stopped", nil
	case "failed":
		return "failed", nil
	default:
		return status, nil
	}
}

func (s *systemd) Remove(svc *service.Service) error {
	_ = s.Stop(svc)

	// Remove symlink from systemd user directory
	linkPath := unitLinkPath(svc)
	if _, err := os.Lstat(linkPath); err == nil {
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("cannot remove systemd symlink: %w", err)
		}
		exec.Command("systemctl", "--user", "daemon-reload").Run()
	}

	// Remove unit file from service directory
	filePath := unitFilePath(svc)
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("cannot remove service file: %w", err)
		}
	}

	fmt.Printf("Service '%s' removed\n", svc.Name)
	return nil
}

func (s *systemd) IsConfigured(svc *service.Service) bool {
	_, err := os.Stat(unitFilePath(svc))
	return err == nil
}

func buildUnit(svc *service.Service) string {
	var envLines strings.Builder
	if len(svc.EnvVars) > 0 {
		envLines.WriteString("Environment=")
		for i, env := range svc.EnvVars {
			if i > 0 {
				envLines.WriteString(" ")
			}
			envLines.WriteString(fmt.Sprintf("\"%s\"", env))
		}
		envLines.WriteString("\n")
	}

	return fmt.Sprintf(`[Unit]
Description=mini_launch service: %s
After=network.target

[Service]
Type=simple
WorkingDirectory=%s
ExecStart=%s
%sStandardOutput=append:%s
StandardError=append:%s
Restart=always
RestartSec=5

[Install]
WantedBy=default.target
`,
		svc.Name,
		svc.Dir,
		svc.Executable,
		envLines.String(),
		svc.LogFile,
		svc.LogFile,
	)
}
