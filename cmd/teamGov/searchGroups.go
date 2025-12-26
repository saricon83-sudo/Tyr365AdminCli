/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var searchString string

// searchGroupsCmd represents the searchGroups command
var searchGroupsCmd = &cobra.Command{
	Use:   "searchGroups",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
				logger := logging.GetLogger()
		if cmd.Flags().Changed("searchText") {
			body, err := teamGovHttp.Get("GetGroups", map[string]string{"searchText": searchString})
			if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetGroups",
				"method": "POST",
				"status": "Error",
			}).Error(err)
				return
			}
			groups, err := teamGovHttp.UnmarshalGroups(&body)
			if err != nil {
				logger.WithFields(log.Fields{
				"url":    "/api/teams/GetGroups",
				"method": "POST",
				"status": "Error",
			}).Error(err)
				return
			}
			
      ViewTable(&groups)
		}
	},
}

func init() {
	searchGroupsCmd.Flags().StringVarP(&searchString, "searchText", "s", "", "searchText")
	TeamGovCmd.AddCommand(searchGroupsCmd)

}


// func RenderGroups(groups []teamGovHttp.UnifiedGroup) {
// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetHeader([]string{"GroupId", "DisplayName", "Alias", "Description", "CreatedDate", "SharePointUrl", "Visibility", "Team", "Yammer", "Label"}) // Customize the table header as needed

// 	// Populate the table with data from the response
// 	for _, req := range groups {
// 		row := []string{
// 			req.GroupId,
// 			req.DisplayName,
// 			req.Alias,
// 			req.Description,
// 			req.CreatedDate,
// 			req.SharePointUrl,
// 			req.Visibility,
// 			req.Team,
// 			fmt.Sprintf("%v", req.Yammer),
// 			fmt.Sprintf("%v", req.Label),
// 		}
// 		table.Append(row)
// 	}

// 	// Render the table
// 	table.Render()
// }
