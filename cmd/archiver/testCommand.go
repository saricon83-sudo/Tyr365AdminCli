package archivercmd

import (
	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Invoke the Archiver test endpoint with a job id and group id",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		groupID, err := cmd.Flags().GetString("groupId")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.Test(id, groupID)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("Test endpoint executed successfully.")
			}
			return nil
		})
	},
}

func init() {
	testCmd.Flags().Int("id", 0, "Identifier of the archive job to test")
	testCmd.Flags().String("groupId", "", "Group identifier to scope the test call")
	_ = testCmd.MarkFlagRequired("id")
	_ = testCmd.MarkFlagRequired("groupId")

	ArchiverCmd.AddCommand(testCmd)
}
