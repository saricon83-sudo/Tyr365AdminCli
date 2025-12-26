/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// postToTeamsCmd represents the postToTeams command
var postToTeamsCmd = &cobra.Command{
	Use:   "postToTeams",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		var webhookURL = `https://lindendevelop.webhook.office.com/webhookb2/26aa87f9-9e27-4017-9ab4-fa925d1bce3e@fd6c9945-d5d1-423c-acf9-2a5431314398/IncomingWebhook/e9abe5737f5c42a39f22d96401bcf802/527d6a2e-8ddd-4d17-993d-99aad9a74283`

		if err := postMessageToTeamsWebhook(webhookURL, "Hello from the CLI!"); err != nil {
			logger.WithFields(log.Fields{
				"url":    webhookURL,
				"method": "POST",
				"status": "Error",
			}).Error(err)
		}
	},
}

func init() {
	TeamGovCmd.AddCommand(postToTeamsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postToTeamsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postToTeamsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func postMessageToTeamsWebhook(webhookURL, message string) error {
	fmt.Println("Posting message to Teams channel...")
	payload := map[string]interface{}{
		"text": message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Teams webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	fmt.Println("Message successfully posted to Teams channel.")
	return nil
}
