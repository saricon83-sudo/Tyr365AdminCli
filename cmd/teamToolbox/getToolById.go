/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamToolboxCmd

import (
	"fmt"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var toolId string

// getToolByIdCmd represents the getToolById command
var getToolByIdCmd = &cobra.Command{
	Use:   "getToolById",
	Short: "Get details of a tool by its ID (interactive fuzzy select supported)",
	Long: `Get details of a specific tool request from the Team Toolbox catalog.
If you do not specify a --toolId flag, the tool will launch an interactive catalog selection wizard.`,
	Run: func(cmd *cobra.Command, args []string) {
		adminClient, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Team Toolbox API: %v\n", err)
			return
		}

		if toolId == "" {
			spinner, _ := pterm.DefaultSpinner.Start("Loading tools catalog...")
			tools, err := adminClient.GetTools()
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to load tools: %v", err))
				return
			}
			spinner.Success("Loaded tools catalog successfully!")

			var options []string
			toolMap := make(map[string]string)
			for _, t := range tools {
				optionStr := fmt.Sprintf("[%d] %s (Topic: %s)", t.Id, t.ToolName, t.TopicName)
				options = append(options, optionStr)
				toolMap[optionStr] = fmt.Sprintf("%d", t.Id)
			}

			if len(options) == 0 {
				pterm.Warning.Println("No tools found in your database catalog.")
				return
			}

			selectedOption, err := pterm.DefaultInteractiveSelect.
				WithOptions(options).
				Show("Select a tool to view details")
			if err != nil {
				pterm.Error.Printf("Interaction cancelled: %v\n", err)
				return
			}

			toolId = toolMap[selectedOption]
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching details for tool ID %s...", toolId))
		response, err := client.GetToolById(toolId)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to fetch tool: %v", err))
			return
		}
		spinner.Success("Fetched details successfully!")
		
		fmt.Println()
		ViewTable(response)
	},
}

func ViewTable(d teamGovHttp.Printer) {
	output.PrintResult(d, d.PrintTable)
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
