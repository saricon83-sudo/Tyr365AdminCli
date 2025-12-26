package archivercmd

import (
	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

var createArchiveJobsCmd = &cobra.Command{
	Use:   "createArchiveJobs",
	Short: "Trigger creation of archive jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.CreateArchiveJobs()
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("CreateArchiveJobs executed successfully.")
			}
			return nil
		})
	},
}

var createArchiveSubJobsCmd = &cobra.Command{
	Use:   "createArchiveSubJobs",
	Short: "Trigger creation of archive sub-jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.CreateArchiveSubJobs()
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("CreateArchiveSubJobs executed successfully.")
			}
			return nil
		})
	},
}

func init() {
	ArchiverCmd.AddCommand(createArchiveJobsCmd)
	ArchiverCmd.AddCommand(createArchiveSubJobsCmd)
}
