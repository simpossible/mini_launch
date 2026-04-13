package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/simpossible/mini_launch/internal/platform"
	"github.com/simpossible/mini_launch/internal/service"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured services",
	Long:  `List all services found under $HOME/servers along with their status.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	p := platform.New()

	services, err := service.DiscoverServices()
	if err != nil {
		return err
	}

	if len(services) == 0 {
		fmt.Println("No services found under $HOME/servers")
		return nil
	}

	fmt.Printf("%-30s %-12s %s\n", "SERVICE", "STATUS", "DIRECTORY")
	fmt.Printf("%s\n", "------------------------------------------------------------------------")

	for _, svc := range services {
		status := "not configured"
		if p.IsConfigured(svc) {
			s, _ := p.Status(svc)
			status = s
		}
		fmt.Printf("%-30s %-12s %s\n", svc.Name, status, svc.Dir)
	}

	return nil
}
