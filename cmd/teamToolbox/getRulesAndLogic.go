/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamToolboxCmd

import (
	"fmt"

	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

// getRulesAndLogicCmd represents the getRulesAndLogic command
var getRulesAndLogicCmd = &cobra.Command{
	Use:   "getRulesAndLogic",
	Short: "Displays current Rules and logic from Verktygslådan",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		response, err := client.GetRulesAndLogic()
		if err != nil {
			fmt.Println(err)
			return
		}
		ViewTable(response)
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getRulesAndLogicCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getRulesAndLogicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getRulesAndLogicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
