package teamToolboxCmd

import (
	"fmt"
	"os"

	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
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
			fmt.Fprintf(os.Stderr, "Error creating admin API client: %v\n", err)
			return
		}

		stats, err := adminAPI.GetDashboard()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching dashboard: %v\n", err)
			return
		}

		stats.PrintTable()
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getDashboard)
}
