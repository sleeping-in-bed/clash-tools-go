package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sleeping-in-bed/clash-tools-go/internal/clash_tools/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage clash configuration",
	Long:  "Manage clash configuration files.",
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print clash configuration path",
	Long:  "Print the absolute path to the clash configuration file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(config.ConfigPath)
		return nil
	},
}

var configCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "Print clash configuration content",
	Long:  "Print the content of the clash configuration file to stdout.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(config.ConfigPath)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit clash configuration",
	Long:  "Edit clash configuration using the default editor.",
	RunE: func(cmd *cobra.Command, args []string) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = os.Getenv("VISUAL")
		}
		if editor == "" {
			editor = "nano"
		}

		return config.RunCommand(editor, config.ConfigPath)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset clash configuration",
	Long:  "Reset clash configuration to the embedded template.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ResetConfig(); err != nil {
			return err
		}
		fmt.Println("configuration reset from template")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configCatCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configResetCmd)
}
