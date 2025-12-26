/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	"fmt"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getaverageprovisioningtimeCmd represents the getaverageprovisioningtime command
var getaverageprovisioningtimeCmd = &cobra.Command{
	Use:   "getavg",
	Short: "Get the average provisioning time in the Teams Governance API.",
	Long: `This command gets the average provisioning time in the Teams Governance API. For example: 365Admin teamGov getaverageprovisioningtime
	The time is an average of the time it takes from when a request is queued (Created time) to when it is succeeded (Modified time) `,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		if cmd.Flag("help").Changed {
			cmd.Help()
		}
		body, err := teamGovHttp.Get("AverageProvisionTime")
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/AverageProvisiovTime",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		fmt.Println("", string(body))
	},
}

func init() {
	TeamGovCmd.AddCommand(getaverageprovisioningtimeCmd)

}
