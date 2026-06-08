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
	rId            int
	rHours         int
	rMinRetries    int
	rCount         int
	rIncludeHidden bool
	rGroupId       string
	rEmail         string
	rEndpointName  string
	rStatus        string
	rMessage       string
	rPriority      int
	rHidden        bool
	rIds           string
	rStep          string

	// List filters
	rQueued    bool
	rRunning   bool
	rStuck     bool
	rHighRetry bool
	rSlowest   bool
	rErrors    bool
)

// requestsCmd represents the requests command group
var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "Manage and inspect provisioning requests",
}

// listRequestsCmd represents the list command
var listRequestsCmd = &cobra.Command{
	Use:   "list",
	Short: "List requests filtered by status or characteristics",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		if rQueued {
			spinner, _ := pterm.DefaultSpinner.Start("Fetching queued requests...")
			res, err := adminAPI.GetQueued()
			if err != nil {
				spinner.Fail(err.Error())
				return
			}
			spinner.Success("Queued requests loaded!")
			output.PrintResult(res, func() { printQueued(res) })
		} else if rRunning {
			spinner, _ := pterm.DefaultSpinner.Start("Fetching running requests...")
			res, err := adminAPI.GetRunning()
			if err != nil {
				spinner.Fail(err.Error())
				return
			}
			spinner.Success("Running requests loaded!")
			output.PrintResult(res, func() { printRunning(res) })
		} else if rStuck {
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests stuck for > %d hours...", rHours))
			res, err := adminAPI.GetStuck(rHours)
			if err != nil {
				spinner.Fail(err.Error())
				return
			}
			spinner.Success("Stuck requests loaded!")
			output.PrintResult(res, func() { printRequests(res) })
		} else if rHighRetry {
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests with > %d retries...", rMinRetries))
			res, err := adminAPI.GetHighRetry(rMinRetries)
			if err != nil {
				spinner.Fail(err.Error())
				return
			}
			spinner.Success("High retry requests loaded!")
			output.PrintResult(res, func() { printRequests(res) })
		} else if rSlowest {
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching slowest %d requests...", rCount))
			res, err := adminAPI.GetSlowest(rCount)
			if err != nil {
				spinner.Fail(err.Error())
				return
			}
			spinner.Success("Slowest requests loaded!")
			output.PrintResult(res, func() { printSlowest(res) })
		} else {
			// Default to Errors
			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching error requests (IncludeHidden: %v)...", rIncludeHidden))
			res, err := adminAPI.GetErrors(rIncludeHidden)
			if err != nil {
				spinner.Fail(err.Error())
				return
			}
			spinner.Success("Error requests loaded!")
			output.PrintResult(res, func() { printErrors(res) })
		}
	},
}

// getRequestDetailCmd represents the get command
var getRequestDetailCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a request by ID",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching details for request ID %d...", rId))
		req, err := adminAPI.GetRequest(rId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Request loaded!")

		tableData := pterm.TableData{
			{"Field", "Value"},
			{"ID", strconv.Itoa(req.Id)},
			{"Group ID", req.GroupId},
			{"Endpoint", req.Endpoint},
			{"Status", req.Status},
			{"Priority", strconv.Itoa(req.Priority)},
			{"Retry Count", strconv.Itoa(req.RetryCount)},
			{"Hidden", strconv.FormatBool(req.Hidden)},
			{"Message", req.Message},
			{"Initiated By", req.InitiatedBy},
			{"Created", req.Created.Format("2006-01-02 15:04:05")},
			{"Modified", req.Modified.Format("2006-01-02 15:04:05")},
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

		if len(req.RequestSteps) > 0 {
			fmt.Println("\nRequest Steps:")
			output.PrintResult(req.RequestSteps, func() { printRequestSteps(req.RequestSteps) })
		}
	},
}

// byGroupCmd represents the by-group command
var byGroupCmd = &cobra.Command{
	Use:   "by-group",
	Short: "Get requests by team Group ID",
	Run: func(cmd *cobra.Command, args []string) {
		if rGroupId == "" {
			var err error
			rGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || rGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests for group %s...", rGroupId))
		res, err := adminAPI.GetRequestsByGroup(rGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Requests loaded!")
		output.PrintResult(res, func() { printRequests(res) })
	},
}

// byInitiatorCmd represents the by-initiator command
var byInitiatorCmd = &cobra.Command{
	Use:   "by-initiator",
	Short: "Get requests by initiator email",
	Run: func(cmd *cobra.Command, args []string) {
		if rEmail == "" {
			var err error
			rEmail, err = pterm.DefaultInteractiveTextInput.Show("Enter Initiator Email")
			if err != nil || rEmail == "" {
				pterm.Error.Println("Email is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests for %s...", rEmail))
		res, err := adminAPI.GetRequestsByInitiator(rEmail)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Requests loaded!")
		output.PrintResult(res, func() { printRequests(res) })
	},
}

// byEndpointCmd represents the by-endpoint command
var byEndpointCmd = &cobra.Command{
	Use:   "by-endpoint",
	Short: "Get requests by target endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		if rEndpointName == "" {
			var err error
			rEndpointName, err = pterm.DefaultInteractiveTextInput.Show("Enter Endpoint Name")
			if err != nil || rEndpointName == "" {
				pterm.Error.Println("Endpoint name is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests for endpoint %s...", rEndpointName))
		res, err := adminAPI.GetRequestsByEndpoint(rEndpointName)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Requests loaded!")
		output.PrintResult(res, func() { printRequests(res) })
	},
}

// stepsCmd represents the steps command
var stepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Get steps of a specific request",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching steps for request %d...", rId))
		res, err := adminAPI.GetRequestSteps(rId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Steps loaded!")
		output.PrintResult(res, func() { printRequestSteps(res) })
	},
}

// updateStatusCmd represents the update-status command
var updateStatusCmd = &cobra.Command{
	Use:   "update-status",
	Short: "Update request status",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}
		if rStatus == "" {
			var err error
			rStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Pending", "Processing", "Completed", "Error", "Cancelled"}).Show("Select Status")
			if err != nil {
				pterm.Error.Println("Status is required")
				return
			}
		}
		if rMessage == "" {
			var err error
			rMessage, err = pterm.DefaultInteractiveTextInput.Show("Enter change message")
			if err != nil {
				pterm.Error.Println("Message is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating request %d status to %s...", rId, rStatus))
		resp, err := adminAPI.UpdateRequestStatus(rId, rStatus, rMessage)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Status updated!")
		pterm.Info.Println(resp.Message)
	},
}

// updatePriorityCmd represents the update-priority command
var updatePriorityCmd = &cobra.Command{
	Use:   "update-priority",
	Short: "Update request priority",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}
		if rPriority == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Priority (integer)")
			if err != nil {
				pterm.Error.Println("Priority is required")
				return
			}
			rPriority, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating request %d priority to %d...", rId, rPriority))
		resp, err := adminAPI.UpdateRequestPriority(rId, rPriority)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Priority updated!")
		pterm.Info.Println(resp.Message)
	},
}

// retryCmd represents the retry command
var retryCmd = &cobra.Command{
	Use:   "retry",
	Short: "Retry a failed request",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Queuing request %d for retry...", rId))
		resp, err := adminAPI.RetryRequest(rId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Request queued!")
		pterm.Info.Println(resp.Message)
	},
}

// updateHiddenCmd represents the update-hidden command
var updateHiddenCmd = &cobra.Command{
	Use:   "update-hidden",
	Short: "Hide/Unhide a request",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Setting request %d hidden to %v...", rId, rHidden))
		resp, err := adminAPI.UpdateRequestHidden(rId, rHidden)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Hidden setting updated!")
		pterm.Info.Println(resp.Message)
	},
}

// bulkRetryCmd represents the bulk-retry command
var bulkRetryCmd = &cobra.Command{
	Use:   "bulk-retry",
	Short: "Bulk retry multiple requests",
	Run: func(cmd *cobra.Command, args []string) {
		if rIds == "" {
			var err error
			rIds, err = pterm.DefaultInteractiveTextInput.Show("Enter comma-separated Request IDs (e.g. 1,2,3)")
			if err != nil || rIds == "" {
				pterm.Error.Println("IDs are required")
				return
			}
		}

		parts := strings.Split(rIds, ",")
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

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Bulk retrying %d requests...", len(ids)))
		res, err := adminAPI.BulkRetryRequests(ids)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Bulk operation completed!")

		tableData := pterm.TableData{
			{"Metric", "Value"},
			{"Total Requested", strconv.Itoa(res.TotalRequested)},
			{"Succeeded", strconv.Itoa(res.Succeeded)},
			{"Failed", strconv.Itoa(res.Failed)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

		if len(res.Errors) > 0 {
			pterm.Warning.Println("\nErrors:")
			errTable := pterm.TableData{{"Request ID", "Error Message"}}
			for _, e := range res.Errors {
				errTable = append(errTable, []string{strconv.Itoa(e.Id), e.Error})
			}
			pterm.DefaultTable.WithHasHeader().WithData(errTable).Render()
		}
	},
}

// bulkHideCmd represents the bulk-hide command
var bulkHideCmd = &cobra.Command{
	Use:   "bulk-hide",
	Short: "Bulk hide/unhide multiple requests",
	Run: func(cmd *cobra.Command, args []string) {
		if rIds == "" {
			var err error
			rIds, err = pterm.DefaultInteractiveTextInput.Show("Enter comma-separated Request IDs")
			if err != nil || rIds == "" {
				pterm.Error.Println("IDs are required")
				return
			}
		}

		parts := strings.Split(rIds, ",")
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

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Bulk setting hidden = %v on %d requests...", rHidden, len(ids)))
		res, err := adminAPI.BulkHideRequests(ids, rHidden)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Bulk hide completed!")

		tableData := pterm.TableData{
			{"Metric", "Value"},
			{"Total Requested", strconv.Itoa(res.TotalRequested)},
			{"Succeeded", strconv.Itoa(res.Succeeded)},
			{"Failed", strconv.Itoa(res.Failed)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	},
}

// addStepCmd represents the add-step command
var addStepCmd = &cobra.Command{
	Use:   "add-step",
	Short: "Add a step to a request (for manual intervention logging)",
	Run: func(cmd *cobra.Command, args []string) {
		if rId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			rId, _ = strconv.Atoi(val)
		}
		if rStep == "" {
			var err error
			rStep, err = pterm.DefaultInteractiveTextInput.Show("Enter Step Name (e.g. ManualIntervention)")
			if err != nil || rStep == "" {
				pterm.Error.Println("Step name is required")
				return
			}
		}
		if rStatus == "" {
			var err error
			rStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Completed", "Failed", "Warning"}).Show("Select Step Status")
			if err != nil {
				pterm.Error.Println("Status is required")
				return
			}
		}
		if rMessage == "" {
			var err error
			rMessage, err = pterm.DefaultInteractiveTextInput.Show("Enter message detail")
			if err != nil {
				pterm.Error.Println("Message is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Adding request step...")
		res, err := adminAPI.AddRequestStep(rId, rStep, rStatus, rMessage)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Step added!")

		pterm.Success.Printf("Step Added: %s (Status: %s, Message: %s)\n", res.Step, res.Status, res.Message)
	},
}

// Helper printer functions using pterm
func printRequests(requests []teamToolboxHelper.Request) {
	if len(requests) == 0 {
		pterm.Info.Println("No requests found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Endpoint", "Status", "Priority", "Retries", "Created"},
	}
	for _, r := range requests {
		tableData = append(tableData, []string{
			strconv.Itoa(r.Id),
			r.GroupId,
			r.Endpoint,
			r.Status,
			strconv.Itoa(r.Priority),
			strconv.Itoa(r.RetryCount),
			r.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printQueued(requests []teamToolboxHelper.ViewQueuedRequest) {
	if len(requests) == 0 {
		pterm.Info.Println("No queued requests found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Endpoint", "Status", "Priority", "Initiated By", "Created"},
	}
	for _, r := range requests {
		tableData = append(tableData, []string{
			strconv.Itoa(r.Id),
			r.GroupId,
			r.Endpoint,
			r.Status,
			strconv.Itoa(r.Priority),
			r.InitiatedBy,
			r.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printRunning(requests []teamToolboxHelper.ViewRunningRequest) {
	if len(requests) == 0 {
		pterm.Info.Println("No running requests found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Endpoint", "Status", "Initiated By", "Created", "Modified"},
	}
	for _, r := range requests {
		tableData = append(tableData, []string{
			strconv.Itoa(r.Id),
			r.GroupId,
			r.Endpoint,
			r.Status,
			r.InitiatedBy,
			r.Created.Format("2006-01-02 15:04"),
			r.Modified.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printErrors(requests []teamToolboxHelper.ViewErrorRequest) {
	if len(requests) == 0 {
		pterm.Info.Println("No error requests found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Endpoint", "Status", "Retries", "Hidden", "Created"},
	}
	for _, r := range requests {
		tableData = append(tableData, []string{
			strconv.Itoa(r.Id),
			r.GroupId,
			r.Endpoint,
			r.Status,
			strconv.Itoa(r.RetryCount),
			strconv.FormatBool(r.Hidden),
			r.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printSlowest(requests []teamToolboxHelper.RequestDurationInfo) {
	if len(requests) == 0 {
		pterm.Info.Println("No records found.")
		return
	}
	tableData := pterm.TableData{
		{"Request ID", "Group ID", "Endpoint", "Status", "Duration (Min)", "Created"},
	}
	for _, r := range requests {
		tableData = append(tableData, []string{
			strconv.Itoa(r.RequestId),
			r.GroupId,
			r.Endpoint,
			r.Status,
			fmt.Sprintf("%.1f", r.DurationMinutes),
			r.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printRequestSteps(steps []teamToolboxHelper.RequestStep) {
	if len(steps) == 0 {
		pterm.Info.Println("No steps recorded.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Step", "Status", "Message", "Created"},
	}
	for _, s := range steps {
		tableData = append(tableData, []string{
			strconv.Itoa(s.Id),
			s.Step,
			s.Status,
			s.Message,
			s.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	listRequestsCmd.Flags().BoolVar(&rQueued, "queued", false, "Filter queued requests")
	listRequestsCmd.Flags().BoolVar(&rRunning, "running", false, "Filter running requests")
	listRequestsCmd.Flags().BoolVar(&rStuck, "stuck", false, "Filter stuck requests")
	listRequestsCmd.Flags().IntVar(&rHours, "hours", 2, "Hours threshold for stuck requests (default: 2)")
	listRequestsCmd.Flags().BoolVar(&rHighRetry, "high-retry", false, "Filter requests with high retry counts")
	listRequestsCmd.Flags().IntVar(&rMinRetries, "min-retries", 3, "Minimum retry count for high-retry (default: 3)")
	listRequestsCmd.Flags().BoolVar(&rSlowest, "slowest", false, "Filter slowest requests by duration")
	listRequestsCmd.Flags().IntVar(&rCount, "count", 20, "Number of results for slowest (default: 20)")
	listRequestsCmd.Flags().BoolVar(&rIncludeHidden, "include-hidden", false, "Include hidden requests in error list")

	getRequestDetailCmd.Flags().IntVar(&rId, "id", 0, "Request ID")
	byGroupCmd.Flags().StringVar(&rGroupId, "groupId", "", "Group GUID")
	byInitiatorCmd.Flags().StringVar(&rEmail, "email", "", "Initiator UPN / Email")
	byEndpointCmd.Flags().StringVar(&rEndpointName, "name", "", "Endpoint name")
	stepsCmd.Flags().IntVar(&rId, "id", 0, "Request ID")

	updateStatusCmd.Flags().IntVar(&rId, "id", 0, "Request ID")
	updateStatusCmd.Flags().StringVar(&rStatus, "status", "", "New status")
	updateStatusCmd.Flags().StringVar(&rMessage, "message", "", "Reason message")

	updatePriorityCmd.Flags().IntVar(&rId, "id", 0, "Request ID")
	updatePriorityCmd.Flags().IntVar(&rPriority, "priority", 0, "Priority value")

	retryCmd.Flags().IntVar(&rId, "id", 0, "Request ID")

	updateHiddenCmd.Flags().IntVar(&rId, "id", 0, "Request ID")
	updateHiddenCmd.Flags().BoolVar(&rHidden, "hidden", false, "Hidden status")

	bulkRetryCmd.Flags().StringVar(&rIds, "ids", "", "Comma-separated request IDs")
	bulkHideCmd.Flags().StringVar(&rIds, "ids", "", "Comma-separated request IDs")
	bulkHideCmd.Flags().BoolVar(&rHidden, "hidden", false, "Hidden status")

	addStepCmd.Flags().IntVar(&rId, "id", 0, "Request ID")
	addStepCmd.Flags().StringVar(&rStep, "step", "", "Step name")
	addStepCmd.Flags().StringVar(&rStatus, "status", "", "Step status")
	addStepCmd.Flags().StringVar(&rMessage, "message", "", "Step message")

	requestsCmd.AddCommand(listRequestsCmd)
	requestsCmd.AddCommand(getRequestDetailCmd)
	requestsCmd.AddCommand(byGroupCmd)
	requestsCmd.AddCommand(byInitiatorCmd)
	requestsCmd.AddCommand(byEndpointCmd)
	requestsCmd.AddCommand(stepsCmd)
	requestsCmd.AddCommand(updateStatusCmd)
	requestsCmd.AddCommand(updatePriorityCmd)
	requestsCmd.AddCommand(retryCmd)
	requestsCmd.AddCommand(updateHiddenCmd)
	requestsCmd.AddCommand(bulkRetryCmd)
	requestsCmd.AddCommand(bulkHideCmd)
	requestsCmd.AddCommand(addStepCmd)

	TeamToolboxCmd.AddCommand(requestsCmd)
}
