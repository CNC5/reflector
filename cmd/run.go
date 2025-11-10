package cmd

import (
	"fmt"
	"reflector/controller"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the reflector",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		r := controller.NewReflector()
		r.RunWithSignalHandling()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
