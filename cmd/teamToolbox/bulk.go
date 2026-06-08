package teamToolboxCmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var (
	bIds    string
	bStatus string
	bToolId int
)

// bulkCmd represents the bulk command group
var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Perform bulk administration actions on requests and tools",
}

// bulkUpdateStatusCmd represents the update-status command
var bulkUpdateStatusCmd = &cobra.Command{
	Use:   "update-status",
	Short: "Bulk update status for multiple requests",
	Run: func(cmd *cobra.Command, args []string) {
		if bIds == "" {
			var err error
			bIds, err = pterm.DefaultInteractiveTextInput.Show("Enter comma-separated Request IDs")
			if err != nil || bIds == "" {
				pterm.Error.Println("IDs are required")
				return
			}
		}
		if bStatus == "" {
			var err error
			bStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Pending", "Processing", "Completed", "Error", "Cancelled"}).Show("Select New Status")
			if err != nil {
				pterm.Error.Println("Status is required")
				return
			}
		}

		parts := strings.Split(bIds, ",")
		var ids []int
		for _, part := range parts {
			val := strings.TrimSpace(part)
			if id, err := strconv.Atoi(val); err == nil {
				ids = append(ids, id)
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Bulk updating status to %s for %d requests...", bStatus, len(ids)))
		res, err := adminAPI.BulkUpdateStatus(ids, bStatus)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printBulkResult(res) })
	},
}

// bulkEnableToolsCmd represents the enable-tools command
var bulkEnableToolsCmd = &cobra.Command{
	Use:   "enable-tools",
	Short: "Bulk enable multiple tools in the catalog",
	Run: func(cmd *cobra.Command, args []string) {
		if bIds == "" {
			var err error
			bIds, err = pterm.DefaultInteractiveTextInput.Show("Enter comma-separated Tool IDs")
			if err != nil || bIds == "" {
				pterm.Error.Println("Tool IDs are required")
				return
			}
		}

		parts := strings.Split(bIds, ",")
		var ids []int
		for _, part := range parts {
			val := strings.TrimSpace(part)
			if id, err := strconv.Atoi(val); err == nil {
				ids = append(ids, id)
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Bulk enabling %d tools...", len(ids)))
		res, err := adminAPI.BulkEnableTools(ids)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printBulkResult(res) })
	},
}

// bulkDisableToolsCmd represents the disable-tools command
var bulkDisableToolsCmd = &cobra.Command{
	Use:   "disable-tools",
	Short: "Bulk disable multiple tools in the catalog",
	Run: func(cmd *cobra.Command, args []string) {
		if bIds == "" {
			var err error
			bIds, err = pterm.DefaultInteractiveTextInput.Show("Enter comma-separated Tool IDs")
			if err != nil || bIds == "" {
				pterm.Error.Println("Tool IDs are required")
				return
			}
		}

		parts := strings.Split(bIds, ",")
		var ids []int
		for _, part := range parts {
			val := strings.TrimSpace(part)
			if id, err := strconv.Atoi(val); err == nil {
				ids = append(ids, id)
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Bulk disabling %d tools...", len(ids)))
		res, err := adminAPI.BulkDisableTools(ids)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printBulkResult(res) })
	},
}

// bulkRetryFailedForToolCmd represents the retry-failed-for-tool command
var bulkRetryFailedForToolCmd = &cobra.Command{
	Use:   "retry-failed-for-tool",
	Short: "Retry all failed requests for a specific tool",
	Run: func(cmd *cobra.Command, args []string) {
		if bToolId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			bToolId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Retrying all failed requests for tool %d...", bToolId))
		res, err := adminAPI.BulkRetryFailedForTool(bToolId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		output.PrintResult(res, func() { printBulkResult(res) })
	},
}

func printBulkResult(res *teamToolboxHelper.BulkOperationResult) {
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Total Requested", strconv.Itoa(res.TotalRequested)},
		{"Succeeded", strconv.Itoa(res.Succeeded)},
		{"Failed", strconv.Itoa(res.Failed)},
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	if len(res.Errors) > 0 {
		pterm.Warning.Println("\nErrors:")
		errTable := pterm.TableData{{"ID", "Error Message"}}
		for _, e := range res.Errors {
			errTable = append(errTable, []string{strconv.Itoa(e.Id), e.Error})
		}
		pterm.DefaultTable.WithHasHeader().WithData(errTable).Render()
	}
}

func init() {
	bulkUpdateStatusCmd.Flags().StringVar(&bIds, "ids", "", "Comma-separated request IDs")
	bulkUpdateStatusCmd.Flags().StringVar(&bStatus, "status", "", "New status")

	bulkEnableToolsCmd.Flags().StringVar(&bIds, "ids", "", "Comma-separated tool IDs")
	bulkDisableToolsCmd.Flags().StringVar(&bIds, "ids", "", "Comma-separated tool IDs")

	bulkRetryFailedForToolCmd.Flags().IntVar(&bToolId, "toolId", 0, "Tool ID")

	bulkCmd.AddCommand(bulkUpdateStatusCmd)
	bulkCmd.AddCommand(bulkEnableToolsCmd)
	bulkCmd.AddCommand(bulkDisableToolsCmd)
	bulkCmd.AddCommand(bulkRetryFailedForToolCmd)

	TeamToolboxCmd.AddCommand(bulkCmd)
}
