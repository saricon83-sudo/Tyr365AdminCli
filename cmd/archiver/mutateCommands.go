package archivercmd

import (
	"encoding/json"
	"fmt"
	"os"

	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

var createArchiveJobbetCmd = &cobra.Command{
	Use:   "createArchiveJobbet",
	Short: "Create a new archive job from a JSON payload",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		payload, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed reading payload file: %w", err)
		}

		var job archiver.ArchiveJob
		if err := json.Unmarshal(payload, &job); err != nil {
			return fmt.Errorf("invalid archive job payload: %w", err)
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.CreateArchiveJobbet(job)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("CreateArchiveJobbet executed successfully.")
			}
			return nil
		})
	},
}

var updateArchiveJobCmd = &cobra.Command{
	Use:   "updateArchiveJob",
	Short: "Update the status of an archive job",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.UpdateArchiveJob(id, status)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("UpdateArchiveJob executed successfully.")
			}
			return nil
		})
	},
}

var updateArchiveSubJobCmd = &cobra.Command{
	Use:   "updateArchiveSubJob",
	Short: "Update the status of an archive sub-job",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.UpdateArchiveSubJob(id, status)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("UpdateArchiveSubJob executed successfully.")
			}
			return nil
		})
	},
}

var deleteArchiveJobCmd = &cobra.Command{
	Use:   "deleteArchiveJob",
	Short: "Delete an archive job by identifier",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.DeleteArchiveJob(id)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("DeleteArchiveJob executed successfully.")
			}
			return nil
		})
	},
}

var deleteArchiveSubJobCmd = &cobra.Command{
	Use:   "deleteArchiveSubJob",
	Short: "Delete an archive sub-job by identifier",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		return withArchiverClient(cmd, func(client *archiver.ArchiverClient) error {
			body, err := client.DeleteArchiveSubJob(id)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				cmd.Println(string(body))
			} else {
				cmd.Println("DeleteArchiveSubJob executed successfully.")
			}
			return nil
		})
	},
}

func init() {
	createArchiveJobbetCmd.Flags().String("input", "", "Path to JSON file describing the archive job payload")
	_ = createArchiveJobbetCmd.MarkFlagRequired("input")

	updateArchiveJobCmd.Flags().Int("id", 0, "Identifier of the archive job to update")
	updateArchiveJobCmd.Flags().String("status", "", "New status value")
	_ = updateArchiveJobCmd.MarkFlagRequired("id")
	_ = updateArchiveJobCmd.MarkFlagRequired("status")

	updateArchiveSubJobCmd.Flags().Int("id", 0, "Identifier of the archive sub-job to update")
	updateArchiveSubJobCmd.Flags().String("status", "", "New status value")
	_ = updateArchiveSubJobCmd.MarkFlagRequired("id")
	_ = updateArchiveSubJobCmd.MarkFlagRequired("status")

	deleteArchiveJobCmd.Flags().Int("id", 0, "Identifier of the archive job to delete")
	_ = deleteArchiveJobCmd.MarkFlagRequired("id")

	deleteArchiveSubJobCmd.Flags().Int("id", 0, "Identifier of the archive sub-job to delete")
	_ = deleteArchiveSubJobCmd.MarkFlagRequired("id")

	ArchiverCmd.AddCommand(createArchiveJobbetCmd)
	ArchiverCmd.AddCommand(updateArchiveJobCmd)
	ArchiverCmd.AddCommand(updateArchiveSubJobCmd)
	ArchiverCmd.AddCommand(deleteArchiveJobCmd)
	ArchiverCmd.AddCommand(deleteArchiveSubJobCmd)
}
