package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var statusCmd = &cobra.Command{
	Use:   "status [service]",
	Short: "Show service status",
	Long: `Show the status of a service. If a service name is provided, show that service's status.
If no argument is given, show the status of all configured services.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	p := platform.New()

	if len(args) > 0 {
		svc, err := service.ResolveService(args[0])
		if err != nil {
			return err
		}

		status, err := p.Status(svc)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", svc.Name, status)
		return nil
	}

	// Show all services
	services, err := service.DiscoverServices()
	if err != nil {
		return err
	}

	if len(services) == 0 {
		fmt.Println("No services found under $HOME/servers")
		return nil
	}

	for _, svc := range services {
		configured := "not configured"
		if p.IsConfigured(svc) {
			status, _ := p.Status(svc)
			configured = status
		}
		fmt.Printf("%-30s %s\n", svc.Name, configured)
	}

	return nil
}
