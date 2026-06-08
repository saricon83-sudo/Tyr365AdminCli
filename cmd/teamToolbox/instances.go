package teamToolboxCmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var (
	iId         int
	iToolId     int
	iGroupId    string
	iTemplateId int
)

// instancesCmd represents the instances command group
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage and inspect tool instances",
}

// listInstancesCmd represents the list command
var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all tool instances with optional filtering",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching tool instances...")
		res, err := adminAPI.GetInstances(iToolId, iGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Instances loaded!")
		output.PrintResult(res, func() { printInstances(res) })
	},
}

// getInstanceCmd represents the get command
var getInstanceCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific tool instance",
	Run: func(cmd *cobra.Command, args []string) {
		if iId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Instance ID")
			if err != nil {
				pterm.Error.Println("Instance ID is required")
				return
			}
			iId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching instance ID %d...", iId))
		inst, err := adminAPI.GetInstance(iId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Instance details loaded!")

		output.PrintResult(inst, func() {
			tableData := pterm.TableData{
				{"Field", "Value"},
				{"ID", strconv.Itoa(inst.Id)},
				{"Group ID", inst.GroupId},
				{"Tool ID", strconv.Itoa(inst.ToolId)},
				{"Template Version", strconv.Itoa(inst.TemplateVersion)},
				{"Created", inst.Created.Format("2006-01-02 15:04:05")},
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// getInstanceMetadataCmd represents the metadata command
var getInstanceMetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Get metadata for a specific tool instance",
	Run: func(cmd *cobra.Command, args []string) {
		if iId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Instance ID")
			if err != nil {
				pterm.Error.Println("Instance ID is required")
				return
			}
			iId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching metadata for instance ID %d...", iId))
		meta, err := adminAPI.GetInstanceMetadata(iId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Metadata loaded!")

		if len(meta) == 0 {
			pterm.Info.Println("No metadata keys found for this instance.")
			return
		}

		output.PrintResult(meta, func() {
			tableData := pterm.TableData{
				{"ID", "Meta Key", "Meta Value"},
			}
			for _, m := range meta {
				tableData = append(tableData, []string{
					strconv.Itoa(m.Id),
					m.Key,
					m.Value,
				})
			}
			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		})
	},
}

// deleteInstanceCmd represents the delete command
var deleteInstanceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tool instance",
	Run: func(cmd *cobra.Command, args []string) {
		if iId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Instance ID to delete")
			if err != nil {
				pterm.Error.Println("Instance ID is required")
				return
			}
			iId, _ = strconv.Atoi(val)
		}

		// Double check confirm
		confirm, err := pterm.DefaultInteractiveConfirm.Show(fmt.Sprintf("Are you sure you want to delete instance %d?", iId))
		if err != nil || !confirm {
			pterm.Info.Println("Deletion cancelled.")
			return
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting instance ID %d...", iId))
		resp, err := adminAPI.DeleteInstance(iId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Instance deleted successfully!")
		pterm.Info.Println(resp.Message)
	},
}

// orphanedInstancesCmd represents the orphaned command
var orphanedInstancesCmd = &cobra.Command{
	Use:   "orphaned",
	Short: "Get orphaned instances (group no longer exists)",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching orphaned tool instances...")
		res, err := adminAPI.GetOrphanedInstances()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Orphaned instances loaded!")

		output.PrintResult(res, func() { printInstances(res) })
	},
}

// outdatedInstancesCmd represents the outdated command
var outdatedInstancesCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Get instances with outdated template versions",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching outdated tool instances...")
		res, err := adminAPI.GetOutdatedInstances()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Outdated instances loaded!")

		output.PrintResult(res, func() { printInstances(res) })
	},
}

// addInstanceCmd represents the add command
var addInstanceCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new tool instance manually",
	Run: func(cmd *cobra.Command, args []string) {
		if iGroupId == "" {
			var err error
			iGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || iGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}
		if iToolId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Tool ID (integer)")
			if err != nil {
				pterm.Error.Println("Tool ID is required")
				return
			}
			iToolId, _ = strconv.Atoi(val)
		}
		if iTemplateId == 0 {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Template ID (integer)")
			if err != nil {
				pterm.Error.Println("Template ID is required")
				return
			}
			iTemplateId, _ = strconv.Atoi(val)
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		newInstance := teamToolboxHelper.TblToolInstance{
			GroupId:         iGroupId,
			ToolId:          iToolId,
			TemplateVersion: iTemplateId,
			Created:         teamToolboxHelper.CustomTime{Time: time.Now()},
		}

		spinner, _ := pterm.DefaultSpinner.Start("Registering new tool instance...")
		res, err := adminAPI.AddInstance(newInstance)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Instance registered!")

		pterm.Success.Printf("Registered Instance ID: %d (Group: %s, Tool: %d, Template Version: %d)\n", res.Id, res.GroupId, res.ToolId, res.TemplateVersion)
	},
}

func init() {
	listInstancesCmd.Flags().IntVar(&iToolId, "toolId", 0, "Filter by tool ID")
	listInstancesCmd.Flags().StringVar(&iGroupId, "groupId", "", "Filter by group ID")

	getInstanceCmd.Flags().IntVar(&iId, "id", 0, "Instance ID")
	getInstanceMetadataCmd.Flags().IntVar(&iId, "id", 0, "Instance ID")
	deleteInstanceCmd.Flags().IntVar(&iId, "id", 0, "Instance ID")

	addInstanceCmd.Flags().StringVar(&iGroupId, "groupId", "", "Group ID of the team")
	addInstanceCmd.Flags().IntVar(&iToolId, "toolId", 0, "Tool ID to register")
	addInstanceCmd.Flags().IntVar(&iTemplateId, "templateId", 0, "Template ID used")

	instancesCmd.AddCommand(listInstancesCmd)
	instancesCmd.AddCommand(getInstanceCmd)
	instancesCmd.AddCommand(getInstanceMetadataCmd)
	instancesCmd.AddCommand(deleteInstanceCmd)
	instancesCmd.AddCommand(orphanedInstancesCmd)
	instancesCmd.AddCommand(outdatedInstancesCmd)
	instancesCmd.AddCommand(addInstanceCmd)

	TeamToolboxCmd.AddCommand(instancesCmd)
}
