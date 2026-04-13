package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var initialCmd = &cobra.Command{
	Use:   "initial",
	Short: "Scan current directory and generate daemon configuration",
	Long: `Scan the current directory for an executable file and generate
the platform-native daemon configuration (launchd plist on macOS, systemd unit on Linux).

Must be run from a directory under $HOME/servers.`,
	RunE: runInitial,
}

func init() {
	rootCmd.AddCommand(initialCmd)
}

func runInitial(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Verify we're under $HOME/servers
	base, err := service.ServersBase()
	if err != nil {
		return err
	}

	svc, err := service.ServiceFromDir(cwd)
	if err != nil {
		return err
	}

	// Check if already configured
	p := platform.New()
	if p.IsConfigured(svc) {
		fmt.Printf("Service '%s' is already configured. Use 'mini_launch remove %s' first to reconfigure.\n", svc.Name, svc.Name)
		return nil
	}

	fmt.Printf("Service name: %s\n", svc.Name)
	fmt.Printf("Directory:    %s\n", svc.Dir)
	fmt.Printf("Executable:   %s\n", svc.Executable)
	fmt.Printf("Log file:     %s\n", svc.LogFile)
	fmt.Printf("Base dir:     %s\n", base)

	if err := p.Generate(svc); err != nil {
		return err
	}

	fmt.Printf("\nService '%s' initialized successfully.\n", svc.Name)
	fmt.Printf("Use 'mini_launch start %s' to start the service.\n", svc.Name)
	return nil
}
