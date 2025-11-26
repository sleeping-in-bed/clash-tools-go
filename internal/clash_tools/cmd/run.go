package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sleeping-in-bed/clash-tools-go/internal/clash_tools/config"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run clash server",
	Long:  "Run clash server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.RunCommand("sudo", config.ClashBinaryPath, "-d", config.GlobalConfigDir)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
