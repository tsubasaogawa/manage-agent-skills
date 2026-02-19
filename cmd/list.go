package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsubasaogawa/manage-agent-skills/pkg/skills"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all downloaded skills",
	Long:  `List all skills that have been downloaded to ~/.local/src/manage-agent-skills`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillList, err := skills.List()
		if err != nil {
			return fmt.Errorf("failed to list skills: %w", err)
		}

		if len(skillList) == 0 {
			fmt.Println("No skills downloaded yet")
			return nil
		}

		fmt.Println("Downloaded skills:")
		for _, skill := range skillList {
			fmt.Printf("  - %s\n", skill)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
