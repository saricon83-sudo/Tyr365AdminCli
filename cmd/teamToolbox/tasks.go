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
	tskStatus    int
	tskJobType   int
	tskProjectNo string
)

// tasksCmd represents the tasks command group
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Inspect background tasks in Team Toolbox",
}

// listTasksCmd represents the list command
var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all background tasks with status and jobType filters",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching tasks...")
		res, err := adminAPI.GetTasks(tskStatus, tskJobType)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tasks loaded!")

		output.PrintResult(res, func() { printTasksList(res) })
	},
}

// pendingTasksCmd represents the pending command
var pendingTasksCmd = &cobra.Command{
	Use:   "pending",
	Short: "Get all pending background tasks",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching pending tasks...")
		res, err := adminAPI.GetPendingTasks()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Pending tasks loaded!")

		output.PrintResult(res, func() { printTasksList(res) })
	},
}

// failedTasksCmd represents the failed command
var failedTasksCmd = &cobra.Command{
	Use:   "failed",
	Short: "Get all failed background tasks",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching failed tasks...")
		res, err := adminAPI.GetFailedTasks()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Failed tasks loaded!")

		output.PrintResult(res, func() { printTasksList(res) })
	},
}

// tasksByProjectCmd represents the by-project command
var tasksByProjectCmd = &cobra.Command{
	Use:   "by-project",
	Short: "Get background tasks filtered by project number",
	Run: func(cmd *cobra.Command, args []string) {
		if tskProjectNo == "" {
			var err error
			tskProjectNo, err = pterm.DefaultInteractiveTextInput.Show("Enter Project Number")
			if err != nil || tskProjectNo == "" {
				pterm.Error.Println("Project number is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching tasks for project %s...", tskProjectNo))
		res, err := adminAPI.GetTasksByProject(tskProjectNo)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tasks loaded!")

		output.PrintResult(res, func() { printTasksList(res) })
	},
}

func printTasksList(tasks []teamToolboxHelper.TblTask) {
	if len(tasks) == 0 {
		pterm.Info.Println("No background tasks found matching criteria.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Project No", "Status", "Job Type", "Message", "Created"},
	}
	for _, t := range tasks {
		tableData = append(tableData, []string{
			strconv.Itoa(t.Id),
			t.ProjectNo,
			strconv.Itoa(t.Status),
			strconv.Itoa(t.JobType),
			t.Message,
			t.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	listTasksCmd.Flags().IntVar(&tskStatus, "status", 0, "Filter by status code")
	listTasksCmd.Flags().IntVar(&tskJobType, "jobType", 0, "Filter by job type code")
	tasksByProjectCmd.Flags().StringVar(&tskProjectNo, "projectNo", "", "Project number")

	tasksCmd.AddCommand(listTasksCmd)
	tasksCmd.AddCommand(pendingTasksCmd)
	tasksCmd.AddCommand(failedTasksCmd)
	tasksCmd.AddCommand(tasksByProjectCmd)

	TeamToolboxCmd.AddCommand(tasksCmd)
}
