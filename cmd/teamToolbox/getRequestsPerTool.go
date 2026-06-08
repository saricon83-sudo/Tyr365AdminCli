package teamToolboxCmd

import (
	"github.com/pterm/pterm"
	teamToolBoxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var getRequestPerTool = &cobra.Command{
	Use:   "getRequestPerTool",
	Short: "Gets requests for a specific tool",
	Long:  `This command queries and returns all entries of a specific tool from the DB`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := teamToolBoxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching requests per tool...")
		response, err := client.GetRequestsByTool()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool request counts loaded!")

		output.PrintResult(response, func() {
			teamToolBoxHelper.PrintToolRequestCountTable(response)
		})
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getRequestPerTool)
}
