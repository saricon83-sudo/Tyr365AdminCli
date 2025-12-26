package teamGov

import (
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/spf13/cobra"
)

// getprocessingrequestsCmd represents the getprocessingrequests command
var getprocessingrequestsCmd = &cobra.Command{
	Use:   "countprequests",
	Short: "Get the number of processing requests in the Teams Governance API.",
	Long:  `This command gets the number of processing requests in the Teams Governance API. For example: 365Admin teamGov getprocessingrequests`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		//If flag --help is used, print the help message
		if cmd.Flag("help").Changed {
			cmd.Help()
		}
		body, err := teamGovHttp.Get("GetProcessingRequests")
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetProcessingRequests",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		requests, err := teamGovHttp.UnmarshalInteger(&body)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetProcessingRequests",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		logger.WithFields(log.Fields{
			"url":    "/api/teams/GetProcessingRequests",
			"method": "GET",
			"status": "Success",
		}).Info(requests)
	},
}

func init() {
	TeamGovCmd.AddCommand(getprocessingrequestsCmd)

}
