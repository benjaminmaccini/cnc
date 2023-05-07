package cmd

import (
	"cnc/pkg/utils"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var cfgFile string
var debug bool
var internal bool

var rootCmd = &cobra.Command{
	Use: "cnc",
	Long: `A productivity utility that can also be exposed as a web server for certain routes

Ex.
cnc chat --query='What is a vector database?' --config='$HOME/.config/cnc/config.yaml'

or set the values with a config file (default $HOME/.config/cnc/config.yaml)
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/cnc/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Set log level to debug")
	rootCmd.PersistentFlags().BoolVarP(&internal, "internal", "", false, "Set the command to internal mode. Restricts certain usages")
}

// Reads in the config and sets it to a globally available struct
func initConfig() {
	utils.L = &utils.BaseLogger{L: log.New()}

	// Set the log level
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search $HOME/.config/cnc/config.yaml
		viper.AddConfigPath(home + "/.config/cnc")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"fileName": viper.ConfigFileUsed()}).Debug("Config successfully loaded")
	}

	// Convert the config to the global struct
	utils.C = &utils.Config{}
	err := viper.Unmarshal(utils.C)
	if err != nil {
		log.Fatal(err)
	}
}
