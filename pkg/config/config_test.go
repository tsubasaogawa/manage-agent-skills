package config

import (
	"net/http"
	"net/http/httptest"
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

func TestLoadDownloadsDefaultWhenMissing(t *testing.T) {
	// Setup a test HTTP server to serve the default config
	configContent := `# Configuration file for manage-agent-skills
[agents]
claude = "~/.claude/skills"
codex = "~/.codex/skills"
gemini = "~/.gemini/skills"
copilot = "~/.copilot/skills"
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(configContent))
	}))
	defer server.Close()

	// Temporarily override the default URL
	originalURL := DefaultConfigURL
	DefaultConfigURL = server.URL
	defer func() {
		DefaultConfigURL = originalURL
	}()

	// Move existing config if it exists
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	backupPath := configPath + ".backup"
	configExists := false
	if _, err := os.Stat(configPath); err == nil {
		configExists = true
		if err := os.Rename(configPath, backupPath); err != nil {
			t.Fatalf("Failed to backup config: %v", err)
		}
	}

	// Cleanup after test
	defer func() {
		if configExists {
			// Restore original config
			if err := os.Rename(backupPath, configPath); err != nil {
				t.Errorf("Failed to restore config: %v", err)
			}
		} else {
			// Clean up test config
			os.Remove(configPath)
			os.RemoveAll(filepath.Dir(configPath))
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.Agents == nil {
		t.Error("Load() returned config with nil Agents map")
	}

	// Should have downloaded the default config with agents
	if len(cfg.Agents) != 4 {
		t.Errorf("Load() should have downloaded default config with 4 agents, got %d", len(cfg.Agents))
	}

	// Verify specific agents exist
	expectedAgents := []string{"claude", "codex", "gemini", "copilot"}
	for _, agent := range expectedAgents {
		if _, exists := cfg.Agents[agent]; !exists {
			t.Errorf("Expected agent %q not found in config", agent)
		}
	}

	// Verify config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should have been created")
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
