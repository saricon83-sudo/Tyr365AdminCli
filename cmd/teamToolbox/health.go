package teamToolboxCmd

import (
	"strconv"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command group
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Inspect database health and run cleanups",
}

// dbHealthCmd represents the db command
var dbHealthCmd = &cobra.Command{
	Use:   "db",
	Short: "Check the health and connection status of the backing database",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Checking database health...")
		res, err := adminAPI.GetDatabaseHealth()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Database health check complete!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Metric", "Value"},
				{"Status", res.Status},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// cleanupOrphanedCmd represents the cleanup command
var cleanupOrphanedCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Scan for and get details of orphaned records in the database",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Scanning for orphaned database records...")
		res, err := adminAPI.GetOrphanedRecords()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Scan complete!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Record Category", "Orphaned Count"},
				{"Orphaned Tool Instances", strconv.Itoa(res.OrphanedToolInstances)},
				{"Orphaned Tool Metadata", strconv.Itoa(res.OrphanedToolMetadata)},
				{"Orphaned Requests", strconv.Itoa(res.OrphanedRequests)},
				{"Orphaned Tool Requests", strconv.Itoa(res.OrphanedToolRequests)},
				{"Orphaned Request Steps", strconv.Itoa(res.OrphanedRequestSteps)},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

func init() {
	healthCmd.AddCommand(dbHealthCmd)
	healthCmd.AddCommand(cleanupOrphanedCmd)

	TeamToolboxCmd.AddCommand(healthCmd)
}
