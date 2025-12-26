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

var siteUrl string

// searchGroupsCmd represents the searchGroups command
var testEndpointCmd = &cobra.Command{
	Use:   "testEndpoint",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()

		body, err := teamGovHttp.Get("TestEndpoint", map[string]string{"siteUrl": siteUrl})
		logEntry := logger.WithFields(log.Fields{
			"url":    "/api/teams/TestEndpoint",
			"method": "GET",
		})

		if err != nil {
			errorEntry := logEntry.WithField("status", "Error").WithError(err)
			if len(body) > 0 {
				errorEntry = errorEntry.WithField("responseBody", string(body))
			}
			errorEntry.Error("TestEndpoint request failed")
			return
		}

		logEntry.WithFields(log.Fields{
			"status":       "Success",
			"responseBody": string(body),
		}).Info("TestEndpoint request succeeded")
	},
}

func init() {
	testEndpointCmd.Flags().StringVarP(&siteUrl, "siteUrl", "", "", "siteUrl")
	TeamGovCmd.AddCommand(testEndpointCmd)

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
