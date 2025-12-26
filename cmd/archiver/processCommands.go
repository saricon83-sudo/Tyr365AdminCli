package archivercmd

import (
	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

var processArchiveSubJobsCmd = &cobra.Command{
	Use:   "processArchiveSubJobs",
	Short: "Process archive sub-jobs in the orchestrator",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.ProcessArchiveSubJobs()
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("ProcessArchiveSubJobs executed successfully.")
			}
			return nil
		})
	},
}

var processStatusOfSubJobsCmd = &cobra.Command{
	Use:   "processStatusOfSubJobs",
	Short: "Update the status of archive sub-jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.ProcessStatusOfSubJobs()
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("ProcessStatusOfSubJobs executed successfully.")
			}
			return nil
		})
	},
}

var reprocessSubJobsCmd = &cobra.Command{
	Use:   "reprocessSubJobs",
	Short: "Reprocess failed archive sub-jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.ReprocessSubJobs()
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("ReprocessSubJobs executed successfully.")
			}
			return nil
		})
	},
}

var finishExportCmd = &cobra.Command{
	Use:   "finishExport",
	Short: "Advance export jobs to the final stage",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.FinishExport()
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("FinishExport executed successfully.")
			}
			return nil
		})
	},
}

func init() {
	ArchiverCmd.AddCommand(processArchiveSubJobsCmd)
	ArchiverCmd.AddCommand(processStatusOfSubJobsCmd)
	ArchiverCmd.AddCommand(reprocessSubJobsCmd)
	ArchiverCmd.AddCommand(finishExportCmd)
}
