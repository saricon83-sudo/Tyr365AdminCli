package teamToolboxCmd

import (
	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

// getRulesAndLogicCmd represents the getRulesAndLogic command
var getRulesAndLogicCmd = &cobra.Command{
	Use:   "getRulesAndLogic",
	Short: "Displays current Rules and logic from Verktygslådan",
	Long:  `Displays current Rules and logic from Verktygslådan by querying the client API.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching rules and logic...")
		response, err := client.GetRulesAndLogic()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Rules and logic loaded!")

		output.PrintResult(response, func() {
			ViewTable(response)
		})
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getRulesAndLogicCmd)
}
