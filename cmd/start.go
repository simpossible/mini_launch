package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var startCmd = &cobra.Command{
	Use:   "start [service]",
	Short: "Start a service daemon",
	Long: `Start a service daemon. If a service name is provided, start that service.
If run from a service directory, start the service in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}

	svc, err := service.ResolveService(arg)
	if err != nil {
		return err
	}

	p := platform.New()
	if !p.IsConfigured(svc) {
		return fmt.Errorf("service '%s' is not configured. Run 'mini_launch initial' first", svc.Name)
	}

	return p.Start(svc)
}
