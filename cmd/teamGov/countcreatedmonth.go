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

var Createdmonth int32

var countcreatedmonthCmd = &cobra.Command{
	Use:   "countcreatedmonth",
	Short: "Returns the number of requests created in the current month.",
	Long: `Returns the number of requests created in the current month. For example: 365Admin teamGov countcreatedmonth
	The endpoints counted against are the following: Create, ApplySPTemplate, ApplyTeamTemplate, Group`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		numbm, err := teamGovHttp.Get("CountEntriesMonth")
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CountEntriesMonth",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		errAtParsing := json.Unmarshal(numbm, &Createdmonth)
		if errAtParsing != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CountEntriesMonth",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		logger.WithFields(log.Fields{
			"url":    "/api/teams/CountEntriesMonth",
			"method": "GET",
			"status": "Error",
		}).Info(Createdmonth)
	},
}

func init() {
	TeamGovCmd.AddCommand(countcreatedmonthCmd)
}
