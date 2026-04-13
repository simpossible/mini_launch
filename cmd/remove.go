package cmd

import (
	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var removeCmd = &cobra.Command{
	Use:   "remove [service]",
	Short: "Remove a service configuration",
	Long: `Remove a service's daemon configuration and stop it.
If run from a service directory, remove the service in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}

	svc, err := service.ResolveService(arg)
	if err != nil {
		return err
	}

	p := platform.New()
	return p.Remove(svc)
}
