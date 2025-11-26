package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/sleeping-in-bed/clash-tools-go/internal/clash_tools/config"
)

type clashConfig struct {
	Port      int `yaml:"port"`
	SocksPort int `yaml:"socks-port"`
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Print shell proxy exports for clash",
	Long:  "Print shell proxy export commands based on clash config so you can use eval \"$(clash-tools proxy)\".",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpPort, socksPort, err := readPortsFromConfig(config.ConfigPath)
		if err != nil {
			return err
		}

		fmt.Printf("export http_proxy='http://127.0.0.1:%d'\n", httpPort)
		fmt.Printf("export https_proxy='http://127.0.0.1:%d'\n", httpPort)
		fmt.Printf("export HTTP_PROXY='http://127.0.0.1:%d'\n", httpPort)
		fmt.Printf("export HTTPS_PROXY='http://127.0.0.1:%d'\n", httpPort)
		fmt.Printf("export all_proxy='socks5://127.0.0.1:%d'\n", socksPort)
		fmt.Printf("export ALL_PROXY='socks5://127.0.0.1:%d'\n", socksPort)
		fmt.Println("export no_proxy='localhost,127.0.0.1,::1'")
		fmt.Println("export NO_PROXY='localhost,127.0.0.1,::1'")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}

// readPortsFromConfig extracts port and socks-port from a YAML config file.
func readPortsFromConfig(path string) (int, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, fmt.Errorf("config file not found: %s", path)
		}
		return 0, 0, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg clashConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return 0, 0, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	if cfg.Port == 0 {
		return 0, 0, fmt.Errorf("port not found in config: %s", path)
	}
	if cfg.SocksPort == 0 {
		return 0, 0, fmt.Errorf("socks-port not found in config: %s", path)
	}

	return cfg.Port, cfg.SocksPort, nil
}
