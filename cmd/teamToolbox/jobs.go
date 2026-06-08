package teamToolboxCmd

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var (
	jId      int
	jStatus  string
	jGroupId string
)

// jobsCmd represents the jobs command group
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage and inspect archive, export, and clean site jobs",
}

// getArchiveJobsCmd represents the archive-list command
var getArchiveJobsCmd = &cobra.Command{
	Use:   "archive-list",
	Short: "Get archive jobs with optional status filter",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching archive jobs...")
		res, err := adminAPI.GetArchiveJobs(jStatus)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printArchiveJobs(res) })
	},
}

// getArchiveJobCmd represents the archive-get command
var getArchiveJobCmd = &cobra.Command{
	Use:   "archive-get",
	Short: "Get a specific archive job by ID along with its sub-jobs",
	Run: func(cmd *cobra.Command, args []string) {
		if jId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Archive Job ID")
			if err != nil {
				pterm.Error.Println("Job ID is required")
				return
			}
			jId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching archive job ID %d...", jId))
		job, err := adminAPI.GetArchiveJob(jId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Archive job loaded!")

		output.PrintResult(job, func() {
			pterm.DefaultSection.Println("Job Information")
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"ID", strconv.Itoa(job.Id)},
				{"Group ID", job.GroupId},
				{"Job Type", job.JobType},
				{"Message", job.Message},
				{"Status", job.Status},
				{"Created", job.Created.Format("2006-01-02 15:04:05")},
				{"Modified", job.Modified.Format("2006-01-02 15:04:05")},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			if len(job.ArchiveSubJobs) > 0 {
				pterm.DefaultSection.Println("\nArchive Sub-Jobs")
				subTable := pterm.TableData{
					{"ID", "SubJob Type", "Status", "Message", "Modified"},
				}
				for _, sj := range job.ArchiveSubJobs {
					subTable = append(subTable, []string{
						strconv.Itoa(sj.Id),
						sj.SubJobType,
						sj.Status,
						sj.Message,
						sj.Modified.Format("2006-01-02 15:04"),
					})
				}
				pterm.DefaultTable.WithHasHeader().WithData(subTable).Render()
			}
		})
	},
}

// updateArchiveJobStatusCmd represents the archive-status command
var updateArchiveJobStatusCmd = &cobra.Command{
	Use:   "archive-status",
	Short: "Update archive job status",
	Run: func(cmd *cobra.Command, args []string) {
		if jId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Archive Job ID")
			if err != nil {
				pterm.Error.Println("Job ID is required")
				return
			}
			jId, _ = strconv.Atoi(val)
		}
		if jStatus == "" {
			var err error
			jStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Pending", "Running", "Completed", "Failed", "Cancelled"}).Show("Select Status")
			if err != nil {
				pterm.Error.Println("Status is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating archive job %d status to %s...", jId, jStatus))
		resp, err := adminAPI.UpdateArchiveJobStatus(jId, jStatus)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Status updated!")
		pterm.Info.Println(resp.Message)
	},
}

// getRequiredArchiveJobsCmd represents the archive-required command
var getRequiredArchiveJobsCmd = &cobra.Command{
	Use:   "archive-required",
	Short: "Get all required archive jobs",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching required archive jobs...")
		res, err := adminAPI.GetRequiredArchiveJobs()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Required jobs loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Group ID", "Team Name", "Retention", "Expiry Date", "Status"},
			}
			for _, r := range res {
				tableData = append(tableData, []string{
					r.GroupId,
					r.TeamName,
					r.Retention,
					r.ExpiryDate.Format("2006-01-02 15:04"),
					r.Status,
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getExportJobsCmd represents the export-list command
var getExportJobsCmd = &cobra.Command{
	Use:   "export-list",
	Short: "Get export jobs with optional Group ID filter",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching export data jobs...")
		res, err := adminAPI.GetExportJobs(jGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Export jobs loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"ID", "Request ID", "Group ID", "Status", "File Path", "File Size (Bytes)", "Message", "Modified"},
			}
			for _, e := range res {
				tableData = append(tableData, []string{
					strconv.Itoa(e.Id),
					strconv.Itoa(e.RequestId),
					e.GroupId,
					e.Status,
					e.FilePath,
					strconv.FormatInt(e.FileSize, 10),
					e.Message,
					e.Modified.Format("2006-01-02 15:04"),
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getClearSiteJobsCmd represents the clearsite-list command
var getClearSiteJobsCmd = &cobra.Command{
	Use:   "clearsite-list",
	Short: "Get clear site cleanup jobs",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching clear site jobs...")
		res, err := adminAPI.GetClearSiteJobs(jStatus)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Cleanup jobs loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"ID", "Group ID", "Status", "Storage Released (Bytes)", "Files Deleted", "Message", "Modified"},
			}
			for _, c := range res {
				tableData = append(tableData, []string{
					strconv.Itoa(c.Id),
					c.GroupId,
					c.Status,
					strconv.FormatInt(c.StorageReleased, 10),
					strconv.Itoa(c.FilesDeleted),
					c.Message,
					c.Modified.Format("2006-01-02 15:04"),
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getClearSiteJobsSummaryCmd represents the clearsite-summary command
var getClearSiteJobsSummaryCmd = &cobra.Command{
	Use:   "clearsite-summary",
	Short: "Get clear site storage and job statistics summary",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching clear site jobs summary...")
		res, err := adminAPI.GetClearSiteJobsSummary()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Summary loaded!")

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Metric", "Value"},
				{"Total Storage Released", formatBytes(res.TotalStorageReleased)},
				{"Total Files Deleted", strconv.Itoa(res.TotalFilesDeleted)},
				{"Total Jobs", strconv.Itoa(res.TotalJobs)},
				{"Completed Jobs", strconv.Itoa(res.CompletedJobs)},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// Helper formats
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func printArchiveJobs(jobs []teamToolboxHelper.ArchiveJob) {
	if len(jobs) == 0 {
		pterm.Info.Println("No archive jobs found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Job Type", "Status", "Message", "Created"},
	}
	for _, j := range jobs {
		tableData = append(tableData, []string{
			strconv.Itoa(j.Id),
			j.GroupId,
			j.JobType,
			j.Status,
			j.Message,
			j.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	getArchiveJobsCmd.Flags().StringVar(&jStatus, "status", "", "Filter by status")
	getArchiveJobCmd.Flags().IntVar(&jId, "id", 0, "Archive Job ID")

	updateArchiveJobStatusCmd.Flags().IntVar(&jId, "id", 0, "Archive Job ID")
	updateArchiveJobStatusCmd.Flags().StringVar(&jStatus, "status", "", "New status")

	getExportJobsCmd.Flags().StringVar(&jGroupId, "groupId", "", "Group ID GUID")
	getClearSiteJobsCmd.Flags().StringVar(&jStatus, "status", "", "Filter by status")

	jobsCmd.AddCommand(getArchiveJobsCmd)
	jobsCmd.AddCommand(getArchiveJobCmd)
	jobsCmd.AddCommand(updateArchiveJobStatusCmd)
	jobsCmd.AddCommand(getRequiredArchiveJobsCmd)
	jobsCmd.AddCommand(getExportJobsCmd)
	jobsCmd.AddCommand(getClearSiteJobsCmd)
	jobsCmd.AddCommand(getClearSiteJobsSummaryCmd)

	TeamToolboxCmd.AddCommand(jobsCmd)
}
