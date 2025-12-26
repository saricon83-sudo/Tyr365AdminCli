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

// countprocessingjobsCmd represents the countprocessingjobs command
var countprocessingjobsCmd = &cobra.Command{
	Use:   "countp",
	Short: "The commands counts the number of processing jobs in the Teams Governance API.",
	Long:  `This command counts the number of processing jobs in the Teams Governance API. For example: 365Admin teamGov countprocessingjobs`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
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
		requests, err := teamGovHttp.UnmarshalRequests(&body)
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
			"status": "Succeess",
		}).Info(requests)
	},
}

func init() {
	TeamGovCmd.AddCommand(countprocessingjobsCmd)
}
