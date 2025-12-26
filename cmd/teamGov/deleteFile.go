/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteFileCmd represents the deleteFile command
var deleteFileCmd = &cobra.Command{
	Use:   "deleteFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		//if flag --file is used, delete the file
		if cmd.Flag("file").Changed {
			err := saveToFile.DeleteFile(fileName)
			if err != nil {
				logger.WithFields(log.Fields{
					"method": "Delete File",
					"status": "Error",
				}).Errorf("Error deleting file: %s\n", err)
				return
			}
			logger.WithFields(log.Fields{
				"method": "Delete File",
				"status": "Success",
			}).Info("File deleted")
		}
	},
}

func init() {
	TeamGovCmd.AddCommand(deleteFileCmd)
	deleteFileCmd.Flags().StringVarP(&fileName, "file", "f", "", "The name of the file you want to delete")

}
