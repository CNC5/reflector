package cmd

import (
	"fmt"
	"reflector/logic"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(logic.VersionString())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
