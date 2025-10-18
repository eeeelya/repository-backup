package command

import (
	"github.com/eeeelya/repository-backup/internal/core"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	sshKeyPath string
)

var saveCmd = &cobra.Command{
	Use:   "save [repository-url] [output-path]",
	Short: "Save a repository as a zip archive",
	Long: `Save a Git repository as a zip archive.

This command clones a Git repository to a temporary directory, 
creates a zip archive of its contents, and then cleans up the temporary files. 
The .git directory is excluded by default to keep the archive size smaller.

Arguments:
  repository-url    URL of the Git repository (HTTPS or SSH)
  output-path       Path where the zip archive will be saved`,
	Example: `
  # Basic usage
  repobackup save https://github.com/user/repo .

  # Use specific SSH key
  repobackup save git@github.com:user/repo.git /home/user --ssh-key ~/.ssh/custom_key
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repo := core.Repo{URL: args[0]}
		outputPath := args[1]

		if err := repo.Save(outputPath, sshKeyPath); err != nil {
			slog.Error("failed to save repository", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	saveCmd.Flags().StringVar(&sshKeyPath, "ssh-key", "~/.ssh/id_rsa", "Path to SSH private key (optional, defaults to ~/.ssh/id_rsa)")

	rootCmd.AddCommand(saveCmd)
}
