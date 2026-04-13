package cmd

import (
	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var stopCmd = &cobra.Command{
	Use:   "stop [service]",
	Short: "Stop a service daemon",
	Long: `Stop a running service daemon. If a service name is provided, stop that service.
If run from a service directory, stop the service in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) error {
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}

	svc, err := service.ResolveService(arg)
	if err != nil {
		return err
	}

	p := platform.New()
	return p.Stop(svc)
}
