package teamToolboxCmd

import (
	"fmt"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/output"
	"github.com/spf13/cobra"
)

var (
	geoGroupId string
)

// geobimCmd represents the geobim command group
var geobimCmd = &cobra.Command{
	Use:   "geobim",
	Short: "Inspect GeoBIM tools and group associations",
}

// listGeoBIMCmd represents the list command
var listGeoBIMCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all GeoBIM records",
	Run: func(cmd *cobra.Command, args []string) {
		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching GeoBIM records...")
		res, err := adminAPI.GetGeoBIM()
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("GeoBIM records loaded!")

		output.PrintResult(res, func() { printGeoBIMList(res) })
	},
}

// getGeoBIMCmd represents the get command
var getGeoBIMCmd = &cobra.Command{
	Use:   "get",
	Short: "Get GeoBIM records for a specific group ID",
	Run: func(cmd *cobra.Command, args []string) {
		if geoGroupId == "" {
			var err error
			geoGroupId, err = pterm.DefaultInteractiveTextInput.Show("Enter Group ID (GUID)")
			if err != nil || geoGroupId == "" {
				pterm.Error.Println("Group ID is required")
				return
			}
		}

		adminAPI, err := teamToolboxHelper.CreateAdminAPI()
		if err != nil {
			pterm.Error.Printf("Failed to connect to Admin API: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching GeoBIM record for group %s...", geoGroupId))
		res, err := adminAPI.GetGeoBIMByGroup(geoGroupId)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("GeoBIM record loaded!")

		output.PrintResult([]teamToolboxHelper.ViewToolGeoBim{*res}, func() { printGeoBIMList([]teamToolboxHelper.ViewToolGeoBim{*res}) })
	},
}

func printGeoBIMList(items []teamToolboxHelper.ViewToolGeoBim) {
	if len(items) == 0 {
		pterm.Info.Println("No GeoBIM records found matching criteria.")
		return
	}
	tableData := pterm.TableData{
		{"Group ID", "Team Name", "Project No", "Project Name", "GeoBIM Data"},
	}
	for _, item := range items {
		tableData = append(tableData, []string{
			item.GroupId,
			item.TeamName,
			item.ProjectNo,
			item.ProjectName,
			item.GeoBimData,
		})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func init() {
	getGeoBIMCmd.Flags().StringVar(&geoGroupId, "groupId", "", "Group ID GUID")

	geobimCmd.AddCommand(listGeoBIMCmd)
	geobimCmd.AddCommand(getGeoBIMCmd)

	TeamToolboxCmd.AddCommand(geobimCmd)
}
