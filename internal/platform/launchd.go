//go:build darwin

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/simpossible/mini_launch/internal/service"
)

const (
	launchdLabelPrefix = "com.minilaunch"
)

type launchd struct{}

func newPlatform() Platform {
	return &launchd{}
}

func plistPath(svc *service.Service) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", plistFilename(svc))
}

func plistFilename(svc *service.Service) string {
	return fmt.Sprintf("%s.%s.plist", launchdLabelPrefix, svc.Name)
}

func label(svc *service.Service) string {
	return fmt.Sprintf("%s.%s", launchdLabelPrefix, svc.Name)
}

func (l *launchd) Generate(svc *service.Service) error {
	path := plistPath(svc)

	// Ensure LaunchAgents directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create LaunchAgents directory: %w", err)
	}

	plist := buildPlist(svc)
	if err := os.WriteFile(path, []byte(plist), 0644); err != nil {
		return fmt.Errorf("cannot write plist file: %w", err)
	}

	fmt.Printf("Generated: %s\n", path)
	return nil
}

func (l *launchd) Start(svc *service.Service) error {
	path := plistPath(svc)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("service not configured. Run 'mini_launch initial' first")
	}

	// Unload first in case it's already loaded (ignore errors)
	exec.Command("launchctl", "unload", path).Run()

	cmd := exec.Command("launchctl", "load", "-w", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start service: %s\n%s", err, string(out))
	}

	fmt.Printf("Service '%s' started\n", svc.Name)
	return nil
}

func (l *launchd) Stop(svc *service.Service) error {
	path := plistPath(svc)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("service not configured")
	}

	cmd := exec.Command("launchctl", "unload", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop service: %s\n%s", err, string(out))
	}

	fmt.Printf("Service '%s' stopped\n", svc.Name)
	return nil
}

func (l *launchd) Restart(svc *service.Service) error {
	if err := l.Stop(svc); err != nil {
		// It's okay if stop fails (service might not be running)
		_ = err
	}
	return l.Start(svc)
}

func (l *launchd) Status(svc *service.Service) (string, error) {
	lbl := label(svc)
	cmd := exec.Command("launchctl", "list", lbl)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "not running (not loaded)", nil
	}

	output := string(out)
	if strings.Contains(output, "PID") {
		return "running", nil
	}
	if strings.Contains(output, "Last Exit Status") {
		return "stopped", nil
	}
	return "loaded", nil
}

func (l *launchd) Remove(svc *service.Service) error {
	// Try to stop first
	_ = l.Stop(svc)

	path := plistPath(svc)
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("cannot remove plist: %w", err)
		}
	}

	fmt.Printf("Service '%s' removed\n", svc.Name)
	return nil
}

func (l *launchd) IsConfigured(svc *service.Service) bool {
	_, err := os.Stat(plistPath(svc))
	return err == nil
}

func buildPlist(svc *service.Service) string {
	var envEntries strings.Builder
	if len(svc.EnvVars) > 0 {
		envEntries.WriteString("\t\t<key>EnvironmentVariables</key>\n")
		envEntries.WriteString("\t\t<dict>\n")
		for _, env := range svc.EnvVars {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				envEntries.WriteString(fmt.Sprintf("\t\t\t<key>%s</key>\n", parts[0]))
				envEntries.WriteString(fmt.Sprintf("\t\t\t<string>%s</string>\n", escapeXML(parts[1])))
			}
		}
		envEntries.WriteString("\t\t</dict>\n")
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>WorkingDirectory</key>
	<string>%s</string>
	<key>StandardOutPath</key>
	<string>%s</string>
	<key>StandardErrorPath</key>
	<string>%s</string>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
%s</dict>
</plist>
`,
		label(svc),
		escapeXML(svc.Executable),
		escapeXML(svc.Dir),
		escapeXML(svc.LogFile),
		escapeXML(svc.LogFile),
		envEntries.String(),
	)
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
