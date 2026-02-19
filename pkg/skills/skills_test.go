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

	expected := filepath.Join(homeDir, ".local", "bin", "manage-agent-skills")
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

func TestDownload_InvalidRepoFormat(t *testing.T) {
	tests := []struct {
		name string
		repo string
	}{
		{"no slash", "invalidrepo"},
		{"too many slashes", "owner/repo/extra"},
		{"empty", ""},
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
