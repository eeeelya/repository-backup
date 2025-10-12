package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		readFile, err := os.ReadFile("VERSION")

		if err != nil {
			fmt.Println("Error reading VERSION file:", err)
			return
		}
		version := string(readFile)

		fullVersionString := fmt.Sprintf("v%s", version)
		fmt.Println(fullVersionString)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
