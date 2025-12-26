package teamToolboxCmd

import (
	"fmt"

	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

var getToolsByDay = &cobra.Command{
	Use:   "addToolToDb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// logger := logging.GetLogger()
		queryParams := make(map[string]interface{})
		if toolName != "" {
			queryParams["toolName"] = toolName
		}
		if currentTemplateId >= 0 {
			queryParams["currentTemlplateId"] = currentTemplateId
		}
		if topicName != "" {
			queryParams["topicName"] = topicName
		}

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			fmt.Println(err)
		}
		jsonBody, err := teamToolboxHelper.MarshalToJSON(queryParams)
		if err != nil {
			fmt.Println(err)
			return
		}
		response, err := client.PostWithJSONBody("addToolToDb", jsonBody)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(response))
	},
}

func init() {
	TeamToolboxCmd.AddCommand(addToolToDbCmd)

}
