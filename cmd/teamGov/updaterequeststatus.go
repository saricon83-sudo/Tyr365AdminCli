package teamGov

import (
	"fmt"

	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	"github.com/spf13/cobra"
)

// updaterequeststatusCmd represents the updaterequeststatus command
var updaterequeststatusCmd = &cobra.Command{
	Use:   "updaterequeststatus",
	Short: "Updates the status of a request in the Teams Governance API. - Flag : --requestId number, --status string",
	Long: `This request updates the status of a request in the Teams Governance API. The requestId is required to update the status of the request. The status can (read should) be updated to the following values: "Queued", "Processing", "Succeeded", "Error". 
    The status should only be updated to "Queued" if the current status is "Error" or in the very rare case of thread death when "Processing". 
    The status should never be updated to "Processing". This should only be done by the Teams Governance API.
    The status should only be updated to "Error" if the current status is "Processing" or "Queued". 
    The status can only be updated to "Succeeded" if any edits have been made manually beyond the Team Governance API.`,
	Run: func(cmd *cobra.Command, args []string) {
		//If flag --help is used, print the help message
		if cmd.Flag("help").Changed {
			cmd.Help()
		}
		_, err := teamGovHttp.Post("UpdateRequestStatus", map[string]string{"requestId": fmt.Sprintf("%d", requestId), "status": status})
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	updaterequeststatusCmd.Flags().Int32VarP(&requestId, "requestId", "r", 0, "The id of the request to update the status of")
	updaterequeststatusCmd.Flags().StringVarP(&status, "status", "s", "", "The new status of the request")
	TeamGovCmd.AddCommand(updaterequeststatusCmd)
}
