package cmd

import (
	"os"
	"reflector/log"

	"github.com/spf13/cobra"
)

var debug *bool
var reflectorConfigLocation *string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reflector",
	Short: "",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if *debug {
			log.SetDefaultLogger(log.NewLogger("reflector", log.DEBUG))
		} else {
			log.SetDefaultLogger(log.NewLogger("reflector", log.INFO))
		}
		log.GetDefaultLogger().Debug().Msg("Debugging enabled")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	debug = rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debugging")
	reflectorConfigLocation =
		rootCmd.PersistentFlags().StringP(
			"config", "c",
			"./config.yaml",
			"Reflector config location",
		)
}
