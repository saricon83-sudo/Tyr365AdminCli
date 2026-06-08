package teamToolboxCmd

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
	teamToolboxHelper "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"
	"github.com/spf13/cobra"
)

var toolName string
var currentTemplateId int32
var topicName string

// addToolToDbCmd represents the addToolToDb command
var addToolToDbCmd = &cobra.Command{
	Use:   "addToolToDb",
	Short: "Add a new tool to the database",
	Long:  `Submits a request to the backend API to register a new tool with the specified name, template ID, and topic name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if toolName == "" {
			var err error
			toolName, err = pterm.DefaultInteractiveTextInput.Show("Enter Tool Name")
			if err != nil || toolName == "" {
				pterm.Error.Println("Tool Name is required")
				return
			}
		}

		if !cmd.Flags().Changed("templateId") {
			var err error
			val, err := pterm.DefaultInteractiveTextInput.Show("Enter Current Template ID")
			if err != nil || val == "" {
				pterm.Error.Println("Template ID is required")
				return
			}
			idVal, err := strconv.Atoi(val)
			if err != nil {
				pterm.Error.Printf("Invalid Template ID: %v\n", err)
				return
			}
			currentTemplateId = int32(idVal)
		}

		if topicName == "" {
			var err error
			topicName, err = pterm.DefaultInteractiveTextInput.Show("Enter Topic Name")
			if err != nil || topicName == "" {
				pterm.Error.Println("Topic Name is required")
				return
			}
		}

		queryParams := make(map[string]interface{})
		queryParams["toolName"] = toolName
		queryParams["currentTemplateId"] = currentTemplateId
		queryParams["topicName"] = topicName

		client, err := teamToolboxHelper.CreateClient()
		if err != nil {
			pterm.Error.Printf("Failed to create client: %v\n", err)
			return
		}

		jsonBody, err := teamToolboxHelper.MarshalToJSON(queryParams)
		if err != nil {
			pterm.Error.Printf("Failed to marshal request: %v\n", err)
			return
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Adding tool %s to database...", toolName))
		response, err := client.PostWithJSONBody("addToolToDb", jsonBody)
		if err != nil {
			spinner.Fail(err.Error())
			return
		}
		spinner.Success("Tool added successfully!")

		pterm.Info.Println("Response from server:")
		pterm.Println(string(response))
	},
}

func init() {
	addToolToDbCmd.Flags().StringVar(&toolName, "toolName", "", "Tool Name")
	addToolToDbCmd.Flags().Int32Var(&currentTemplateId, "templateId", 0, "Current Template ID")
	addToolToDbCmd.Flags().StringVar(&topicName, "topicName", "", "Topic Name")

	TeamToolboxCmd.AddCommand(addToolToDbCmd)
}
