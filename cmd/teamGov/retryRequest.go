package teamGov

import (
	"fmt"

	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/spf13/cobra"
)

// requeueCmd represents the requeue command
var retryRequestCmd = &cobra.Command{
	Use:   "retryRequest",
	Short: "This command requeues a request in the Teams Governance API. Flag : --requestId number",
	Long:  `This command requeues a request in the Teams Governance API. For example: 365Admin teamGov retryRequest --requestId 147999`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()

		if cmd.Flag("help").Changed {
			cmd.Help()
		}
		body, err := teamGovHttp.Post("RetryRequest", map[string]string{"requestId": fmt.Sprintf("%d", requestId)})
		if err != nil {
		logger.WithFields(log.Fields{
				"url":    "/api/teams/RetryRequest",
				"method": "POST",
				"status": "Error",
			}).Error(err)
		}

		fmt.Println(string(body))
	},
}

func init() {
	retryRequestCmd.Flags().Int32VarP(&requestId, "requestId", "r", 0, "The id of the request to requeue")
	if err := retryRequestCmd.MarkFlagRequired("requestId"); err != nil {
		fmt.Println(err)
	}

	TeamGovCmd.AddCommand(retryRequestCmd)

}
