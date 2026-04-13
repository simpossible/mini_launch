package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ServersDir = "servers"
)

type Service struct {
	Name       string   // service name (relative path with / replaced by _)
	Dir        string   // absolute path to service directory
	Executable string   // absolute path to executable
	LogFile    string   // absolute path to std.log
	EnvVars    []string // environment variables from shell config
}

// ServersBase returns the base directory $HOME/servers
func ServersBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ServersDir), nil
}

// NameFromDir computes the service name from a directory path relative to $HOME/servers.
// Path separators are replaced with underscores.
func NameFromDir(relPath string) string {
	relPath = strings.TrimPrefix(relPath, string(filepath.Separator))
	relPath = strings.TrimSuffix(relPath, string(filepath.Separator))
	return strings.ReplaceAll(relPath, string(filepath.Separator), "_")
}

// ServiceFromName reconstructs service info from a service name.
func ServiceFromName(name string) (*Service, error) {
	base, err := ServersBase()
	if err != nil {
		return nil, err
	}

	// name uses _ as separator, convert back to path separator
	relPath := strings.ReplaceAll(name, "_", string(filepath.Separator))
	dir := filepath.Join(base, relPath)

	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("service directory not found: %s", dir)
	}

	exec, err := findExecutable(dir)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:       name,
		Dir:        dir,
		Executable: exec,
		LogFile:    filepath.Join(dir, "std.log"),
		EnvVars:    CollectEnvVars(),
	}, nil
}

// ServiceFromDir creates a Service from a directory path.
func ServiceFromDir(dir string) (*Service, error) {
	base, err := ServersBase()
	if err != nil {
		return nil, err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve path: %w", err)
	}

	relPath, err := filepath.Rel(base, absDir)
	if err != nil {
		return nil, fmt.Errorf("cannot compute relative path: %w", err)
	}

	if strings.HasPrefix(relPath, "..") {
		return nil, fmt.Errorf("directory %s is not under $HOME/servers", absDir)
	}

	name := NameFromDir(relPath)
	exec, err := findExecutable(absDir)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:       name,
		Dir:        absDir,
		Executable: exec,
		LogFile:    filepath.Join(absDir, "std.log"),
		EnvVars:    CollectEnvVars(),
	}, nil
}

// ResolveService resolves a service from either a name argument or the current directory.
func ResolveService(arg string) (*Service, error) {
	if arg != "" {
		return ServiceFromName(arg)
	}
	// Use current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot get current directory: %w", err)
	}
	return ServiceFromDir(cwd)
}
