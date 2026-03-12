package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSkillsDir(t *testing.T) {
	dir, err := GetSkillsDir()
	if err != nil {
		t.Fatalf("GetSkillsDir() failed: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	expected := filepath.Join(homeDir, ".local", "src", "manage-agent-skills")
	if dir != expected {
		t.Errorf("GetSkillsDir() = %q, want %q", dir, expected)
	}
}

func TestList_EmptyDirectory(t *testing.T) {
	// This test assumes the skills directory doesn't exist or is empty
	skills, err := List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	// Should return empty list if directory doesn't exist
	if skills == nil {
		t.Error("List() returned nil, want empty slice")
	}
}

func TestDelete_NotFound(t *testing.T) {
	err := Delete("nonexistent-skill-xyz", []string{})
	if err == nil {
		t.Error("Delete() with nonexistent skill should return error")
	}
}

func TestDelete_RemovesSkillAndSymlink(t *testing.T) {
	// Use a temporary HOME so GetSkillsDir points to a controlled location
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	skillName := "test-skill"
	skillsDir := filepath.Join(tmpHome, ".local", "src", "manage-agent-skills")
	skillDir := filepath.Join(skillsDir, skillName)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}

	// Set up a temporary agent directory with a symlink
	agentDir := t.TempDir()
	symlinkPath := filepath.Join(agentDir, skillName)
	if err := os.Symlink(skillDir, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	if err := Delete(skillName, []string{agentDir}); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
		t.Error("symlink should not exist after delete")
	}
	if _, err := os.Stat(skillDir); !os.IsNotExist(err) {
		t.Error("skill directory should not exist after delete")
	}
}

func TestDownload_InvalidRepoFormat(t *testing.T) {
	tests := []struct {
		name string
		repo string
	}{
		{"no slash", "invalidrepo"},
		{"too many parts", "owner/repo/extra/parts"},
		{"empty", ""},
		{"missing repo name", "github.com/owner"},
		{"only domain", "github.com/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Download(tt.repo)
			if err == nil {
				t.Error("Download() with invalid repo format should return error")
			}
		})
	}
}
