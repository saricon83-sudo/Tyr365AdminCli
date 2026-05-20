/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package graphCommands

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// OwnedTeam represents a Team where the user is an owner
type OwnedTeam struct {
	GroupId     string `json:"groupId"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Mail        string `json:"mail"`
	Visibility  string `json:"visibility"`
	CreatedDate string `json:"createdDate"`
	IsArchived  bool   `json:"isArchived"`
}

type OwnedTeamSlice []OwnedTeam

var userUpn string

// getTeamsOwnedByUserCmd represents the command to get teams owned by a user
var getTeamsOwnedByUserCmd = &cobra.Command{
	Use:   "getTeamsOwnedByUser",
	Short: "Get all Teams where a user is an owner",
	Long: `Queries Microsoft Graph to find all Teams/Groups where the specified user is an owner.
	
Examples:
  365Admin graph getTeamsOwnedByUser --user john.doe@company.com --print
  365Admin graph getTeamsOwnedByUser --user john.doe@company.com --json
  365Admin graph getTeamsOwnedByUser --user john.doe@company.com --excel
  365Admin graph getTeamsOwnedByUser --user john.doe@company.com --csv`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()

		if userUpn == "" {
			fmt.Println("Error: --user flag is required")
			cmd.Help()
			return
		}

		// Get user ID from UPN
		userId, err := getUserIdFromUpn(userUpn)
		if err != nil {
			logger.WithFields(log.Fields{
				"user":   userUpn,
				"method": "getUserIdFromUpn",
				"status": "Error",
			}).Error(err)
			fmt.Printf("Error finding user %s: %v\n", userUpn, err)
			return
		}

		// Get owned objects (groups/teams)
		teams, err := getTeamsOwnedByUser(userId)
		if err != nil {
			logger.WithFields(log.Fields{
				"user":   userUpn,
				"userId": userId,
				"method": "getTeamsOwnedByUser",
				"status": "Error",
			}).Error(err)
			fmt.Printf("Error getting owned teams: %v\n", err)
			return
		}

		if len(teams) == 0 {
			fmt.Printf("User %s is not an owner of any Teams.\n", userUpn)
			return
		}

		logger.WithFields(log.Fields{
			"user":       userUpn,
			"teamsCount": len(teams),
			"method":     "getTeamsOwnedByUser",
			"status":     "Success",
		}).Info("Successfully retrieved owned teams")

		// Handle output formats
		if cmd.Flag("json").Changed {
			var fileName string
			fmt.Println("Enter a name for the JSON file (without extension):")
			fmt.Scanln(&fileName)
			if fileName == "" {
				fileName = fmt.Sprintf("teams_owned_by_%s", userUpn)
			}
			err := saveToFile.SaveDataToJSONFile(teams, fileName+".json")
			if err != nil {
				logger.WithFields(log.Fields{
					"method": "SaveDataToJSONFile",
					"status": "Error",
				}).Error(err)
				return
			}
		}

		if cmd.Flag("excel").Changed {
			var fileName string
			fmt.Println("Enter a name for the Excel file (without extension):")
			fmt.Scanln(&fileName)
			if fileName == "" {
				fileName = fmt.Sprintf("teams_owned_by_%s", userUpn)
			}
			err := saveToFile.SaveToExcel(teams, fileName)
			if err != nil {
				logger.WithFields(log.Fields{
					"method": "SaveToExcel",
					"status": "Error",
				}).Error(err)
				return
			}
		}

		if cmd.Flag("csv").Changed {
			var fileName string
			fmt.Println("Enter a name for the CSV file (without extension):")
			fmt.Scanln(&fileName)
			if fileName == "" {
				fileName = fmt.Sprintf("teams_owned_by_%s", userUpn)
			}
			err := saveToCSV(teams, fileName+".csv")
			if err != nil {
				logger.WithFields(log.Fields{
					"method": "saveToCSV",
					"status": "Error",
				}).Error(err)
				return
			}
		}

		if cmd.Flag("print").Changed {
			printOwnedTeamsTable(teams)
		}

		// Default: print summary if no output flag specified
		if !cmd.Flag("json").Changed && !cmd.Flag("excel").Changed && !cmd.Flag("csv").Changed && !cmd.Flag("print").Changed {
			fmt.Printf("\nFound %d Teams where %s is an owner:\n\n", len(teams), userUpn)
			printOwnedTeamsTable(teams)
		}
	},
}

func init() {
	getTeamsOwnedByUserCmd.Flags().StringVarP(&userUpn, "user", "u", "", "User's email/UPN (e.g., john.doe@company.com)")
	getTeamsOwnedByUserCmd.Flags().Bool("print", false, "Print results as a table")
	getTeamsOwnedByUserCmd.Flags().Bool("json", false, "Export results to JSON file")
	getTeamsOwnedByUserCmd.Flags().Bool("excel", false, "Export results to Excel file")
	getTeamsOwnedByUserCmd.Flags().Bool("csv", false, "Export results to CSV file")

	if err := getTeamsOwnedByUserCmd.MarkFlagRequired("user"); err != nil {
		fmt.Println(err)
	}

	GraphCmd.AddCommand(getTeamsOwnedByUserCmd)
}

// getUserIdFromUpn retrieves the user's object ID from their UPN
func getUserIdFromUpn(upn string) (string, error) {
	err := graphHelper.InitializeGraphForAppAuth()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Graph client: %w", err)
	}

	filter := fmt.Sprintf("userPrincipalName eq '%s'", upn)
	selectProps := []string{"id", "displayName", "userPrincipalName"}
	amount := int32(1)

	users, err := graphHelper.GetUsers(selectProps, &amount, filter)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	userList := users.GetValue()
	if len(userList) == 0 {
		return "", fmt.Errorf("user not found: %s", upn)
	}

	return *userList[0].GetId(), nil
}

// getTeamsOwnedByUser retrieves all Teams/Groups where the user is an owner
func getTeamsOwnedByUser(userId string) (OwnedTeamSlice, error) {
	err := graphHelper.InitializeGraphForAppAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Graph client: %w", err)
	}

	// Get owned objects using the new GraphHelper method
	ownedObjects, err := graphHelper.GetUserOwnedObjects(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned objects: %w", err)
	}

	var teams OwnedTeamSlice

	for _, obj := range ownedObjects {
		// Check if the object is a group (Teams are unified groups)
		odataType := obj.GetOdataType()
		if odataType != nil && *odataType == "#microsoft.graph.group" {
			// Get the group details
			groupId := obj.GetId()
			if groupId == nil {
				continue
			}

			// Fetch full group details
			group, err := graphHelper.GetGroupById(*groupId)
			if err != nil {
				continue // Skip if we can't get details
			}

			// Check if it's a unified group (Microsoft 365 Group / Team)
			groupTypes := group.GetGroupTypes()
			isUnifiedGroup := false
			for _, gt := range groupTypes {
				if gt == "Unified" {
					isUnifiedGroup = true
					break
				}
			}

			if !isUnifiedGroup {
				continue // Skip non-unified groups
			}

			// Build the OwnedTeam struct
			team := OwnedTeam{
				GroupId: *groupId,
			}

			if displayName := group.GetDisplayName(); displayName != nil {
				team.DisplayName = *displayName
			}
			if description := group.GetDescription(); description != nil {
				team.Description = *description
			}
			if mail := group.GetMail(); mail != nil {
				team.Mail = *mail
			}
			if visibility := group.GetVisibility(); visibility != nil {
				team.Visibility = *visibility
			}
			if createdDate := group.GetCreatedDateTime(); createdDate != nil {
				team.CreatedDate = createdDate.Format("2006-01-02 15:04:05")
			}

			// Check if Team is archived by trying to get Team details
			teamDetails, err := graphHelper.GetTeamById(*groupId)
			if err == nil && teamDetails != nil {
				if isArchived := teamDetails.GetIsArchived(); isArchived != nil {
					team.IsArchived = *isArchived
				}
			}

			teams = append(teams, team)
		}
	}

	return teams, nil
}

// printOwnedTeamsTable prints the teams in a formatted table
func printOwnedTeamsTable(teams OwnedTeamSlice) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Group ID", "Display Name", "Mail", "Visibility", "Created Date", "Archived"})
	table.SetAutoWrapText(false)
	table.SetRowLine(true)

	for _, team := range teams {
		archivedStr := "No"
		if team.IsArchived {
			archivedStr = "Yes"
		}
		row := []string{
			team.GroupId,
			team.DisplayName,
			team.Mail,
			team.Visibility,
			team.CreatedDate,
			archivedStr,
		}
		table.Append(row)
	}

	table.Render()
	fmt.Printf("\nTotal: %d Teams\n", len(teams))
}

// saveToCSV exports the teams to a CSV file
func saveToCSV(teams OwnedTeamSlice, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"GroupId", "DisplayName", "Description", "Mail", "Visibility", "CreatedDate", "IsArchived"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data rows
	for _, team := range teams {
		archivedStr := "false"
		if team.IsArchived {
			archivedStr = "true"
		}
		row := []string{
			team.GroupId,
			team.DisplayName,
			team.Description,
			team.Mail,
			team.Visibility,
			team.CreatedDate,
			archivedStr,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %w", err)
		}
	}

	fmt.Printf("CSV file created: %s\n", filename)
	return nil
}

// PrintTable implements the Printer interface for OwnedTeamSlice
func (teams *OwnedTeamSlice) PrintTable() {
	printOwnedTeamsTable(*teams)
}
