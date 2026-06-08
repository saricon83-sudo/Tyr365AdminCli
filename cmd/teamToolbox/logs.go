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
	lSubject string
	lStatus  string
	lFrom    string
	lTo      string
	lLimit   int
	lCount   int
	lMessage string
)

// logsCmd represents the logs command group
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Inspect and create log entries",
}

// listLogsCmd represents the list command
var listLogsCmd = &cobra.Command{
	Use:   "list",
	Short: "List logs with optional filters for subject, status, dates, and limits",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching logs...")
		res, err := adminAPI.GetLogs(lSubject, lStatus, lFrom, lTo, lLimit)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Logs loaded!")

		output.PrintResult(res, func() { printLoggerEntries(res) })
	},
}

// recentLogsCmd represents the recent command
var recentLogsCmd = &cobra.Command{
	Use:   "recent",
	Short: "Get N most recent log entries",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching recent %d logs...", lCount))
		res, err := adminAPI.GetRecentLogs(lCount)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Recent logs loaded!")

		output.PrintResult(res, func() { printLoggerEntries(res) })
	},
}

// logsBySubjectCmd represents the by-subject command
var logsBySubjectCmd = &cobra.Command{
	Use:   "by-subject",
	Short: "Get log entries matching a specific subject",
	Run: func(cmd *cobra.Command, args []string) {
		if lSubject == "" {
			var err error
			lSubject, err = pterm.DefaultInteractiveTextInput.Show("Enter Subject/Topic")
			if err != nil || lSubject == "" {
				pterm.Error.Println("Subject is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching logs for subject '%s'...", lSubject))
		res, err := adminAPI.GetLogsBySubject(lSubject)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Logs loaded!")

		output.PrintResult(res, func() { printLoggerEntries(res) })
	},
}

// addLogCmd represents the add command
var addLogCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new log entry manually",
	Run: func(cmd *cobra.Command, args []string) {
		if lSubject == "" {
			var err error
			lSubject, err = pterm.DefaultInteractiveTextInput.Show("Enter Subject/Topic")
			if err != nil || lSubject == "" {
				pterm.Error.Println("Subject is required")
				return
			}
		}
		if lMessage == "" {
			var err error
			lMessage, err = pterm.DefaultInteractiveTextInput.Show("Enter Log Message")
			if err != nil || lMessage == "" {
				pterm.Error.Println("Message is required")
				return
			}
		}
		if lStatus == "" {
			var err error
			lStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Information", "Warning", "Error", "Success"}).Show("Select Status/Severity")
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

		entry := teamToolboxHelper.LogEntry{
			Subject: lSubject,
			Message: lMessage,
			Status:  lStatus,
		}

		spinner, _ := pterm.DefaultSpinner.Start("Adding log entry...")
		res, err := adminAPI.AddLog(entry)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Log entry added successfully!")

		pterm.Success.Printf("Registered Log: ID: %d, Subject: %s, Message: %s (Status: %s)\n", res.Id, res.Subject, res.Message, res.Status)
	},
}

func printLoggerEntries(logs []teamToolboxHelper.TblToolBoxLogger) {
	if len(logs) == 0 {
		pterm.Info.Println("No logs found matching criteria.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Subject", "Message", "Status", "Created"},
	}
	for _, l := range logs {
		tableData = append(tableData, []string{
			strconv.Itoa(l.Id),
			l.Subject,
			l.Message,
			l.Status,
			l.Created.Format("2006-01-02 15:04:05"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	listLogsCmd.Flags().StringVar(&lSubject, "subject", "", "Filter by subject")
	listLogsCmd.Flags().StringVar(&lStatus, "status", "", "Filter by status")
	listLogsCmd.Flags().StringVar(&lFrom, "from", "", "Filter from date (YYYY-MM-DD)")
	listLogsCmd.Flags().StringVar(&lTo, "to", "", "Filter to date (YYYY-MM-DD)")
	listLogsCmd.Flags().IntVar(&lLimit, "limit", 100, "Limit number of results (default: 100)")

	recentLogsCmd.Flags().IntVar(&lCount, "count", 25, "Number of recent logs to fetch (default: 25)")

	logsBySubjectCmd.Flags().StringVar(&lSubject, "subject", "", "Subject keyword")

	addLogCmd.Flags().StringVar(&lSubject, "subject", "", "Log subject")
	addLogCmd.Flags().StringVar(&lMessage, "message", "", "Log message text")
	addLogCmd.Flags().StringVar(&lStatus, "status", "", "Log severity/status")

	logsCmd.AddCommand(listLogsCmd)
	logsCmd.AddCommand(recentLogsCmd)
	logsCmd.AddCommand(logsBySubjectCmd)
	logsCmd.AddCommand(addLogCmd)

	TeamToolboxCmd.AddCommand(logsCmd)
}
