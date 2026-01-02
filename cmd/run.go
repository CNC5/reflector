package cmd

import (
	"os"
	"reflector/log"
	"reflector/logic"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the reflector",
	Run: func(cmd *cobra.Command, args []string) {
		r := logic.NewReflector(
			"v2.9.0",
			"v25.9.11",
		)
		config, err := os.Open(*reflectorConfigLocation)
		if err != nil {
			log.GetDefaultLogger().
				Error().
				Update("err", err.Error()).
				Msgf(
					"failed to open %s",
					*reflectorConfigLocation)
			os.Exit(1)
		}
		err = r.ParseReflectorConfigV1(config)
		if err != nil {
			log.GetDefaultLogger().
				Error().
				Update("err", err.Error()).
				Msg("failed to parse config")
			os.Exit(1)
		}
		r.RunWithSignalHandling()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
