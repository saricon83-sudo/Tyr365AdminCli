package teamGov

import (
	"fmt"

	"github.com/pterm/pterm"
	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var callerID string

// getfailedrequestsCmd represents the getfailedrequests command
var getfailedrequestsCmd = &cobra.Command{
	Use:   "getfailedrequests",
	Short: "Gets failed requests in the Teams Governance API by callerID",
	Long: `This command gets failed requests in the Teams Governance API by callerID. For example: 365Admin teamGov getfailedrequests.
		The response is a table with the following columns: ID, Created, GroupID, TeamName, Endpoint, CallerID, Status, ProvisioningStep, Message, InitiatedBy, Modified, RetryCount, QueuePriority.
		You specify the source by using the flag --callerID. For example: 365Admin teamGov getfailedrequests --callerID "Tyra".
	`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		response, err := teamGovHttp.Get("GetFailedRequests", map[string]string{"callerID": callerID})
		if err != nil {
			logger.WithFields(log.Fields{
				"url":      "/api/teams/GetFailedRequests",
				"method":   "GET",
				"status":   "Error",
				"CallerID": callerID,
			}).Error("error retrivering failed requests")
			return
		}
		requests, errs := teamGovHttp.UnmarshalRequests(&response)
		if errs != nil {
			logger.WithFields(log.Fields{
				"url":      "/api/teams/GetFailedRequests",
				"method":   "GET",
				"status":   "Error",
				"CallerID": callerID,
			}).Error(errs)
			return
		}

		if cmd.Flag("excel").Changed {
			var fileName string
			fmt.Println("Name your new excel file:")
			fmt.Scanln(&fileName)
			saveToFile.SaveToExcel(requests, fileName)
		}
		if cmd.Flag("print").Changed {
			ViewTable(&requests)
		}
		if cmd.Flag("json").Changed {
			var fileName string
			fmt.Println("Enter a name for the JSON file (without extension):")
			fmt.Scanln(&fileName)

			err := saveToFile.SaveDataToJSONFile(requests, fileName+".json")
			if err != nil {
				logger.WithFields(log.Fields{

					"method": "SaveDataTOJsonFile",
					"status": "Error",
				}).Error(err)

				return
			}
			logger.WithFields(log.Fields{

				"method": "SaveDataToJsonFile",
				"status": "Success",
			}).Error("Data successfully saved to JSON file:", fileName+".json")
		}
		if cmd.Flag("interactive").Changed {
			requeueSelectedTeams(requests)
		}

	},
}

func init() {
	getfailedrequestsCmd.Flags().StringVarP(&callerID, "callerID", "c", "", "The callerID to filter the failed requests")
	getfailedrequestsCmd.Flags().Bool("print", false, "Print the response as a table")
	getfailedrequestsCmd.Flags().Bool("excel", false, "Save the response to an Excel file")
	getfailedrequestsCmd.Flags().Bool("json", false, "Save the response to a JSON file")
	getfailedrequestsCmd.Flags().Bool("interactive", false, "interactive mode")
	TeamGovCmd.AddCommand(getfailedrequestsCmd)

}

func requeueSelectedTeams(requests []teamGovHttp.Request) {
	var options []string
	teamNameToGroupId := make(map[string]int) // Map to associate team names with their GroupIds

	// Populate the options slice and the map
	for _, request := range requests {
		option := fmt.Sprintf("Option %s", request.TeamName)
		options = append(options, option)
		teamNameToGroupId[option] = request.ID // Use the formatted option as key for consistency
	}

	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(false).
		WithCheckmark(&pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")})

	selectedOptions, _ := printer.Show()

	var selectedGroupIds []int
	for _, selectedOption := range selectedOptions {
		if groupId, exists := teamNameToGroupId[selectedOption]; exists {
			selectedGroupIds = append(selectedGroupIds, groupId)
		}
	}

	pterm.Info.Printfln("Selected GroupIds: %s", pterm.Green(selectedGroupIds))

	if len(selectedGroupIds) > 0 {
		for _, id := range selectedGroupIds {
			_, err := teamGovHttp.Post("RetryRequest", map[string]string{"requestId": fmt.Sprintf("%d", id)})
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Request with ID", id, "requeued successfully")
		}
	}
}
