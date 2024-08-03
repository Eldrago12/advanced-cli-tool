package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "act",
	Short: "Advanced CLI tool",
	Long:  `A CLI tool for advanced system tasks and GitHub operations.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(githubCmd)
}
