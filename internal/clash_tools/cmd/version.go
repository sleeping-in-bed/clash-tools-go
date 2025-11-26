package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// getBuildVersion returns the version string from Go build info.
func getBuildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return "unknown"
	}
	if info.Main.Version == "" {
		return "(devel)"
	}
	return info.Main.Version
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print clash-tools version",
	Long:  "Print the current clash-tools binary version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(getBuildVersion())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
