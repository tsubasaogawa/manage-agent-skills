package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsubasaogawa/manage-agent-skills/pkg/config"
	"github.com/tsubasaogawa/manage-agent-skills/pkg/skills"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [skill-name]",
	Short: "Delete a downloaded skill",
	Long:  `Delete a skill directory and any symlinks pointing to it from all configured agent directories`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Extract agent paths (values) from the agents map
		agentPaths := make([]string, 0, len(cfg.Agents))
		for _, path := range cfg.Agents {
			agentPaths = append(agentPaths, path)
		}

		if err := skills.Delete(skillName, agentPaths); err != nil {
			return fmt.Errorf("failed to delete skill: %w", err)
		}

		fmt.Printf("Successfully deleted skill '%s'\n", skillName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
