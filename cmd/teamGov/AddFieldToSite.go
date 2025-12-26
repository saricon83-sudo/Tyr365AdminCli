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

var alias string

// AddFieldToSiteCmd represents the AddFieldToSite command
var AddFieldToSiteCmd = &cobra.Command{
	Use:   "AddFieldToSite",
	Short: "Calls the AddFieldToSite endpoint with alias parameter",
	Long:  `The cmd calls an endpoint in the Teams Governance API`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		if cmd.Flag("alias").Changed {
			queryParams := make(map[string]string)
			queryParams["alias"] = alias
			_, err := teamGovHttp.Get("AddFieldToSite", queryParams)

			if err != nil {
				logger.WithFields(log.Fields{
					"url":             "/api/teams/AddFieldToSite",
					"method":          "Get",
					"error":           err,
					"queryparameters": queryParams["alias"],
				}).Error("Failed to execute request")

				return
			}

			logger.WithFields(log.Fields{
				"url":             "/api/teams/AddFieldToSite",
				"method":          "Get",
				"status":          "success",
				"queryparameters": queryParams["alias"],
			}).Infof("Successfully added field to site %s\n", alias)

		}
	},
}

func init() {
	AddFieldToSiteCmd.Flags().StringVarP(&alias, "alias", "a", "", "alias of sp site")
	TeamGovCmd.AddCommand(AddFieldToSiteCmd)

}
