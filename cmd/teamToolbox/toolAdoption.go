package teamToolboxCmd

import (
	"github.com/pterm/pterm"
	teamToolBoxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var toolAdoption = &cobra.Command{
	Use:   "toolAdoption",
	Short: "Get tool adoption rates",
	Long:  "Get the adoption rates of all tools available in the Team Toolbox",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := teamToolBoxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching tool adoption statistics...")
		response, err := client.GetToolAdoption()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool adoption stats loaded!")

		output.PrintResult(response, func() {
			teamToolBoxHelper.PrintToolAdoptionStatsTable(response)
		})
	},
}

func init() {
	TeamToolboxCmd.AddCommand(toolAdoption)
}
