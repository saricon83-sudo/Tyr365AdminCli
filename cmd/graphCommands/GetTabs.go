/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package graphCommands

import (
	"fmt"

	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/pterm/pterm"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/spf13/cobra"
)

var teamGuid string
var searchString string

// ensureFilesFolderCmd represents the ensureFilesFolder command
var getTabsCmd = &cobra.Command{
	Use:   "getTabs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		body, err := teamGovHttp.Get("GetGroups", map[string]string{"searchText": searchString})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		groups, err := teamGovHttp.UnmarshalGroups(&body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		teams, err := selectTeams(groups)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, team := range teams {
			allChannels, err := graphHelper.GetAllChannels(team)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			selectedChannels, err := selectChannels(allChannels)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			for _, channel := range selectedChannels {
				tabs, err := graphHelper.GetTabs(team, channel)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				for _, tab := range tabs.GetValue() {
					fmt.Println("Tab:", *tab.GetDisplayName())
					fmt.Println("TabId:", *tab.GetId())
					fmt.Println("TabWebUrl:", *tab.GetWebUrl())

				}

				// fmt.Println("Drive:", *drive.GetName())
				// fmt.Println("DriveId:", *drive.GetId())
				// fmt.Println("DriveWebUrl:", *drive.GetWebUrl())
				// fmt.Println("DriveRetentionLabel:", drive.GetRetentionLabel())
				// fmt.Println("DriveLastModifiedBy:", drive.GetLastModifiedBy())
				// fmt.Println("DriveCreatedBy:", drive.GetCreatedBy())
				// fmt.Println("DriveGetSpecialFolders:", drive.GetSpecialFolder())

			}
		}
	},
}

func init() {
	//add flag to search for groups
	getTabsCmd.Flags().StringVarP(&searchString, "searchText", "s", "", "searchText")
	getTabsCmd.Flags().StringVarP(&teamGuid, "teamGuid", "t", "", "The guid of the team to ensure the files folder for")
	GraphCmd.AddCommand(getTabsCmd)

}
func selectTeams(groups []teamGovHttp.UnifiedGroup) ([]string, error) {

	var options []string
	teamNameToGroupId := make(map[string]string) // Map to associate team names with their GroupIds

	// Populate the options slice and the map
	for _, group := range groups {
		option := fmt.Sprintf("Option %s", group.DisplayName)
		options = append(options, option)
		teamNameToGroupId[option] = group.GroupId // Use the formatted option as key for consistency
	}

	// Create a new interactive multiselect printer with the options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(false).
		WithCheckmark(&pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")})

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	// Initialize a slice for selected GroupIds based on selected options
	var selectedGroupIds []string
	for _, selectedOption := range selectedOptions {
		if groupId, exists := teamNameToGroupId[selectedOption]; exists {
			selectedGroupIds = append(selectedGroupIds, groupId)
		}
	}

	pterm.Info.Printfln("Selected GroupIds: %s", pterm.Green(selectedGroupIds))

	if len(selectedGroupIds) == 0 {
		return nil, fmt.Errorf("No groups selected")
	}
	return selectedGroupIds, nil
}
func selectChannels(teamChannels models.ChannelCollectionResponseable) ([]string, error) {
	var options []string
	teamNameToGroupId := make(map[string]string) // Map to associate team names with their GroupIds

	// Populate the options slice and the map
	for _, channel := range teamChannels.GetValue() {
		option := fmt.Sprintf("Option %s", *channel.GetDisplayName())
		options = append(options, option)
		teamNameToGroupId[option] = *channel.GetId() // Use the formatted option as key for consistency
	}

	// Create a new interactive multiselect printer with the options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(false).
		WithCheckmark(&pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")})

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	// Initialize a slice for selected GroupIds based on selected options
	var selectedGroupIds []string
	for _, selectedOption := range selectedOptions {
		if groupId, exists := teamNameToGroupId[selectedOption]; exists {
			selectedGroupIds = append(selectedGroupIds, groupId)
		}
	}

	pterm.Info.Printfln("Selected GroupIds: %s", pterm.Green(selectedGroupIds))

	if len(selectedGroupIds) == 0 {
		return nil, fmt.Errorf("No groups selected")
	}
	return selectedGroupIds, nil
}
