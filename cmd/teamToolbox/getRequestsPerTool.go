package teamToolboxCmd

import (
	"fmt"
	teamToolBoxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

var getRequestPerTool = &cobra.Command{
	Use:   "getRequestPerTool",
	Short: "Gets requests for a specific tool",
	Long:  `This command queries and returns all entries of a specific tool from the DB`,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := teamToolBoxHelper.CreateAdminAPI()

		if err != nil {
			fmt.Println(err)
		}
		response, err := client.GetRequestsByTool()
		if err != nil {
			fmt.Println(err)
			return
		}
		teamToolBoxHelper.PrintToolRequestCountTable(response)

	},
}

func init() {
	TeamToolboxCmd.AddCommand(getRequestPerTool)
}
