package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manage-agent-skills",
	Short: "A tool to manage agent skills",
	Long:  `A tool to download, install, and manage agent skills from GitHub repositories.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
