package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"clash-tools-go/internal/clash_tools/config"
)

const (
	dockerServiceDir      = "/etc/systemd/system/docker.service.d"
	dockerProxyConfigName = "http-proxy.conf"
	defaultHTTPProxy      = "http://127.0.0.1:7890"
	defaultHTTPSProxy     = "http://127.0.0.1:7890"
	defaultNoProxy        = "localhost,127.0.0.1,::1"
)

var (
	dockerProxyConfigPath = filepath.Join(dockerServiceDir, dockerProxyConfigName)
	dockerHTTPProxy       string
	dockerHTTPSProxy      string
	dockerNoProxy         string
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker daemon proxy",
	Long:  "Manage Docker daemon HTTP/HTTPS proxy configuration for systemd-managed docker service.",
}

var dockerEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable Docker daemon proxy",
	Long:  "Enable Docker daemon proxy by writing systemd drop-in configuration and restarting Docker.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dockerHTTPProxy == "" {
			dockerHTTPProxy = defaultHTTPProxy
		}
		if dockerHTTPSProxy == "" {
			dockerHTTPSProxy = defaultHTTPSProxy
		}
		if dockerNoProxy == "" {
			dockerNoProxy = defaultNoProxy
		}

		if err := os.MkdirAll(dockerServiceDir, 0o755); err != nil {
			return fmt.Errorf("failed to create docker service directory: %w", err)
		}

		content := fmt.Sprintf(`[Service]
Environment="HTTP_PROXY=%s"
Environment="HTTPS_PROXY=%s"
Environment="NO_PROXY=%s"
`, dockerHTTPProxy, dockerHTTPSProxy, dockerNoProxy)

		if err := os.WriteFile(dockerProxyConfigPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write docker proxy config: %w", err)
		}

		if err := restartDocker(); err != nil {
			return err
		}

		fmt.Printf("Enabled Docker daemon proxy\n")
		fmt.Printf("  HTTP_PROXY=%s\n", dockerHTTPProxy)
		fmt.Printf("  HTTPS_PROXY=%s\n", dockerHTTPSProxy)
		fmt.Printf("  NO_PROXY=%s\n", dockerNoProxy)
		return nil
	},
}

var dockerDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable Docker daemon proxy",
	Long:  "Disable Docker daemon proxy by removing systemd drop-in configuration and restarting Docker.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.Remove(dockerProxyConfigPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove docker proxy config: %w", err)
		}

		if err := restartDocker(); err != nil {
			return err
		}

		fmt.Printf("Disabled Docker daemon proxy\n")
		return nil
	},
}

var dockerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker daemon proxy status",
	Long:  "Show Docker daemon proxy configuration status by printing the content of the systemd drop-in file if it exists.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(dockerProxyConfigPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Docker daemon proxy: Disabled")
				return nil
			}
			return fmt.Errorf("failed to read docker proxy config: %w", err)
		}

		fmt.Println("Docker daemon proxy: Enabled")
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.PersistentFlags().StringVar(
		&dockerHTTPProxy,
		"http-proxy",
		"",
		fmt.Sprintf("HTTP proxy for Docker daemon (default %s)", defaultHTTPProxy),
	)
	dockerCmd.PersistentFlags().StringVar(
		&dockerHTTPSProxy,
		"https-proxy",
		"",
		fmt.Sprintf("HTTPS proxy for Docker daemon (default %s)", defaultHTTPSProxy),
	)
	dockerCmd.PersistentFlags().StringVar(
		&dockerNoProxy,
		"no-proxy",
		"",
		fmt.Sprintf("NO_PROXY for Docker daemon (default %s)", defaultNoProxy),
	)

	dockerCmd.AddCommand(dockerEnableCmd)
	dockerCmd.AddCommand(dockerDisableCmd)
	dockerCmd.AddCommand(dockerStatusCmd)
}

// restartDocker reloads systemd and restarts the Docker service.
func restartDocker() error {
	fmt.Println("Restarting Docker daemon...")
	if err := config.RunCommand("sudo", "systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd daemon: %w", err)
	}
	if err := config.RunCommand("sudo", "systemctl", "restart", "docker"); err != nil {
		return fmt.Errorf("failed to restart docker service: %w", err)
	}
	fmt.Println("Docker daemon restarted")
	return nil
}
