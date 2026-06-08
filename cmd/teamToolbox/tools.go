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
	tId                int
	tEnabled           bool
	tTopic             string
	tTemplateId        int
	tName              string
	tDesc              string
	tUrl               string
	tRequiresArchiving bool
	tDays              int
)

// toolsCmd represents the tools command group
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Manage and inspect tools in the catalog",
}

// listToolsAdminCmd represents the list command
var listToolsAdminCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tools in the database",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching tools database...")
		tools, err := adminAPI.GetTools()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tools loaded!")

		if len(tools) == 0 {
			pterm.Info.Println("No tools found.")
			return
		}

		output.PrintResult(tools, func() { printPtermTools(tools) })
	},
}

// getToolDetailsCmd represents the get command
var getToolDetailsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get full details of a specific tool",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching details for tool ID %d...", tId))
		tool, err := adminAPI.GetToolFullDetails(tId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool details loaded!")

		output.PrintResult(tool, func() {
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"ID", strconv.Itoa(tool.Id)},
				{"Tool Name", tool.ToolName},
				{"Description", tool.ToolDescription},
				{"Topic Name", tool.TopicName},
				{"Info Page URL", tool.InfoPageUrl},
				{"Form Template", tool.FormTemplate},
				{"Current Template ID", strconv.Itoa(int(tool.CurrentTemplateId))},
				{"Enabled", strconv.FormatBool(tool.Enabled)},
				{"Requires Archiving", strconv.FormatBool(tool.RequiresArchiving)},
				{"Instance Count", strconv.Itoa(tool.InstanceCount)},
				{"Request Count", strconv.Itoa(tool.RequestCount)},
			}

			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			if len(tool.Rules) > 0 {
				pterm.Println("\nRules and logic for this tool:")
				ruleTable := pterm.TableData{
					{"Rule ID", "Rule Name", "Logic", "Value"},
				}
				for _, r := range tool.Rules {
					ruleTable = append(ruleTable, []string{
						strconv.Itoa(r.RuleId),
						r.RuleName,
						r.Logic,
						r.Value,
					})
				}
				pterm.DefaultTable.WithHasHeader().WithData(ruleTable).Render()
			}

			if len(tool.ExtendedRequirements) > 0 {
				pterm.Println("\nExtended requirements for this tool:")
				reqTable := pterm.TableData{
					{"Requirement Name", "Value"},
				}
				for _, r := range tool.ExtendedRequirements {
					reqTable = append(reqTable, []string{
						r.RequirementName,
						r.RequirementValue,
					})
				}
				pterm.DefaultTable.WithHasHeader().WithData(reqTable).Render()
			}
		})
	},
}

// updateToolEnabledCmd represents the update-enabled command
var updateToolEnabledCmd = &cobra.Command{
	Use:   "update-enabled",
	Short: "Enable or disable a tool",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Setting tool %d enabled status to %v...", tId, tEnabled))
		resp, err := adminAPI.UpdateToolEnabled(tId, tEnabled)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool status updated!")
		pterm.Info.Println(resp.Message)
	},
}

// updateToolTopicCmd represents the update-topic command
var updateToolTopicCmd = &cobra.Command{
	Use:   "update-topic",
	Short: "Update tool topic name",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}
		if tTopic == "" {
			var err error
			tTopic, err = pterm.DefaultInteractiveTextInput.Show("Enter New Topic Name")
			if err != nil || tTopic == "" {
				pterm.Error.Println("Topic is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating tool %d topic to %s...", tId, tTopic))
		resp, err := adminAPI.UpdateToolTopic(tId, tTopic)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool topic updated!")
		pterm.Info.Println(resp.Message)
	},
}

// updateToolTemplateCmd represents the update-template command
var updateToolTemplateCmd = &cobra.Command{
	Use:   "update-template",
	Short: "Update tool template version ID",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}
		if tTemplateId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Template ID")
			if err != nil {
				pterm.Error.Println("Template ID is required")
				return
			}
			tTemplateId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating tool %d template version to %d...", tId, tTemplateId))
		resp, err := adminAPI.UpdateToolTemplate(tId, tTemplateId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool template updated!")
		pterm.Info.Println(resp.Message)
	},
}

// updateToolCmd represents the update command
var updateToolCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tool properties",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID to update")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		// Load current details
		spinner, _ := pterm.DefaultSpinner.Start("Loading current tool details...")
		curr, err := adminAPI.GetToolFullDetails(tId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Loaded current details!")

		// Use flag values if modified, else prompt or use current values
		if !cmd.Flags().Changed("name") {
			tName, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue(curr.ToolName).Show("Tool Name")
		}
		if !cmd.Flags().Changed("desc") {
			tDesc, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue(curr.ToolDescription).Show("Tool Description")
		}
		if !cmd.Flags().Changed("topic") {
			tTopic, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue(curr.TopicName).Show("Topic Name")
		}
		if !cmd.Flags().Changed("url") {
			tUrl, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue(curr.InfoPageUrl).Show("Info Page URL")
		}
		if !cmd.Flags().Changed("templateId") {
			val, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(strconv.Itoa(int(curr.CurrentTemplateId))).Show("Template ID")
			tTemplateId, _ = strconv.Atoi(val)
		}
		if !cmd.Flags().Changed("enabled") {
			sel, _ := pterm.DefaultInteractiveSelect.WithOptions([]string{"True", "False"}).WithDefaultOption(strconv.FormatBool(curr.Enabled)).Show("Enabled")
			tEnabled = (sel == "True")
		}
		if !cmd.Flags().Changed("requires-archiving") {
			sel, _ := pterm.DefaultInteractiveSelect.WithOptions([]string{"True", "False"}).WithDefaultOption(strconv.FormatBool(curr.RequiresArchiving)).Show("Requires Archiving")
			tRequiresArchiving = (sel == "True")
		}

		updateDto := teamToolboxHelper.ToolUpdateDto{
			ToolName:          tName,
			ToolDescription:   tDesc,
			TopicName:         tTopic,
			InfoPageUrl:       tUrl,
			FormTemplate:      curr.FormTemplate,
			CurrentTemplateId: tTemplateId,
			Enabled:           tEnabled,
			RequiresArchiving: tRequiresArchiving,
		}

		spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Updating tool %d in catalog...", tId))
		resp, err := adminAPI.UpdateTool(tId, updateDto)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool updated!")
		pterm.Info.Println(resp.Message)
	},
}

// toolInstancesCmd represents the instances subcommand
var toolInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Get all instances of a specific tool",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching instances for tool %d...", tId))
		res, err := adminAPI.GetToolInstances(tId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Instances loaded!")

		output.PrintResult(res, func() { printInstances(res) })
	},
}

// toolRequestsCmd represents the requests subcommand
var toolRequestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "Get all requests for a specific tool",
	Run: func(cmd *cobra.Command, args []string) {
		if tId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			tId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests for tool %d...", tId))
		res, err := adminAPI.GetToolRequestsForTool(tId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool requests loaded!")

		output.PrintResult(res, func() { printToolRequests(res) })
	},
}

// unusedToolsCmd represents the unused command
var unusedToolsCmd = &cobra.Command{
	Use:   "unused",
	Short: "Get tools with no requests in the last N days",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching unused tools in last %d days...", tDays))
		tools, err := adminAPI.GetUnusedTools(tDays)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Unused tools loaded!")

		output.PrintResult(tools, func() { printPtermTools(tools) })
	},
}

// addToolCmd represents the add command
var addToolCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new tool to the catalog",
	Run: func(cmd *cobra.Command, args []string) {
		if tName == "" {
			var err error
			tName, err = pterm.DefaultInteractiveTextInput.Show("Enter Tool Name")
			if err != nil || tName == "" {
				pterm.Error.Println("Tool Name is required")
				return
			}
		}
		if tTopic == "" {
			var err error
			tTopic, err = pterm.DefaultInteractiveTextInput.Show("Enter Topic Name")
			if err != nil || tTopic == "" {
				pterm.Error.Println("Topic Name is required")
				return
			}
		}
		if tTemplateId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Current Template ID")
			if err != nil {
				pterm.Error.Println("Template ID is required")
				return
			}
			tTemplateId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		newTool := teamToolboxHelper.TblTool{
			ToolName:         tName,
			TopicName:        tTopic,
			CurrentTempateId: int32(tTemplateId),
		}

		spinner, _ := pterm.DefaultSpinner.Start("Adding tool to catalog...")
		res, err := adminAPI.AddTool(newTool)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool added successfully!")

		pterm.Success.Printf("Registered Tool: ID: %d, Name: %s (Topic: %s)\n", res.Id, res.ToolName, res.TopicName)
	},
}

// Helper printers
func printInstances(instances []teamToolboxHelper.TblToolInstance) {
	if len(instances) == 0 {
		pterm.Info.Println("No instances found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Tool ID", "Template Version", "Created"},
	}
	for _, inst := range instances {
		tableData = append(tableData, []string{
			strconv.Itoa(inst.Id),
			inst.GroupId,
			strconv.Itoa(inst.ToolId),
			strconv.Itoa(inst.TemplateVersion),
			inst.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func printToolRequests(reqs []teamToolboxHelper.TblToolRequest) {
	if len(reqs) == 0 {
		pterm.Info.Println("No tool requests found.")
		return
	}
	tableData := pterm.TableData{
		{"ID", "Group ID", "Tool ID", "Request Data", "Status", "Initiated By", "Created"},
	}
	for _, r := range reqs {
		tableData = append(tableData, []string{
			strconv.Itoa(int(r.Id)),
			r.GroupId,
			strconv.Itoa(int(r.ToolId)),
			r.RequestData,
			r.Status,
			r.InitiatedBy,
			r.Created.Format("2006-01-02 15:04"),
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	getToolDetailsCmd.Flags().IntVar(&tId, "id", 0, "Tool ID")
	updateToolEnabledCmd.Flags().IntVar(&tId, "id", 0, "Tool ID")
	updateToolEnabledCmd.Flags().BoolVar(&tEnabled, "enabled", false, "Enable status (true/false)")

	updateToolTopicCmd.Flags().IntVar(&tId, "id", 0, "Tool ID")
	updateToolTopicCmd.Flags().StringVar(&tTopic, "topic", "", "New topic name")

	updateToolTemplateCmd.Flags().IntVar(&tId, "id", 0, "Tool ID")
	updateToolTemplateCmd.Flags().IntVar(&tTemplateId, "templateId", 0, "New template ID")

	updateToolCmd.Flags().IntVar(&tId, "id", 0, "Tool ID to update")
	updateToolCmd.Flags().StringVar(&tName, "name", "", "New name")
	updateToolCmd.Flags().StringVar(&tDesc, "desc", "", "New description")
	updateToolCmd.Flags().StringVar(&tTopic, "topic", "", "New topic")
	updateToolCmd.Flags().StringVar(&tUrl, "url", "", "New URL")
	updateToolCmd.Flags().IntVar(&tTemplateId, "templateId", 0, "New template version")
	updateToolCmd.Flags().BoolVar(&tEnabled, "enabled", false, "Enable status")
	updateToolCmd.Flags().BoolVar(&tRequiresArchiving, "requires-archiving", false, "Requires archiving status")

	toolInstancesCmd.Flags().IntVar(&tId, "id", 0, "Tool ID")
	toolRequestsCmd.Flags().IntVar(&tId, "id", 0, "Tool ID")
	unusedToolsCmd.Flags().IntVar(&tDays, "days", 90, "Number of days without requests (default: 90)")

	addToolCmd.Flags().StringVar(&tName, "name", "", "Tool Name")
	addToolCmd.Flags().StringVar(&tTopic, "topic", "", "Topic Name")
	addToolCmd.Flags().IntVar(&tTemplateId, "templateId", 0, "Current Template ID")

	toolsCmd.AddCommand(listToolsAdminCmd)
	toolsCmd.AddCommand(getToolDetailsCmd)
	toolsCmd.AddCommand(updateToolEnabledCmd)
	toolsCmd.AddCommand(updateToolTopicCmd)
	toolsCmd.AddCommand(updateToolTemplateCmd)
	toolsCmd.AddCommand(updateToolCmd)
	toolsCmd.AddCommand(toolInstancesCmd)
	toolsCmd.AddCommand(toolRequestsCmd)
	toolsCmd.AddCommand(unusedToolsCmd)
	toolsCmd.AddCommand(addToolCmd)

	TeamToolboxCmd.AddCommand(toolsCmd)
}
