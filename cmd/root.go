package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mini_launch",
	Short: "A simple service daemon manager for macOS and Linux",
	Long: `mini_launch manages service daemons using native platform tools:
  - macOS: launchd (launchctl)
  - Linux: systemd --user

Services are organized under $HOME/servers with one executable per directory.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
