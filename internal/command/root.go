package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-repository-backup",
	Short: "Awesome tool for backing up git repositories",
	Long:  `git-repository-backup is an awesome command-line tool for backing up git repositories`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to git-repository-backup!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
