package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	expected := filepath.Join(homeDir, ".config", "manage-agent-skills", "config.toml")
	if path != expected {
		t.Errorf("GetConfigPath() = %q, want %q", path, expected)
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	// This test assumes config file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() with non-existent file should not fail: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.Agents == nil {
		t.Error("Load() returned config with nil Agents map")
	}

	if len(cfg.Agents) != 0 {
		t.Errorf("Load() returned config with %d agents, want 0", len(cfg.Agents))
	}
}

func TestConfig_Structure(t *testing.T) {
	// Test that we can create a config structure
	cfg := &Config{
		Agents: map[string]string{
			"test-agent": "/path/to/skills",
		},
	}

	if cfg.Agents["test-agent"] != "/path/to/skills" {
		t.Errorf("Config.Agents[test-agent] = %q, want %q",
			cfg.Agents["test-agent"], "/path/to/skills")
	}
}
