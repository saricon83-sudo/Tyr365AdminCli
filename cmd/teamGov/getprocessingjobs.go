/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	"encoding/json"
	"fmt"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getprocessingjobsCmd represents the getprocessingjobs command
var getprocessingjobsCmd = &cobra.Command{
	Use:   "getprocessing",
	Short: "Get the processing jobs in the Teams Governance API.",
	Long: `This command gets the processing jobs in the Teams Governance API. For example: 365Admin teamGov getprocessingjobs. 
    The response is a table with the following columns: ID, Created, GroupID, TeamName, Endpoint, CallerID, Status, ProvisioningStep, Message, InitiatedBy, Modified, RetryCount, QueuePriority.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		if cmd.Flag("help").Changed {
			cmd.Help()
		}
		body, err := teamGovHttp.Get("GetProcessingJobs")
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetProcessingJobs",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		requests, err := teamGovHttp.UnmarshalRequests(&body)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetProcessingJobs",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		//    RenderData(requests)
		if cmd.Flag("showData").Changed {
			outData, _ := json.Marshal(requests)
			fmt.Println(string(outData))
		} else {
			ViewTable(&requests)

		}
	},
}

func init() {

	getprocessingjobsCmd.Flags().Bool("showData", false, "Prints the table")
	TeamGovCmd.AddCommand(readFileCmd)

	TeamGovCmd.AddCommand(getprocessingjobsCmd)

}
