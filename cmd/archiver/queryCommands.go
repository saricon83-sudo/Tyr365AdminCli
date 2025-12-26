package archivercmd

import (
	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

var getArchiveJobsByStatusCmd = &cobra.Command{
	Use:   "getArchiveJobsByStatus",
	Short: "List archive jobs filtered by status",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			jobs, err := client.GetArchiveJobsByStatusTyped(status)
			if err != nil {
				return err
			}

			if len(jobs) == 0 {
				cmd.Printf("No archive jobs found with status %q.\n", status)
				return nil
			}

			jobs.PrintTable()
			return nil
		})
	},
}

var getArchiveJobByIDCmd = &cobra.Command{
	Use:   "getArchiveJobById",
	Short: "Fetch a single archive job by identifier",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			job, err := client.GetArchiveJobByIDTyped(id)
			if err != nil {
				return err
			}
			if job == nil {
				cmd.Printf("No archive job found with id %d.\n", id)
				return nil
			}

			jobSlice := archiver.ArchiveJobSlice{*job}
			jobSlice.PrintTable()
			return nil
		})
	},
}

var getArchiveJobDTOCmd = &cobra.Command{
	Use:   "getArchiveJobDto",
	Short: "Fetch archive job details by group identifier",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := cmd.Flags().GetString("groupId")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			job, err := client.GetArchiveJobDTOTyped(groupID)
			if err != nil {
				return err
			}
			if job == nil {
				cmd.Printf("No archive job DTO found for group %q.\n", groupID)
				return nil
			}

			jobSlice := archiver.ArchiveJobSlice{*job}
			jobSlice.PrintTable()
			return nil
		})
	},
}

func init() {
	getArchiveJobsByStatusCmd.Flags().String("status", "", "Status value to filter jobs (e.g. processing, completed)")
	_ = getArchiveJobsByStatusCmd.MarkFlagRequired("status")

	getArchiveJobByIDCmd.Flags().Int("id", 0, "Identifier of the archive job")
	_ = getArchiveJobByIDCmd.MarkFlagRequired("id")

	getArchiveJobDTOCmd.Flags().String("groupId", "", "Group identifier bound to the archive job")
	_ = getArchiveJobDTOCmd.MarkFlagRequired("groupId")

	ArchiverCmd.AddCommand(getArchiveJobsByStatusCmd)
	ArchiverCmd.AddCommand(getArchiveJobByIDCmd)
	ArchiverCmd.AddCommand(getArchiveJobDTOCmd)
}
