/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	"encoding/json"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var Createdweek int32

// countcreatedweekCmd represents the countcreatedweek command
var countcreatedweekCmd = &cobra.Command{
	Use:   "countcreatedweek",
	Short: "Count the number of created requests in the Teams Governance API for the current week.",
	Long: `This command counts the number of created requests in the Teams Governance API for the current week. For example: 365Admin teamGov countcreatedweek
	The endpoints counted against are the following: Create, ApplySPTemplate, ApplyTeamTemplate, Group`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		numbm, err := teamGovHttp.Get("CountEntriesWeek")
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CountEntriesWeek",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		errAtParsing := json.Unmarshal(numbm, &Createdweek)
		if errAtParsing != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CountEntriesWeek",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		logger.WithFields(log.Fields{
			"url":    "/api/teams/CountEntriesWeek",
			"method": "GET",
			"status": "Success",
		}).Info(Createdweek)
	},
}

func init() {
	TeamGovCmd.AddCommand(countcreatedweekCmd)

}
