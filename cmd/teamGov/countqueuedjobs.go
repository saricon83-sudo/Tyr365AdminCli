package teamGov

import (
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// countqueuedjobsCmd represents the countqueuedjobs command
var countqueuedjobsCmd = &cobra.Command{
	Use:   "countq",
	Short: "counts how many jobs are queued in the Teams Governance API.",
	Long:  `This command counts how many jobs are queued in the Teams Governance API. For example: 365Admin teamGov countqueuedjobs`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		if cmd.Flag("help").Changed {
			cmd.Help()
		}
		body, err := teamGovHttp.Get("CountQueuedJobs")
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CountQueuedJobs",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		resp, err := teamGovHttp.UnmarshalInteger(&body)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CountQueuedJobs",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}
		logger.WithFields(log.Fields{
			"url":    "/api/teams/CountQueuedJobs",
			"method": "GET",
			"status": "Success",
		}).Infof("Number of queued jobs: %d ", resp)
	},
}

func init() {
	TeamGovCmd.AddCommand(countqueuedjobsCmd)
}
