/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package archivercmd

import (
	"fmt"

	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

// getJobsCmd represents the getJobs command
var GetJobsCmd = &cobra.Command{
	Use:   "getJobs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getJobs called")
		client, err := archiver.NewArchiverClient()
		if err != nil {
			fmt.Println("Error creating Archiver client:", err)
			return
		}

		// Example 1: Get single job using the new typed method
		fmt.Println("=== Getting single job by ID (650) ===")
		job, err := client.GetArchiveJobByIDTyped(650)
		if err != nil {
			fmt.Println("Error fetching job:", err)
			return
		}

		if job == nil {
			fmt.Println("No job found with ID 650")
			return
		}

		fmt.Printf("Job details: %+v\n", job)

		// Use the TablePrinter functionality for a single job
		jobSlice := archiver.ArchiveJobSlice{*job}
		jobSlice.PrintTable()

		// Example 2: Get all jobs using the new typed method
		fmt.Println("\n=== Getting all jobs ===")
		allJobs, err := client.GetJobsTyped()
		if err != nil {
			fmt.Println("Error fetching all jobs:", err)
			return
		}

		fmt.Printf("Found %d jobs\n", len(allJobs))
		if len(allJobs) > 0 {
			// Print first 5 jobs as table
			if len(allJobs) > 5 {
				allJobs = allJobs[:5]
				fmt.Println("Showing first 5 jobs:")
			}
			allJobs.PrintTable()
		}

		// Example 3: Get jobs by status using the new typed method
		fmt.Println("\n=== Getting jobs by status (completed) ===")
		completedJobs, err := client.GetArchiveJobsByStatusTyped("completed")
		if err != nil {
			fmt.Println("Error fetching completed jobs:", err)
			// Don't return, just log and continue
		} else {
			fmt.Printf("Found %d completed jobs\n", len(completedJobs))
			if len(completedJobs) > 0 && len(completedJobs) <= 10 {
				completedJobs.PrintTable()
			} else if len(completedJobs) > 10 {
				fmt.Println("Too many results to display in table (>10)")
			}
		}
	},
}

func init() {
	ArchiverCmd.AddCommand(GetJobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getJobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getJobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
