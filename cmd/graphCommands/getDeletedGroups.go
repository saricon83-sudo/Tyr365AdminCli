/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package graphCommands

import (
	"encoding/json"
	"fmt"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// Make struct for the bmodels.Group
type MatchRequest struct {
	Ids []string `json:"ids"`
}

var getDeletedGroupsCmd = &cobra.Command{
	Use:   "getDeletedGroups",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		groups, err := graphHelper.GetDeletedGroups()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var groupIds []string
		for _, group := range groups {
			groupIds = append(groupIds, *group.GetId())
		}

		requestBody := MatchRequest{
			Ids: groupIds,
		}

		body, err := teamGovHttp.PostWithBody("GetManagedTeams", nil, requestBody)
		if err != nil {
			fmt.Println("Failed to get managed teams:", err)
			return
		}
		var managedTeams []teamGovHttp.ManagedTeam
		err = json.Unmarshal(body, &managedTeams)

		if err != nil {
			fmt.Println("Failed to unmarshal managed teams:", err)
			return
		}

		// we add flag to print which team has which Origin and Retention
		for _, team := range managedTeams {
			if team.Origin == "GovPortal" && team.Retention == "Forever" {
				fmt.Println(team.TeamName + " is from " + team.Origin + " and needs to be discussed")
			}
			if team.Origin == "Tyra" && team.Retention == "Forever" {
				fmt.Println(team.TeamName + " is from " + team.Origin + " and needs to be restored")
			} else {
				fmt.Println(team.TeamName + " is from " + team.Origin + " and does not need to be restored")
			}
		}
		restoreSelectedTeams(managedTeams)

	},
}

func init() {
	GraphCmd.AddCommand(getDeletedGroupsCmd)

}

func restoreSelectedTeams(managedTeams []teamGovHttp.ManagedTeam) {
	var options []string
	teamNameToGroupId := make(map[string]string) // Map to associate team names with their GroupIds

	// Populate the options slice and the map
	for _, team := range managedTeams {
		option := fmt.Sprintf("Option %s", team.TeamName)
		options = append(options, option)
		teamNameToGroupId[option] = team.GroupId // Use the formatted option as key for consistency
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

	// Add your logic here to restore teams based on selected GroupIds
}
