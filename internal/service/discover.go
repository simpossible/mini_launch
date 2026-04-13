package service

import (
	"fmt"
	"os"
	"path/filepath"
)

// findExecutable finds exactly one executable file in the given directory.
// It excludes directories, hidden files, and std.log.
func findExecutable(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("cannot read directory %s: %w", dir, err)
	}

	var executables []string
	for _, entry := range entries {
		name := entry.Name()

		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Skip hidden files
		if len(name) > 0 && name[0] == '.' {
			continue
		}

		// Skip log file
		if name == "std.log" {
			continue
		}

	 fullPath := filepath.Join(dir, name)

		// Check if executable
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		// Check execute permission
		if info.Mode()&0111 != 0 {
			executables = append(executables, fullPath)
		}
	}

	if len(executables) == 0 {
		return "", fmt.Errorf("no executable file found in %s", dir)
	}

	if len(executables) > 1 {
		names := make([]string, len(executables))
		for i, p := range executables {
			names[i] = filepath.Base(p)
		}
		return "", fmt.Errorf("multiple executable files found in %s: %v\nPlease keep only one executable per service directory", dir, names)
	}

	return executables[0], nil
}

// DiscoverServices scans $HOME/servers for all directories containing an executable.
func DiscoverServices() ([]*Service, error) {
	base, err := ServersBase()
	if err != nil {
		return nil, err
	}

	var services []*Service
	err = filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}

		if !d.IsDir() {
			return nil
		}

		// Skip hidden directories
		if len(d.Name()) > 0 && d.Name()[0] == '.' {
			return filepath.SkipDir
		}

		// Try to find an executable in this directory
		exec, findErr := findExecutable(path)
		if findErr != nil {
			return nil // not a service directory, continue walking
		}

		relPath, relErr := filepath.Rel(base, path)
		if relErr != nil {
			return nil
		}

		name := NameFromDir(relPath)
		services = append(services, &Service{
			Name:       name,
			Dir:        path,
			Executable: exec,
			LogFile:    filepath.Join(path, "std.log"),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}
