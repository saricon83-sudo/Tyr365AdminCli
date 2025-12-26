/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	archivercmd "github.com/saricon83-sudo/Tyr365AdminCli/cmd/archiver"
	"github.com/saricon83-sudo/Tyr365AdminCli/cmd/azure"
	"github.com/saricon83-sudo/Tyr365AdminCli/cmd/graphCommands"
	"github.com/saricon83-sudo/Tyr365AdminCli/cmd/sp"
	"github.com/saricon83-sudo/Tyr365AdminCli/cmd/teamGov"
	teamToolboxCmd "github.com/saricon83-sudo/Tyr365AdminCli/cmd/teamToolbox"

	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logFile string
)
var cfgFile string
var Output bool
var logFilePath string
var debug bool
var useJSON bool
var (
	fileLog bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "365Admin",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if fileLog {
			today := time.Now().Format("06-01-02.json") // yy-mm-dd.json format
			logging.SetupLogging(today, useJSON)        // Setup logging to JSON file named with today's date
		} else {
			logging.SetupLogging("", useJSON) // Setup default logging to stdout in text format
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("365Admin")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.GetLogger().Fatal(err)
	}
	logging.CloseLogging()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(sp.SpCmd)
	rootCmd.AddCommand(teamGov.TeamGovCmd)
	rootCmd.AddCommand(graphCommands.GraphCmd)
	rootCmd.AddCommand(azure.AzureCmd)
	rootCmd.AddCommand(teamToolboxCmd.TeamToolboxCmd)
	rootCmd.AddCommand(archivercmd.ArchiverCmd)
	// rootCmd.PersistentFlags().StringVarP(&logFilePath, "log", "l", "", "Specify the log file path")
	// rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Specify config file (default is $HOME/.config.json)")
	rootCmd.PersistentFlags().BoolVarP(&Output, "stdout", "", false, "Output to standard output")
	rootCmd.PersistentFlags().BoolVarP(&useJSON, "useJson", "", false, "Output logs in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&fileLog, "fileLog", "", false, "Log to a JSON file named with today's date")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".365Admin" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".365Admin")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
