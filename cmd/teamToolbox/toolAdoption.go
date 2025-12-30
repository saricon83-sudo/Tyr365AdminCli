package teamToolboxCmd

import (
	"fmt"

	teamToolBoxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

var toolAdoption = &cobra.Command{
	Use:   "toolAdoption",
	Short: "Get tool adoption rates",
	Long:  "Get the adoption rates of the tools available in the ToolBox",
	Run: func(cmd *cobra.Command, args []string) {

		client, err := teamToolBoxHelper.CreateAdminAPI()

		if err != nil {
			fmt.Println(err)
		}
		response, err := client.GetToolAdoption()

		if err != nil {
			fmt.Println(err)
		}

		teamToolBoxHelper.PrintToolAdoptionStatsTable(response)

	},
}

func init() {
	TeamToolboxCmd.AddCommand(toolAdoption)
}
