package teamGov

import (
	"fmt"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getprovisioningstatusCmd = &cobra.Command{
	Use:   "getrequest",
	Short: "Get provisioning status of a request. Flag : --requestId number",
	Long:  `Get the provisioning status of a request. For example: 365Admin teamGov getprovisioningstatus --requestId 147999`,

	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		if cmd.Flag("help").Changed {
			cmd.Help()
		}

		body, err := teamGovHttp.Get("GetProvisioningStatus", map[string]string{"requestId": fmt.Sprintf("%d", requestId)})
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetProvisioningStatus",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		requests, err := teamGovHttp.UnmarshalRequests(&body)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/GetProvisioningStatus",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		logger.WithFields(log.Fields{
			"url":    "/api/teams/GetProvisioningStatus",
			"method": "GET",
			"status": "Success",
		}).Info(err)
		ViewTable(&requests)
	},
}

func init() {
	getprovisioningstatusCmd.Flags().Int32VarP(&requestId, "requestId", "r", 0, "The id of the request to requeue")
	if err := getprovisioningstatusCmd.MarkFlagRequired("requestId"); err != nil {
		fmt.Print("Supply a RequestId")
	}

	TeamGovCmd.AddCommand(getprovisioningstatusCmd)
}

func ViewTable(d teamGovHttp.Printer) {
	output.PrintResult(d, d.PrintTable)
}
