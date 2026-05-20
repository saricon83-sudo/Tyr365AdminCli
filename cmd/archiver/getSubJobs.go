/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package archivercmd

import (
	"fmt"

	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

// getSubJobsCmd represents the getSubJobs command
var GetSubJobsCmd = &cobra.Command{
	Use:   "getSubJobs",
	Short: "Get archive sub jobs by status",
	Long: `Get archive sub jobs by status and display them in a formatted table.
Example usage:
  getSubJobs --status completed
  getSubJobs --status processing`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getSubJobs called")

		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			status = "processing" // default status
		}

		client, err := archiver.NewArchiverClient()
		if err != nil {
			fmt.Println("Error creating Archiver client:", err)
			return
		}

		fmt.Printf("=== Getting archive sub jobs with status: %s ===\n", status)
		subJobs, err := client.GetArchiveSubJobsByStatusTyped(status)
		if err != nil {
			fmt.Println("Error fetching sub jobs:", err)
			return
		}

		fmt.Printf("Found %d sub jobs with status '%s'\n", len(subJobs), status)

		if len(subJobs) == 0 {
			fmt.Println("No sub jobs found with that status")
			return
		}

		// Display in table format using the TablePrinter functionality
		if len(subJobs) <= 20 {
			subJobs.PrintTable()
		} else {
			fmt.Printf("Too many results (%d) to display in table. Showing first 20:\n", len(subJobs))
			firstTwenty := subJobs[:20]
			firstTwenty.PrintTable()
		}
	},
}

func init() {
	ArchiverCmd.AddCommand(GetSubJobsCmd)

	// Add flags
	GetSubJobsCmd.Flags().StringP("status", "s", "processing", "Status to filter by (e.g., processing, completed, failed)")
}
