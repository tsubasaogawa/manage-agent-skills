package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsubasaogawa/manage-agent-skills/pkg/skills"
)

var downloadCmd = &cobra.Command{
	Use:   "download [repository]",
	Short: "Download a skill from a GitHub repository",
	Long:  `Download a skill from a GitHub repository (e.g., tsubasaogawa/semantic-commit-helper)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		if err := skills.Download(repo); err != nil {
			return fmt.Errorf("failed to download skill: %w", err)
		}
		fmt.Printf("Successfully downloaded skill from %s\n", repo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
