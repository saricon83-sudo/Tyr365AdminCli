/*
Copyright © 2024 NAME HERE <gustav.linden@tyrens.se>
*/
// QuerymanagedCmd represents the command for querying managed teams with specific criteria.
package teamGov

import (
	"fmt"

	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	groupId   string
	teamName  string
	origin    string
	retention string
	fields    string
)

var querymanagedCmd = &cobra.Command{
	Use:   "querymanaged",
	Short: "Query managed teams with specific criteria",
	Long: `Query managed teams based on groupId, teamName, status, origin, retention, and fields.
Example usage:
teamGov querymanaged --groupId "12345" --teamName "MyTeam" --status "active" --origin "internal" --retention "permanent" --fields "Id,teamName,status"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Processing flags and constructing query parameters map
			logger := logging.GetLogger()
		queryParams := make(map[string]string)
		if groupId != "" {
			queryParams["groupId"] = groupId
		}
		if teamName != "" {
			queryParams["teamName"] = teamName
		}
		if origin != "" {
			queryParams["origin"] = origin
		}
		if retention != "" {
			queryParams["retention"] = retention
		}
		if fields != "" {
			queryParams["fields"] = fields
		}
		if status != "" {
			queryParams["status"] = status
		}
		// Get the response from the API
		body, err := teamGovHttp.GetQuery("QueryManagedTeams", queryParams)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/QueryManagedTeams",
				"method": "GET",
				"status": "Error",
				"queryParams": queryParams,
			}).Error(err)
			return
		}
		// Unmarshal (deserialize for C# people) the response
		managedTeams, err := teamGovHttp.UnmarshalManagedTeams(&body)

		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/QueryManagedTeams",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}

		if cmd.Flag("excel").Changed {
			var fileName string
			fmt.Println("Name your new excel file:")
			fmt.Scanln(&fileName)
			saveToFile.SaveToExcel(managedTeams, fileName)
		}
		if cmd.Flag("print").Changed {
			ViewTable(&managedTeams)

		}
		if cmd.Flag("json").Changed {
			var fileName string
			fmt.Println("Enter a name for the JSON file (without extension):")
			fmt.Scanln(&fileName)

			err := saveToFile.SaveDataToJSONFile(managedTeams, fileName+".json")
			if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/QueryManagedTeams",
				"method": "SaveDataToJsonFile",
				"status": "Error",
			}).Error(err)
				return
			}
			logger.WithFields(log.Fields{
				"url":    "/api/teams/QueryManagedTeams",
				"method": "GET",
				"status": "Success",
			}).Info("Successfulle got ")

		}
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

	},
}

func init() {
	querymanagedCmd.Flags().StringVarP(&groupId, "groupId", "", "", "Group ID to filter")
	querymanagedCmd.Flags().StringVarP(&teamName, "teamName", "", "", "Team name to filter")
	querymanagedCmd.Flags().StringVarP(&status, "status", "", "", "Status to filter (e.g., active, inactive)")
	querymanagedCmd.Flags().StringVarP(&origin, "origin", "", "", "Origin to filter (e.g., internal, external)")
	querymanagedCmd.Flags().StringVarP(&retention, "retention", "", "", "Retention policy to filter (e.g., permanent, temporary)")
	querymanagedCmd.Flags().StringVarP(&fields, "fields", "", "", "Comma-separated list of fields to include in the output")
	querymanagedCmd.Flags().BoolP("excel", "x", false, "Save the output to an Excel file")
	querymanagedCmd.Flags().BoolP("print", "p", false, "Print the output to the console")
	querymanagedCmd.Flags().BoolP("json", "j", false, "Save the output to a JSON file")
	TeamGovCmd.AddCommand(querymanagedCmd)

}

// func renderManagedTeams(managed *teamGovHttp.ManagedTeamSlice) {
// 	// Reflect the slice to work with its elements

// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetHeader([]string{"GroupId", "TeamName", "Status", "Origin", "Retention"}) // Customize the table header as needed

// 	// Populate the table with data from the response
// 	for _, req := range *managed {
// 		row := []string{
// 			// fmt.Sprintf("%d", req.Id),
// 			req.GroupId,
// 			req.TeamName,
// 			req.Status,
// 			req.Origin,
// 			req.Retention,
// 		}
// 		table.Append(row)
// 	}

// 	// Render the table
// 	table.Render()
// }
