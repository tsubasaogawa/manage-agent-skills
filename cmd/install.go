package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsubasaogawa/manage-agent-skills/pkg/config"
	"github.com/tsubasaogawa/manage-agent-skills/pkg/skills"
)

var installCmd = &cobra.Command{
	Use:   "install [agent-name]",
	Short: "Install skills to an agent",
	Long:  `Install all downloaded skills to the specified agent's skill directory`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		agentPath, exists := cfg.Agents[agentName]
		if !exists {
			return fmt.Errorf("agent '%s' not found in config", agentName)
		}
		
		if err := skills.Install(agentPath); err != nil {
			return fmt.Errorf("failed to install skills: %w", err)
		}
		
		fmt.Printf("Successfully installed skills to agent '%s'\n", agentName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
