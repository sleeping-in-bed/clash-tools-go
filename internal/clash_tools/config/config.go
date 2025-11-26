package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// GlobalConfigDir is the directory where clash-tools keeps global configuration and assets.
	GlobalConfigDir = "/var/lib/clash_tools/clash"

	// ConfigFileName is the main configuration file name.
	ConfigFileName = "config.yaml"

	// TemplateConfigFileName is the template config file name shipped with the binary.
	TemplateConfigFileName = "config.template.yaml"

	// CountryMMDBName is the geoip database file name.
	CountryMMDBName = "Country.mmdb"

	// ClashBinaryName is the clash executable file name.
	ClashBinaryName = "clash"

	// templateEmbeddedPath is the embedded asset path for the template configuration file.
	templateEmbeddedPath = "assets/" + TemplateConfigFileName

	// countryMMDBEmbeddedPath is the embedded asset path for the Country.mmdb file.
	countryMMDBEmbeddedPath = "assets/" + CountryMMDBName

	// clashEmbeddedPath is the embedded asset path for the clash executable.
	clashEmbeddedPath = "assets/" + ClashBinaryName
)

var (
	// ConfigPath is the absolute path to the main configuration file.
	ConfigPath = filepath.Join(GlobalConfigDir, ConfigFileName)

	// TemplateConfigPath is the absolute path where the template configuration file is stored.
	TemplateConfigPath = filepath.Join(GlobalConfigDir, TemplateConfigFileName)

	// CountryMMDBPath is the absolute path to the Country.mmdb file.
	CountryMMDBPath = filepath.Join(GlobalConfigDir, CountryMMDBName)

	// ClashBinaryPath is the absolute path to the clash executable file.
	ClashBinaryPath = filepath.Join(GlobalConfigDir, ClashBinaryName)
)

// embeddedFiles contains all static assets that are compiled into the binary.
//
//go:embed assets/*
var embeddedFiles embed.FS

// Init ensures that the global configuration directory and required files exist.
func Init() error {
	if err := os.MkdirAll(GlobalConfigDir, 0o755); err != nil {
		return fmt.Errorf("failed to create global config dir: %w", err)
	}

	// Ensure template config file exists on disk.
	if err := ensureFileFromEmbedded(TemplateConfigPath, templateEmbeddedPath, 0o644); err != nil {
		return fmt.Errorf("failed to ensure template config file: %w", err)
	}

	// Initialize main config file from template if it does not exist.
	if err := ensureFileFromEmbedded(ConfigPath, templateEmbeddedPath, 0o644); err != nil {
		return fmt.Errorf("failed to ensure config file: %w", err)
	}

	// Ensure Country.mmdb exists.
	if err := ensureFileFromEmbedded(CountryMMDBPath, countryMMDBEmbeddedPath, 0o644); err != nil {
		return fmt.Errorf("failed to ensure Country.mmdb: %w", err)
	}

	// Ensure clash executable exists and is executable.
	if err := ensureFileFromEmbedded(ClashBinaryPath, clashEmbeddedPath, 0o755); err != nil {
		return fmt.Errorf("failed to ensure clash binary: %w", err)
	}

	if err := os.Chmod(ClashBinaryPath, 0o755); err != nil {
		return fmt.Errorf("failed to set clash binary executable: %w", err)
	}

	return nil
}

// ResetConfig overwrites the main configuration file with the template config on disk.
func ResetConfig() error {
	// Make sure the template file exists (Init should have done this already).
	if err := ensureFileFromEmbedded(TemplateConfigPath, templateEmbeddedPath, 0o644); err != nil {
		return fmt.Errorf("failed to ensure template config file before reset: %w", err)
	}

	data, err := os.ReadFile(TemplateConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read template config file: %w", err)
	}

	tmpPath := ConfigPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write tmp config file: %w", err)
	}

	if err := os.Rename(tmpPath, ConfigPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to replace config file with template: %w", err)
	}

	return nil
}

// ensureFileFromEmbedded writes an embedded asset to dstPath if dstPath does not exist.
func ensureFileFromEmbedded(dstPath, assetPath string, perm os.FileMode) error {
	if _, err := os.Stat(dstPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	data, err := embeddedFiles.ReadFile(assetPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded asset %s: %w", assetPath, err)
	}

	tmpPath := dstPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return fmt.Errorf("failed to write tmp file %s: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, dstPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to move tmp file to %s: %w", dstPath, err)
	}

	return nil
}
