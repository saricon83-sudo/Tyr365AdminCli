package teamToolboxCmd

import (
	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var getDashboard = &cobra.Command{
	Use:   "dashboard",
	Short: "Get comprehensive dashboard statistics",
	Long: `Retrieves comprehensive dashboard statistics from the admin API.

This command calls GET /dashboard and returns an overview including:
- Total tools (enabled/disabled)
- Total tool instances
- Total managed teams
- Total requests and tool requests
- Pending archive jobs
- Error and stuck request counts
- Requests by status breakdown
- Tool requests by status breakdown

Example:
  365Admin teamToolbox dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching dashboard statistics...")
		stats, err := adminAPI.GetDashboard()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Dashboard loaded!")

		output.PrintResult(stats, stats.PrintTable)
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getDashboard)
}
