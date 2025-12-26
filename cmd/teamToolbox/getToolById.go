/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamToolboxCmd

import (
	"fmt"

	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/spf13/cobra"
)
var toolId string
// getToolByIdCmd represents the getToolById command
var getToolByIdCmd = &cobra.Command{
	Use:   "getToolById",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
client , err :=	teamToolboxHelper.CreateClient()
if err != nil {
	fmt.Println(err)
}

response, err := client.GetToolById(toolId)
if err != nil {
	fmt.Println(err)
}
ViewTable(response)

	},
}

func ViewTable(d teamGovHttp.Printer) {
	d.PrintTable()
}
func init() {
	getToolByIdCmd.Flags().StringVarP(&toolId, "toolId", "", "", "Id of tool")
	TeamToolboxCmd.AddCommand(getToolByIdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getToolByIdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getToolByIdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
