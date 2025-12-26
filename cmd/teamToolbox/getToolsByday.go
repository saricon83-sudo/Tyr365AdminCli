package teamToolboxCmd

import (
	"fmt"
	"os"

	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

var days int

var getRequestsByDay = &cobra.Command{
	Use:   "getRequestsByDay",
	Short: "Get daily request counts for the last N days",
	Long: `Retrieves daily request counts for the last N days from the admin API.

This command calls GET /stats/requests-by-day and returns statistics including:
- Date
- Total count of requests
- Completed requests
- Error count

Example:
  365Admin teamToolbox getRequestsByDay --days 7
  365Admin teamToolbox getRequestsByDay --days 30`,
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating admin API client: %v\n", err)
			return
		}

		results, err := adminAPI.GetRequestsByDay(days)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching requests by day: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("No request data found for the specified period.")
			return
		}

		teamToolboxHelper.PrintDailyRequestCountTable(results)
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getRequestsByDay)
	getRequestsByDay.Flags().IntVar(&days, "days", 30, "Number of days to look back (default: 30)")
}
