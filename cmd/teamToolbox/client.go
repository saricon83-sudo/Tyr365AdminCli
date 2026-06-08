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
	cGroupId     string
	cUserId      string
	cToolId      int
	cTemplateId  int
	cRequestId   int
	cStatus      string
	cSubject     string
	cMessage     string
)

// getToolsForTeamCmd represents the getToolsForTeam command
var getToolsForTeamCmd = &cobra.Command{
	Use:   "getToolsForTeam",
	Short: "Get all available tools for a specific Team",
	Run: func(cmd *cobra.Command, args []string) {
		if cGroupId == "" {
			var err error
			cGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || cGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching tools for team %s...", cGroupId))
		tools, err := client.GetAllTblToolsForTeam(cGroupId)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to fetch tools: %v", err))
			return
		}
		spinner.Success("Tools retrieved successfully!")

		output.PrintResult(tools, func() { printPtermTools(tools) })
	},
}

// checkOwnershipCmd represents the checkOwnership command
var checkOwnershipCmd = &cobra.Command{
	Use:   "checkOwnership",
	Short: "Check if a user is owner of a Team (required before adding tools)",
	Run: func(cmd *cobra.Command, args []string) {
		if cGroupId == "" {
			var err error
			cGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || cGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}
		if cUserId == "" {
			var err error
			cUserId, err = pterm.DefaultInteractiveTextInput.Show("Enter User Principal Name (UPN/Email)")
			if err != nil || cUserId == "" {
				pterm.Error.Println("User ID/UPN is required")
				return
			}
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking ownership of user %s in team %s...", cUserId, cGroupId))
		isOwner, err := client.UserIsOwnerInTeam(cGroupId, cUserId)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to check ownership: %v", err))
			return
		}

		if isOwner {
			spinner.Success(fmt.Sprintf("User %s is an owner in team %s!", cUserId, cGroupId))
			pterm.Success.Println("User has ownership privileges.")
		} else {
			spinner.Warning(fmt.Sprintf("User %s is NOT an owner in team %s.", cUserId, cGroupId))
			pterm.Warning.Println("User does not have ownership privileges.")
		}
	},
}

// requestToolCmd represents the requestTool command
var requestToolCmd = &cobra.Command{
	Use:   "requestTool",
	Short: "Submit a tool provisioning request",
	Run: func(cmd *cobra.Command, args []string) {
		if cGroupId == "" {
			var err error
			cGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || cGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}
		if cToolId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID (integer)")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			cToolId, _ = strconv.Atoi(val)
		}
		if cTemplateId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Template ID (integer)")
			if err != nil {
				pterm.Error.Println("Template ID is required")
				return
			}
			cTemplateId, _ = strconv.Atoi(val)
		}
		if cUserId == "" {
			var err error
			cUserId, err = pterm.DefaultInteractiveTextInput.Show("Enter your UPN/Email (Initiated By)")
			if err != nil || cUserId == "" {
				pterm.Error.Println("Initiator email is required")
				return
			}
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		req := &teamToolboxHelper.TblToolRequest{
			GroupId:     cGroupId,
			ToolId:      cToolId,
			Status:      "Pending",
			InitiatedBy: cUserId,
			RequestData: fmt.Sprintf(`{"templateId": %d}`, cTemplateId),
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Submitting request for Tool %d...", cToolId))
		result, err := client.AddRequestForTool(req)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to submit request: %v", err))
			return
		}
		spinner.Success("Request submitted successfully!")

		pterm.Success.Printf("Created Request ID: %d (Status: %s)\n", result.Id, result.Status)
	},
}

// getRequestCmd represents the getRequest command
var getRequestCmd = &cobra.Command{
	Use:   "getRequest",
	Short: "Get details of a specific tool request",
	Run: func(cmd *cobra.Command, args []string) {
		if cRequestId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			cRequestId, _ = strconv.Atoi(val)
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching request details for ID %d...", cRequestId))
		req, err := client.GetRequestById(cRequestId)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to fetch request: %v", err))
			return
		}
		spinner.Success("Request details retrieved!")

		output.PrintResult(req, func() {
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"ID", strconv.Itoa(int(req.Id))},
				{"Group ID", req.GroupId},
				{"Tool ID", strconv.Itoa(int(req.ToolId))},
				{"Request Data", req.RequestData},
				{"Status", req.Status},
				{"Initiated By", req.InitiatedBy},
				{"Created", req.Created.Format("2006-01-02 15:04:05")},
				{"Modified", req.Modified.Format("2006-01-02 15:04:05")},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// updateRequestStatusCmd represents the updateRequestStatus command
var updateRequestStatusCmd = &cobra.Command{
	Use:   "updateRequestStatus",
	Short: "Update the status of a tool request (worker callback)",
	Run: func(cmd *cobra.Command, args []string) {
		if cRequestId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Request ID")
			if err != nil {
				pterm.Error.Println("Request ID is required")
				return
			}
			cRequestId, _ = strconv.Atoi(val)
		}
		if cStatus == "" {
			var err error
			cStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Pending", "Processing", "Completed", "Error", "Cancelled"}).Show("Select Status")
			if err != nil {
				pterm.Error.Println("Status is required")
				return
			}
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating request %d status to %s...", cRequestId, cStatus))
		resp, err := client.UpdateToolRequestStatus(cRequestId, cStatus)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to update status: %v", err))
			return
		}
		spinner.Success("Status updated successfully!")
		pterm.Info.Println(resp.Message)
	},
}

// addToolInstanceCmd represents the addToolInstance command
var addToolInstanceCmd = &cobra.Command{
	Use:   "addToolInstance",
	Short: "Register a provisioned tool instance",
	Run: func(cmd *cobra.Command, args []string) {
		if cGroupId == "" {
			var err error
			cGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || cGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}
		if cToolId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			cToolId, _ = strconv.Atoi(val)
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Registering tool instance %d for team %s...", cToolId, cGroupId))
		instance, err := client.AddInstanceOfToolToDb(cGroupId, cToolId)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to add tool instance: %v", err))
			return
		}
		spinner.Success("Tool instance registered successfully!")
		pterm.Success.Printf("Instance ID: %d (Group: %s, Tool: %d, Created: %s)\n",
			instance.Id, instance.GroupId, instance.ToolId, instance.Created.Format("2006-01-02 15:04:05"))
	},
}

// logMessageCmd represents the logMessage command
var logMessageCmd = &cobra.Command{
	Use:   "logMessage",
	Short: "Log a provisioning message",
	Run: func(cmd *cobra.Command, args []string) {
		if cSubject == "" {
			var err error
			cSubject, err = pterm.DefaultInteractiveTextInput.Show("Enter Subject/Topic")
			if err != nil || cSubject == "" {
				pterm.Error.Println("Subject is required")
				return
			}
		}
		if cMessage == "" {
			var err error
			cMessage, err = pterm.DefaultInteractiveTextInput.Show("Enter Message Body")
			if err != nil || cMessage == "" {
				pterm.Error.Println("Message body is required")
				return
			}
		}
		if cStatus == "" {
			var err error
			cStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Information", "Warning", "Error", "Success"}).Show("Select Message Severity/Status")
			if err != nil {
				pterm.Error.Println("Status is required")
				return
			}
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Logging provisioning message...")
		entry, err := client.LogMessage(cSubject, cMessage, cStatus)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to log message: %v", err))
			return
		}
		spinner.Success("Message logged successfully!")
		pterm.Info.Printf("Logged: [%s] Subject: %s - %s\n", entry.Status, entry.Subject, entry.Message)
	},
}

// listToolsCmd represents the listTools command
var listToolsCmd = &cobra.Command{
	Use:   "listTools",
	Short: "List all tools in the catalog",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching tools catalog...")
		tools, err := adminAPI.GetTools()
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to fetch tools catalog: %v", err))
			return
		}
		spinner.Success("Catalog loaded!")

		output.PrintResult(tools, func() { printPtermTools(tools) })
	},
}

func printPtermTools(tools []teamToolboxHelper.TblTool) {
	tableData := pterm.TableData{
		{"ID", "Tool Name", "Template ID", "Topic Name"},
	}
	for _, t := range tools {
		tableData = append(tableData, []string{
			strconv.Itoa(int(t.Id)),
			t.ToolName,
			strconv.Itoa(int(t.CurrentTempateId)),
			t.TopicName,
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	getToolsForTeamCmd.Flags().StringVar(&cGroupId, "groupId", "", "Group ID of the team")
	checkOwnershipCmd.Flags().StringVar(&cGroupId, "groupId", "", "Group ID of the team")
	checkOwnershipCmd.Flags().StringVar(&cUserId, "userId", "", "User ID / UPN of the user")
	requestToolCmd.Flags().StringVar(&cGroupId, "groupId", "", "Group ID of the team")
	requestToolCmd.Flags().IntVar(&cToolId, "toolId", 0, "Tool ID to provision")
	requestToolCmd.Flags().IntVar(&cTemplateId, "templateId", 0, "Template ID of the tool")
	requestToolCmd.Flags().StringVar(&cUserId, "userId", "", "UPN/Email of request initiator")
	getRequestCmd.Flags().IntVar(&cRequestId, "requestId", 0, "ID of the request to retrieve")
	updateRequestStatusCmd.Flags().IntVar(&cRequestId, "id", 0, "ID of the request to update")
	updateRequestStatusCmd.Flags().StringVar(&cStatus, "status", "", "New status to assign")
	addToolInstanceCmd.Flags().StringVar(&cGroupId, "groupId", "", "Group ID of the team")
	addToolInstanceCmd.Flags().IntVar(&cToolId, "toolId", 0, "Tool ID to register")
	logMessageCmd.Flags().StringVar(&cSubject, "subject", "", "Log subject")
	logMessageCmd.Flags().StringVar(&cMessage, "message", "", "Log message text")
	logMessageCmd.Flags().StringVar(&cStatus, "status", "", "Log message status")

	TeamToolboxCmd.AddCommand(getToolsForTeamCmd)
	TeamToolboxCmd.AddCommand(checkOwnershipCmd)
	TeamToolboxCmd.AddCommand(requestToolCmd)
	TeamToolboxCmd.AddCommand(getRequestCmd)
	TeamToolboxCmd.AddCommand(updateRequestStatusCmd)
	TeamToolboxCmd.AddCommand(addToolInstanceCmd)
	TeamToolboxCmd.AddCommand(logMessageCmd)
	TeamToolboxCmd.AddCommand(listToolsCmd)
}
