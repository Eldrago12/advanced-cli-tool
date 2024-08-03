package cmd

import (
	"fmt"

	"github.com/Eldrago12/advanced-cli-tool/github"

	"github.com/spf13/cobra"
)

var githubCmd = &cobra.Command{
	Use:   "gt",
	Short: "GitHub operations",
	Long:  "Perform various GitHub operations like listing repos, committing, pulling, pushing, and handling branches.",
}

var listReposCmd = &cobra.Command{
	Use:   "ls",
	Short: "List GitHub repositories",
	Long:  "List all GitHub repositories for the authenticated user.",
	Run: func(cmd *cobra.Command, args []string) {
		github.ListRepos()
	},
}

var pushCmd = &cobra.Command{
	Use:   "push [repo] [branch]",
	Short: "Push to GitHub with automatic rebase",
	Long:  "Push changes to a GitHub repository, automatically handling divergences and using Git LFS for large files.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := args[0]
		branch := ""
		if len(args) > 1 {
			branch = args[1]
		}
		github.Push(repo, branch)
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull [branch]",
	Short: "Pull from GitHub with rebase",
	Long:  "Pull changes from a GitHub repository with rebase.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		github.Pull(branch)
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes",
	Long:  "Commit changes to a GitHub repository.",
	Run: func(cmd *cobra.Command, args []string) {
		var message string
		fmt.Print("Enter commit message: ")
		fmt.Scanln(&message)
		github.Commit(message)
	},
}

var createBranchCmd = &cobra.Command{
	Use:   "checkout [branch]",
	Short: "Create and checkout a new branch",
	Long:  "Create and checkout a new branch in a GitHub repository.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		github.CreateBranch(branch)
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GitHub",
	Long:  "Logout from GitHub by removing stored credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		github.Logout()
	},
}

func init() {
	githubCmd.AddCommand(listReposCmd)
	githubCmd.AddCommand(pushCmd)
	githubCmd.AddCommand(pullCmd)
	githubCmd.AddCommand(commitCmd)
	githubCmd.AddCommand(createBranchCmd)
	githubCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(githubCmd)
}
