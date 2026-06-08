package teamToolboxCmd

import (
	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var includeHidden bool

var getErrors = &cobra.Command{
	Use:   "getErrors",
	Short: "Get requests with status errors",
	Long:  "Get requests from the Governance API that currently has the error status",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching error requests...")
		response, err := client.GetErrors(includeHidden)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Error requests loaded!")

		output.PrintResult(response, func() {
			teamToolboxHelper.PrintViewErrorRequestTable(response)
		})
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getErrors)
}
