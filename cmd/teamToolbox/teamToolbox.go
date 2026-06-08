/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamToolboxCmd

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

// TeamToolboxCmd represents the teamToolbox command
var TeamToolboxCmd = &cobra.Command{
	Use:   "teamToolbox",
	Short: "Manage Team Toolbox – tools, requests, instances, teams, jobs and more",
	Long: `Interactive admin console for the Team Toolbox API.

Run without arguments to open the interactive menu, or call any
subcommand directly (e.g. 365Admin teamToolbox requests list --queued).`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMenu()
	},
}

// ─── menu entry ──────────────────────────────────────────────────────────────

type menuEntry struct {
	label  string
	cobCmd *cobra.Command
}

type menuGroup struct {
	label    string
	children []menuEntry
}

// ─── root menu ───────────────────────────────────────────────────────────────

func runInteractiveMenu() {
	// Print banner
	ptermLogo, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Team", pterm.NewStyle(pterm.FgLightCyan)),
		putils.LettersFromStringWithStyle("Toolbox", pterm.NewStyle(pterm.FgLightMagenta)),
	).Srender()
	pterm.DefaultCenter.Println(ptermLogo)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		pterm.LightBlue("Interactive Admin Console") + "\n",
	)

	groups := buildMenuGroups()

	for {
		// ── Level 1: group selection ──────────────────────────────────────────
		groupLabels := make([]string, len(groups))
		for i, g := range groups {
			groupLabels[i] = g.label
		}
		groupLabels = append(groupLabels, "🚪 Exit")

		selectedGroup, err := pterm.DefaultInteractiveSelect.
			WithOptions(groupLabels).
			WithDefaultText("Select a command group").
			Show()
		if err != nil || selectedGroup == "🚪 Exit" {
			pterm.Info.Println("Goodbye!")
			return
		}

		var chosen *menuGroup
		for i := range groups {
			if groups[i].label == selectedGroup {
				chosen = &groups[i]
				break
			}
		}
		if chosen == nil {
			continue
		}

		// ── Level 2: action selection ─────────────────────────────────────────
		childLabels := make([]string, len(chosen.children))
		for i, c := range chosen.children {
			childLabels[i] = c.label
		}
		childLabels = append(childLabels, "⬅️  Back")

		selectedChild, err := pterm.DefaultInteractiveSelect.
			WithOptions(childLabels).
			WithDefaultText("Select an action").
			Show()
		if err != nil || selectedChild == "⬅️  Back" {
			continue
		}

		var childCmd *cobra.Command
		for _, c := range chosen.children {
			if c.label == selectedChild {
				childCmd = c.cobCmd
				break
			}
		}
		if childCmd == nil || childCmd.Run == nil {
			pterm.Warning.Println("Command is not directly runnable (it is a group). Use the CLI flags directly.")
			continue
		}

		// ── Execute ───────────────────────────────────────────────────────────
		pterm.DefaultSection.Printf("▶  %s\n", childCmd.Short)
		// The Run handler already prompts for any missing parameters.
		childCmd.Run(childCmd, []string{})

		// Pause so the output is readable before the menu redraws.
		pterm.Println()
		_, _ = pterm.DefaultInteractiveTextInput.
			WithDefaultText("Press Enter to return to the menu").
			Show()
		pterm.Println()
	}
}

// buildMenuGroups returns the full two-level menu structure.
func buildMenuGroups() []menuGroup {
	return []menuGroup{
		{
			label: "📊 Dashboard & Stats",
			children: []menuEntry{
				{"📈 Dashboard overview", getDashboard},
				{"📋 Requests by status", statsStatusCmd},
				{"💾 Storage released", statsStorageCmd},
				{"📦 Archive job stats", statsArchiveCmd},
				{"⏳ Pending counts", statsPendingCmd},
				{"📅 Requests by day", getRequestsByDay},
				{"📊 Tool adoption stats", toolAdoption},
			},
		},
		{
			label: "🔁 Requests",
			children: []menuEntry{
				{"📋 List requests (filter: queued/running/stuck/errors…)", listRequestsCmd},
				{"🔍 Get request details by ID", getRequestDetailCmd},
				{"🪜 Get request steps", stepsCmd},
				{"🔄 Retry a request", retryCmd},
				{"✏️  Update request status", updateStatusCmd},
				{"🎚️  Update request priority", updatePriorityCmd},
				{"👁️  Toggle request hidden", updateHiddenCmd},
				{"➕ Add step log to request", addStepCmd},
				{"🔀 Bulk retry requests", bulkRetryCmd},
				{"🙈 Bulk hide/show requests", bulkHideCmd},
				{"👥 Requests by group", byGroupCmd},
				{"👤 Requests by initiator", byInitiatorCmd},
				{"🔗 Requests by endpoint", byEndpointCmd},
			},
		},
		{
			label: "🔧 Tools",
			children: []menuEntry{
				{"📋 List all tools", listToolsAdminCmd},
				{"🔍 Get full tool details", getToolDetailsCmd},
				{"🔍 Get tool by ID (interactive)", getToolByIdCmd},
				{"📦 Instances for a tool", toolInstancesCmd},
				{"🔁 Requests for a tool", toolRequestsCmd},
				{"🚫 Unused tools", unusedToolsCmd},
				{"✅ Enable / Disable tool", updateToolEnabledCmd},
				{"🏷️  Update tool topic", updateToolTopicCmd},
				{"🔢 Update tool template", updateToolTemplateCmd},
				{"✏️  Update tool (name/desc/url)", updateToolCmd},
				{"➕ Add new tool to catalog", addToolCmd},
				{"📋 Requests per tool (legacy)", getRequestPerTool},
				{"📜 Rules & Logic (legacy)", getRulesAndLogicCmd},
			},
		},
		{
			label: "📦 Instances",
			children: []menuEntry{
				{"📋 List instances (filter by tool/group)", listInstancesCmd},
				{"🔍 Get instance details", getInstanceCmd},
				{"📋 Instance metadata", getInstanceMetadataCmd},
				{"⬆️  Outdated instances", outdatedInstancesCmd},
				{"🔗 Orphaned instances", orphanedInstancesCmd},
				{"➕ Register new instance", addInstanceCmd},
				{"🗑️  Delete instance", deleteInstanceCmd},
			},
		},
		{
			label: "👥 Teams & Groups",
			children: []menuEntry{
				{"📋 List managed teams", listTeamsCmd},
				{"🔍 Get team record", getTeamCmd},
				{"📊 Full team details", getTeamDetailsCmd},
				{"📦 Team's tool instances", getTeamInstancesCmd},
				{"🔁 Team's recent requests", getTeamRequestsCmd},
				{"✏️  Update team status", updateTeamStatusCmd},
				{"🔍 Search teams", searchTeamsCmd},
				{"🚫 Teams without tools", teamsWithoutToolsCmd},
				{"🗂️  Get group record", getGroupCmd},
				{"🗂️  Get groups view", getGroupsViewCmd},
			},
		},
		{
			label: "🗄️  Jobs",
			children: []menuEntry{
				{"📋 List archive jobs", getArchiveJobsCmd},
				{"🔍 Get archive job details", getArchiveJobCmd},
				{"✏️  Update archive job status", updateArchiveJobStatusCmd},
				{"📋 Required archive jobs", getRequiredArchiveJobsCmd},
				{"📤 Export data jobs", getExportJobsCmd},
				{"🧹 Clear site jobs", getClearSiteJobsCmd},
				{"📊 Clear site jobs summary", getClearSiteJobsSummaryCmd},
			},
		},
		{
			label: "📋 Logs",
			children: []menuEntry{
				{"📋 List all logs", listLogsCmd},
				{"⏱️  Recent logs", recentLogsCmd},
				{"🔍 Logs by subject", logsBySubjectCmd},
				{"➕ Add log entry", addLogCmd},
			},
		},
		{
			label: "📐 Rules & Requirements",
			children: []menuEntry{
				{"📋 List all rules", listRulesCmd},
				{"🔍 Get rule by ID", getRuleCmd},
				{"🚫 Unused rules", unusedRulesCmd},
				{"📋 All rule logics", getRuleLogicsCmd},
				{"🔍 Logics for a tool", getRuleLogicsByToolCmd},
				{"📋 Extended requirements for tool", getExtendedRequirementsByToolCmd},
			},
		},
		{
			label: "⚙️  Background Tasks",
			children: []menuEntry{
				{"📋 List all tasks", listTasksCmd},
				{"⏳ Pending tasks", pendingTasksCmd},
				{"❌ Failed tasks", failedTasksCmd},
				{"🔍 Tasks by project", tasksByProjectCmd},
			},
		},
		{
			label: "🗺️  GeoBIM",
			children: []menuEntry{
				{"📋 List all GeoBIM records", listGeoBIMCmd},
				{"🔍 GeoBIM by group", getGeoBIMCmd},
			},
		},
		{
			label: "❤️  Health & Cleanup",
			children: []menuEntry{
				{"🩺 Database health check", dbHealthCmd},
				{"🧹 Cleanup orphaned records", cleanupOrphanedCmd},
			},
		},
		{
			label: "🔀 Bulk Operations",
			children: []menuEntry{
				{"✏️  Bulk update request status", bulkUpdateStatusCmd},
				{"✅ Bulk enable tools", bulkEnableToolsCmd},
				{"🚫 Bulk disable tools", bulkDisableToolsCmd},
				{"🔄 Bulk retry failed for tool", bulkRetryFailedForToolCmd},
			},
		},
		{
			label: "📱 Client APIs (end-user facing)",
			children: []menuEntry{
				{"📋 Get tools available for team", getToolsForTeamCmd},
				{"👑 Check team ownership", checkOwnershipCmd},
				{"🔨 Request tool provisioning", requestToolCmd},
				{"🔍 Get tool request by ID", getRequestCmd},
				{"✏️  Update tool request status", updateRequestStatusCmd},
				{"➕ Register tool instance (client)", addToolInstanceCmd},
				{"📝 Log message", logMessageCmd},
				{"📋 List tools (client view)", listToolsCmd},
			},
		},
		{
			label: "📜 Legacy / Standalone Commands",
			children: []menuEntry{
				{"⚠️  Get error requests (legacy)", getErrors},
				{"📊 Get requests per tool (legacy)", getRequestPerTool},
				{"📜 Get rules & logic (legacy)", getRulesAndLogicCmd},
			},
		},
	}
}

func init() {
	// No extra flags needed on the root teamToolbox command itself.
}
