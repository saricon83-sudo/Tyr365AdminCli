package teamToolboxCmd

import (
	"fmt"

	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

var includeHidden bool

var getErrors = &cobra.Command{
	Use:   "getErrors",
	Short: "Get requests with status errors",
	Long:  "Get requests from the Governance API that currently has the error status",
	Run: func(cmd *cobra.Command, args []string) {

		client, err := teamToolboxHelper.CreateAdminAPI()

		if err != nil {
			fmt.Println(err)
		}
		response, err := client.GetErrors(includeHidden)

		if err != nil {
			fmt.Println(err)
		}

		teamToolboxHelper.PrintViewErrorRequestTable(response)
	},
}

func init() {
	TeamToolboxCmd.AddCommand(getErrors)
}
