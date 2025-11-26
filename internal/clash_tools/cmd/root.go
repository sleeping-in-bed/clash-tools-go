package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"clash-tools-go/internal/clash_tools/config"
)

var rootCmd = &cobra.Command{
	Use:   "clash-tools",
	Short: "Clash tools CLI",
	Long:  "Clash tools CLI for managing clash.",
}

func Execute() {
	if os.Geteuid() != 0 {
		fmt.Println("clash-tools must be run as root (use sudo).")
		os.Exit(1)
	}

	cobra.CheckErr(config.Init())
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}
