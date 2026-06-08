package teamToolboxCmd

import (
	"strconv"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command group
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Inspect system, tool, storage, and archive statistics",
}

// statsStatusCmd represents the status subcommand under stats
var statsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get total requests grouped by status",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching requests by status...")
		res, err := adminAPI.GetRequestsByStatus()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Stats loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Request Status", "Count"},
			}
			for status, count := range res {
				tableData = append(tableData, []string{
					status,
					strconv.Itoa(count),
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// statsStorageCmd represents the storage subcommand under stats
var statsStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Get storage released summary from clean site operations",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching storage released summary...")
		res, err := adminAPI.GetStorageReleased()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Stats loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Metric", "Value"},
				{"Total Storage Released", formatBytes(res.TotalStorageReleased)},
				{"Total Files Deleted", strconv.Itoa(res.TotalFilesDeleted)},
				{"Total Jobs", strconv.Itoa(res.TotalJobs)},
				{"Completed Jobs", strconv.Itoa(res.CompletedJobs)},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			if len(res.ByPeriod) > 0 {
				pterm.DefaultSection.Println("\nReleased Storage By Period:")
				periodTable := pterm.TableData{
					{"Period", "Grain", "Storage Released"},
				}
				for _, p := range res.ByPeriod {
					periodTable = append(periodTable, []string{
						p.Period,
						p.Grain,
						formatBytes(p.StorageReleased),
					})
				}
				pterm.DefaultTable.WithHasHeader().WithData(periodTable).Render()
			}
		})
	},
}

// statsArchiveCmd represents the archive subcommand under stats
var statsArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Get archive job statistics",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching archive job stats...")
		res, err := adminAPI.GetArchiveJobStats()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Stats loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Metric", "Value"},
				{"Total Archive Jobs", strconv.Itoa(res.Total)},
				{"Pending Jobs", strconv.Itoa(res.Pending)},
				{"Running Jobs", strconv.Itoa(res.Running)},
				{"Completed Jobs", strconv.Itoa(res.Completed)},
				{"Failed Jobs", strconv.Itoa(res.Failed)},
				{"Sub-Jobs Total", strconv.Itoa(res.SubJobsTotal)},
				{"Sub-Jobs Pending", strconv.Itoa(res.SubJobsPending)},
				{"Sub-Jobs Completed", strconv.Itoa(res.SubJobsCompleted)},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// statsPendingCmd represents the pending subcommand under stats
var statsPendingCmd = &cobra.Command{
	Use:   "pending",
	Short: "Get pending counts across all job types",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching pending job counts...")
		res, err := adminAPI.GetPendingCounts()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Stats loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Job Category", "Pending", "Running"},
				{"Requests", strconv.Itoa(res.PendingRequests), strconv.Itoa(res.RunningRequests)},
				{"Tool Requests", strconv.Itoa(res.PendingToolRequests), "-"},
				{"Archive Jobs", strconv.Itoa(res.PendingArchiveJobs), "-"},
				{"Export Jobs", strconv.Itoa(res.PendingExportJobs), "-"},
				{"Clear Site Jobs", strconv.Itoa(res.PendingClearSiteJobs), "-"},
				{"Tasks", strconv.Itoa(res.PendingTasks), "-"},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

func init() {
	statsCmd.AddCommand(statsStatusCmd)
	statsCmd.AddCommand(statsStorageCmd)
	statsCmd.AddCommand(statsArchiveCmd)
	statsCmd.AddCommand(statsPendingCmd)

	TeamToolboxCmd.AddCommand(statsCmd)
}
