package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
