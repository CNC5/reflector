package cmd

import (
	"reflector/log"
	"reflector/logic"
	"reflector/utils"
	"reflector/xray"

	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a management module",
}

var loadDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect existing servers",
	Run: func(cmd *cobra.Command, args []string) {
		log.
			GetDefaultLogger().Info().
			Update("server", utils.DetectExistingServer()).Msg("detect performed")
	},
}

var loadHTTPServerCmd = &cobra.Command{
	Use:   "httpserver",
	Short: "Detect existing servers and load appropriate modules",
	Run: func(cmd *cobra.Command, args []string) {
		logic.HTTPServerAutoSelect()
	},
}

var loadXrayCmd = &cobra.Command{
	Use:   "xray",
	Short: "Load xray",
	Run: func(cmd *cobra.Command, args []string) {
		xray.NewPortableXray("v25.9.11")
	},
}

// var loadReflectorConfigCmd = &cobra.Command{
// 	Use:   "reflectorconfig",
// 	Short: "Load reflector config",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		rcbytes, err := os.ReadFile(*reflectorConfigLocation)
// 		if err != nil {
// 			log.GetDefaultLogger().Error().
// 				Update("err", err).
// 				Update("location", reflectorConfigLocation).
// 				Msg("failed to open the file")
// 			return
// 		}
// 		c, err := logic.LoadConfig(rcbytes)
// 		if err != nil {
// 			log.GetDefaultLogger().Error().
// 				Update("err", err).
// 				Msg("failed to load the file")
// 			return
// 		}
// 		marshc, _ := yaml.Marshal(c)
// 		xc := xray.NewXrayConfig()
// 		logic.NewReflector()
// 		marshxc, _ := json.Marshal(xc)
// 		log.GetDefaultLogger().Debug().
// 			Update("xrayconfig", fmt.Sprintf("%s", marshxc)).
// 			Update("config", fmt.Sprintf("%s", marshc)).Msg("config loaded")
// 	},
// }

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.AddCommand(loadDetectCmd)
	loadCmd.AddCommand(loadHTTPServerCmd)
	loadCmd.AddCommand(loadXrayCmd)
	// loadCmd.AddCommand(loadReflectorConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
