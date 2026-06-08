package teamToolboxCmd

import (
	"fmt"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var (
	tgGroupId string
	tgStatus  string
	tgOrigin  string
	tgQuery   string
)

// teamsCmd represents the teams command group
var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage and inspect managed teams",
}

// groupsCmd represents the groups command group
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Inspect group records and views",
}

// listTeamsCmd represents the list command
var listTeamsCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed teams with optional status and origin filters",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching managed teams...")
		res, err := adminAPI.GetTeams(tgStatus, tgOrigin)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Teams loaded!")

		output.PrintResult(res, func() { printManagedTeams(res) })
	},
}

// getTeamCmd represents the get command
var getTeamCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific managed team",
	Run: func(cmd *cobra.Command, args []string) {
		if tgGroupId == "" {
			var err error
			tgGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || tgGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching team ID %s...", tgGroupId))
		team, err := adminAPI.GetTeam(tgGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Team loaded!")

		output.PrintResult(team, func() {
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"Group ID", team.GroupId},
				{"Team Name", team.TeamName},
				{"Project No", team.ProjectNo},
				{"Project Name", team.ProjectName},
				{"Status", team.Status},
				{"Origin", team.Origin},
				{"Retention", team.Retention},
				{"Site ID", team.SiteId},
				{"URL", team.Url},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getTeamDetailsCmd represents the details command
var getTeamDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get full details including instances and recent requests for a team",
	Run: func(cmd *cobra.Command, args []string) {
		if tgGroupId == "" {
			var err error
			tgGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || tgGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching full team details for %s...", tgGroupId))
		details, err := adminAPI.GetTeamFullDetails(tgGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Full details loaded!")

		output.PrintResult(details, func() {
			pterm.DefaultSection.Println("General Information")
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"Group ID", details.GroupId},
				{"Team Name", details.TeamName},
				{"Project No", details.ProjectNo},
				{"Project Name", details.ProjectName},
				{"Status", details.Status},
				{"Origin", details.Origin},
				{"Retention", details.Retention},
				{"Site ID", details.SiteId},
				{"URL", details.Url},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			if len(details.ToolInstances) > 0 {
				pterm.DefaultSection.Println("\nInstalled Tool Instances")
				printInstances(details.ToolInstances)
			}

			if len(details.RecentRequests) > 0 {
				pterm.DefaultSection.Println("\nRecent Requests")
				printRequests(details.RecentRequests)
			}

			if len(details.RecentToolRequests) > 0 {
				pterm.DefaultSection.Println("\nRecent Tool Requests")
				printToolRequests(details.RecentToolRequests)
			}
		})
	},
}

// getTeamInstancesCmd represents the instances command
var getTeamInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Get tool instances installed in a team",
	Run: func(cmd *cobra.Command, args []string) {
		if tgGroupId == "" {
			var err error
			tgGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || tgGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching instances for team %s...", tgGroupId))
		res, err := adminAPI.GetTeamInstances(tgGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Instances loaded!")

		output.PrintResult(res, func() { printInstances(res) })
	},
}

// getTeamRequestsCmd represents the requests command
var getTeamRequestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "Get requests related to a team",
	Run: func(cmd *cobra.Command, args []string) {
		if tgGroupId == "" {
			var err error
			tgGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || tgGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching requests for team %s...", tgGroupId))
		res, err := adminAPI.GetTeamRequests(tgGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Requests loaded!")

		output.PrintResult(res, func() { printRequests(res) })
	},
}

// updateTeamStatusCmd represents the status command
var updateTeamStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Update team status",
	Run: func(cmd *cobra.Command, args []string) {
		if tgGroupId == "" {
			var err error
			tgGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || tgGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}
		if tgStatus == "" {
			var err error
			tgStatus, err = pterm.DefaultInteractiveSelect.WithOptions([]string{"Active", "Archived", "Deleted", "Locked"}).Show("Select Status")
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

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Updating team %s status to %s...", tgGroupId, tgStatus))
		resp, err := adminAPI.UpdateTeamStatus(tgGroupId, tgStatus)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Team status updated!")
		pterm.Info.Println(resp.Message)
	},
}

// teamsWithoutToolsCmd represents the without-tools command
var teamsWithoutToolsCmd = &cobra.Command{
	Use:   "without-tools",
	Short: "Get teams without any tool instances installed",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching teams without tools...")
		res, err := adminAPI.GetTeamsWithoutTools()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Teams loaded!")

		output.PrintResult(res, func() { printManagedTeams(res) })
	},
}

// searchTeamsCmd represents the search command
var searchTeamsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search teams by query string (name, group id, project number, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if tgQuery == "" {
			var err error
			tgQuery, err = pterm.DefaultInteractiveTextInput.Show("Enter Search Query")
			if err != nil || tgQuery == "" {
				pterm.Error.Println("Search query is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Searching for '%s'...", tgQuery))
		res, err := adminAPI.SearchTeams(tgQuery)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Search complete!")

		if len(res) == 0 {
			pterm.Info.Println("No matching teams found.")
			return
		}

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Group ID", "Team Name", "Project No", "Project Name", "Status", "Origin", "Matched Field"},
			}
			for _, t := range res {
				tableData = append(tableData, []string{
					t.GroupId,
					t.TeamName,
					t.ProjectNo,
					t.ProjectName,
					t.Status,
					t.Origin,
					t.MatchedField,
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getGroupCmd represents the get command under groups
var getGroupCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific group record by ID",
	Run: func(cmd *cobra.Command, args []string) {
		if tgGroupId == "" {
			var err error
			tgGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || tgGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching group record %s...", tgGroupId))
		grp, err := adminAPI.GetGroup(tgGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Group loaded!")

		output.PrintResult(grp, func() {
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"Group ID", grp.GroupId},
				{"Team Name", grp.TeamName},
				{"Project No", grp.ProjectNo},
				{"Project Name", grp.ProjectName},
				{"Status", grp.Status},
				{"Origin", grp.Origin},
				{"Created", grp.Created.Format("2006-01-02 15:04:05")},
				{"Modified", grp.Modified.Format("2006-01-02 15:04:05")},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getGroupsViewCmd represents the view command under groups
var getGroupsViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Get all groups from ViewGroups",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching groups view...")
		res, err := adminAPI.GetGroupsView()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Groups view loaded!")

		if len(res) == 0 {
			pterm.Info.Println("No groups returned in view.")
			return
		}

		output.PrintResult(res, func() {
			tableData := pterm.TableData{
				{"Group ID", "Team Name", "Project No", "Project Name", "Status", "Origin"},
			}
			for _, g := range res {
				tableData = append(tableData, []string{
					g.GroupId,
					g.TeamName,
					g.ProjectNo,
					g.ProjectName,
					g.Status,
					g.Origin,
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// Helper printers
func printManagedTeams(teams []teamToolboxHelper.ManagedTeam) {
	if len(teams) == 0 {
		pterm.Info.Println("No managed teams found.")
		return
	}
	tableData := pterm.TableData{
		{"Group ID", "Team Name", "Project No", "Status", "Origin", "URL"},
	}
	for _, t := range teams {
		tableData = append(tableData, []string{
			t.GroupId,
			t.TeamName,
			t.ProjectNo,
			t.Status,
			t.Origin,
			t.Url,
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	listTeamsCmd.Flags().StringVar(&tgStatus, "status", "", "Filter by status")
	listTeamsCmd.Flags().StringVar(&tgOrigin, "origin", "", "Filter by origin")

	getTeamCmd.Flags().StringVar(&tgGroupId, "groupId", "", "Group ID")
	getTeamDetailsCmd.Flags().StringVar(&tgGroupId, "groupId", "", "Group ID")
	getTeamInstancesCmd.Flags().StringVar(&tgGroupId, "groupId", "", "Group ID")
	getTeamRequestsCmd.Flags().StringVar(&tgGroupId, "groupId", "", "Group ID")

	updateTeamStatusCmd.Flags().StringVar(&tgGroupId, "groupId", "", "Group ID")
	updateTeamStatusCmd.Flags().StringVar(&tgStatus, "status", "", "New status")

	searchTeamsCmd.Flags().StringVar(&tgQuery, "query", "", "Search query text")

	getGroupCmd.Flags().StringVar(&tgGroupId, "groupId", "", "Group ID")

	teamsCmd.AddCommand(listTeamsCmd)
	teamsCmd.AddCommand(getTeamCmd)
	teamsCmd.AddCommand(getTeamDetailsCmd)
	teamsCmd.AddCommand(getTeamInstancesCmd)
	teamsCmd.AddCommand(getTeamRequestsCmd)
	teamsCmd.AddCommand(updateTeamStatusCmd)
	teamsCmd.AddCommand(teamsWithoutToolsCmd)
	teamsCmd.AddCommand(searchTeamsCmd)

	groupsCmd.AddCommand(getGroupCmd)
	groupsCmd.AddCommand(getGroupsViewCmd)

	TeamToolboxCmd.AddCommand(teamsCmd)
	TeamToolboxCmd.AddCommand(groupsCmd)
}
