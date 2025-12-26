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

var date string

// GetRequestsByDayCmd represents the GetRequestsByDay command
var GetRequestsByDayCmd = &cobra.Command{
	Use:   "GetRequestsByDay",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		queryParams := make(map[string]string)
		if date != "" {
			queryParams["date"] = date
		}

		body, err := teamGovHttp.Get("GetRequestsForDay", queryParams)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetRequestsForDay",
				"method": "GET",
				"status": "Error",
				"query":  queryParams,
			}).Error(err)
			return
		}
		requests, err := teamGovHttp.UnmarshalRequests(&body)

		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetRequestsForDay",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		ViewTable(&requests)

	},
}

func init() {
	GetRequestsByDayCmd.Flags().StringVarP(&date, "date", "d", "", "Date in fashion YYYY/MM/DD")
	TeamGovCmd.AddCommand(GetRequestsByDayCmd)
}
