package teamToolboxCmd

import (
	"fmt"
	"os"

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
			fmt.Fprintf(os.Stderr, "Error creating admin API client: %v\n", err)
			return
		}

		response, err := client.GetRequestsByTool()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching requests by tool: %v\n", err)
			return
		}

		teamToolBoxHelper.PrintToolRequestCountTable(response)

	},
}

func init() {
	TeamToolboxCmd.AddCommand(getRequestPerTool)
}
