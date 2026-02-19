package skills

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetSkillsDir returns the directory where skills are stored
func GetSkillsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local", "src", "manage-agent-skills"), nil
}

// Download clones a GitHub repository to the skills directory
func Download(repo string) error {
	skillsDir, err := GetSkillsDir()
	if err != nil {
		return err
	}

	// Create skills directory if it doesn't exist
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	// Store original input for error messages
	originalInput := repo

	// Extract repository name from repo string (e.g., "github.com/tsubasaogawa/semantic-commit-helper" -> "semantic-commit-helper")
	// Remove "github.com/" prefix if present
	repo = strings.TrimPrefix(repo, "github.com/")

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: expected 'github.com/owner/repo', got '%s'", originalInput)
	}
	repoName := parts[1]

	// Check if repository already exists
	targetDir := filepath.Join(skillsDir, repoName)
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("skill '%s' already exists", repoName)
	}

	// Clone the repository
	gitURL := fmt.Sprintf("https://github.com/%s.git", repo)
	cmd := exec.Command("git", "clone", gitURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// List returns a list of all downloaded skills
func List() ([]string, error) {
	skillsDir, err := GetSkillsDir()
	if err != nil {
		return nil, err
	}

	// Check if skills directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory contents
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	skills := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			skills = append(skills, entry.Name())
		}
	}

	return skills, nil
}

// Delete removes a skill directory and any symlinks pointing to it from the given agent directories
func Delete(skillName string, agentPaths []string) error {
	skillsDir, err := GetSkillsDir()
	if err != nil {
		return err
	}

	skillDir := filepath.Join(skillsDir, skillName)

	// Check if the skill exists
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", skillName)
	}

	// Remove symlinks from all agent directories
	for _, agentPath := range agentPaths {
		if strings.HasPrefix(agentPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			agentPath = filepath.Join(homeDir, agentPath[2:])
		}

		symlinkPath := filepath.Join(agentPath, skillName)
		info, err := os.Lstat(symlinkPath)
		if err != nil {
			// Symlink doesn't exist, skip
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(symlinkPath); err != nil {
				fmt.Printf("  Warning: failed to remove symlink %s: %v\n", symlinkPath, err)
			} else {
				fmt.Printf("  Removed symlink %s\n", symlinkPath)
			}
		}
	}

	// Remove the skill directory
	if err := os.RemoveAll(skillDir); err != nil {
		return fmt.Errorf("failed to remove skill directory: %w", err)
	}

	return nil
}

// Install creates symbolic links for all skills in the agent's skill directory
func Install(agentSkillDir string) error {
	skillsDir, err := GetSkillsDir()
	if err != nil {
		return err
	}

	// Check if skills directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return fmt.Errorf("no skills downloaded yet")
	}

	// Expand tilde in agent skill directory path
	if strings.HasPrefix(agentSkillDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		agentSkillDir = filepath.Join(homeDir, agentSkillDir[2:])
	}

	// Create agent skill directory if it doesn't exist
	if err := os.MkdirAll(agentSkillDir, 0755); err != nil {
		return fmt.Errorf("failed to create agent skill directory: %w", err)
	}

	// Read all skills
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return fmt.Errorf("failed to read skills directory: %w", err)
	}

	installedCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()
		sourcePath := filepath.Join(skillsDir, skillName)
		targetPath := filepath.Join(agentSkillDir, skillName)

		// Check if symlink already exists
		if _, err := os.Lstat(targetPath); err == nil {
			fmt.Printf("  Skipping %s (already exists)\n", skillName)
			continue
		}

		// Create symbolic link
		if err := os.Symlink(sourcePath, targetPath); err != nil {
			fmt.Printf("  Warning: failed to create symlink for %s: %v\n", skillName, err)
			continue
		}

		fmt.Printf("  Installed %s\n", skillName)
		installedCount++
	}

	if installedCount == 0 {
		fmt.Println("  No new skills to install")
	}

	return nil
}
