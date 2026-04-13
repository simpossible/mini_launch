package cmd

import (
	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var restartCmd = &cobra.Command{
	Use:   "restart [service]",
	Short: "Restart a service daemon",
	Long: `Restart a service daemon. If a service name is provided, restart that service.
If run from a service directory, restart the service in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRestart,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func runRestart(cmd *cobra.Command, args []string) error {
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}

	svc, err := service.ResolveService(arg)
	if err != nil {
		return err
	}

	p := platform.New()
	return p.Restart(svc)
}
