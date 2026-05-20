/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package graphCommands

import (
	"fmt"
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	graphhelper "github.com/saricon83-sudo/Tyr365AdminCli/graphHelper"
	"github.com/spf13/cobra"
)

var propertiesFlag []string
var amountFlag int32
var filterFlag string

var getusersCmd = &cobra.Command{
	Use:   "getusers",
	Short: "Retrieves users from Microsoft Graph",
	Long:  `Retrieves a list of users from Microsoft Graph, with customizable properties via the command line.`,
	Run: func(cmd *cobra.Command, args []string) {

		users, err := graphHelper.GetUsers(propertiesFlag, &amountFlag, filterFlag)
		if err != nil {
			fmt.Printf("Error getting users: %v\n", err)
			return
		}
		userse, err := JsonifyUserResponse(users)

		if cmd.Flag("json").Changed {
			var fileName string
			fmt.Println("Enter a name for the JSON file (without extension):")
			fmt.Scanln(&fileName)

			err := saveToFile.SaveDataToJSONFile(userse, fileName+".json")
			if err != nil {
				fmt.Printf("Error saving data to JSON file: %s\n", err)
				return
			}
			fmt.Println("Data successfully saved to JSON file:", fileName+".json")
		}

	},
}

func init() {
	GraphCmd.AddCommand(getusersCmd)
	getusersCmd.Flags().StringSliceVarP(&propertiesFlag, "properties", "p", []string{"displayName", "id", "mail"}, "Properties to select (comma-separated)")
	getusersCmd.Flags().Int32VarP(&amountFlag, "amount", "a", 25, "Amount of users to retrieve")
	getusersCmd.Flags().StringVarP(&filterFlag, "filter", "f", "", "Filter the users by a specific property")
	getusersCmd.Flags().BoolP("json", "j", false, "Save the output to a JSON file")
}

func safeGetString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func safeGetStringSlice(ss []string) []string {
	if ss == nil {
		return []string{}
	}
	return ss
}
func safeFormatTime(t *time.Time) string {
	if t != nil {
		return t.Format(time.RFC3339) // ISO 8601 format
	}
	return ""
}

// Helper function to safely get a bool from a *bool, defaulting to false if nil
func safeGetBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false // default value if nil
}
func stringPointer(s string) *string {
	return &s
}

func JsonifyUserResponse(users models.UserCollectionResponseable) ([]graphhelper.User, error) {
	var mappedUsers []graphhelper.User

	for _, u := range users.GetValue() {
		mappedUser := graphhelper.User{
			ID:                          safeGetString(u.GetId()),
			DeletedDateTime:             safeFormatTime(u.GetDeletedDateTime()),
			AccountEnabled:              safeGetBool(u.GetAccountEnabled()),
			AgeGroup:                    safeGetString(u.GetAgeGroup()),
			BusinessPhones:              safeGetStringSlice(u.GetBusinessPhones()),
			City:                        safeGetString(u.GetCity()),
			CreatedDateTime:             safeFormatTime(u.GetCreatedDateTime()),
			CreationType:                stringPointer(safeGetString(u.GetCreationType())),
			CompanyName:                 safeGetString(u.GetCompanyName()),
			ConsentProvidedForMinor:     stringPointer(safeGetString(u.GetConsentProvidedForMinor())),
			Country:                     safeGetString(u.GetCountry()),
			Department:                  safeGetString(u.GetDepartment()),
			DisplayName:                 safeGetString(u.GetDisplayName()),
			EmployeeId:                  stringPointer(safeGetString(u.GetEmployeeId())),
			EmployeeHireDate:            stringPointer(safeFormatTime(u.GetEmployeeHireDate())),
			EmployeeLeaveDateTime:       stringPointer(safeFormatTime(u.GetEmployeeLeaveDateTime())),
			EmployeeType:                stringPointer(safeGetString(u.GetEmployeeType())),
			FaxNumber:                   stringPointer(safeGetString(u.GetFaxNumber())),
			GivenName:                   safeGetString(u.GetGivenName()),
			ImAddresses:                 safeGetStringSlice(u.GetImAddresses()),
			JobTitle:                    safeGetString(u.GetJobTitle()),
			LegalAgeGroupClassification: safeGetString(u.GetLegalAgeGroupClassification()),
			Mail:                        safeGetString(u.GetMail()),
			MailNickname:                safeGetString(u.GetMailNickname()),
			MobilePhone:                 safeGetString(u.GetMobilePhone()),
			OfficeLocation:              safeGetString(u.GetOfficeLocation()),
			PostalCode:                  safeGetString(u.GetPostalCode()),
			PreferredLanguage:           safeGetString(u.GetPreferredLanguage()),
			ProxyAddresses:              safeGetStringSlice(u.GetProxyAddresses()),
			SecurityIdentifier:          safeGetString(u.GetSecurityIdentifier()),
			State:                       safeGetString(u.GetState()),
			StreetAddress:               safeGetString(u.GetStreetAddress()),
			Surname:                     safeGetString(u.GetSurname()),
			UsageLocation:               safeGetString(u.GetUsageLocation()),
			UserPrincipalName:           safeGetString(u.GetUserPrincipalName()),
			UserType:                    safeGetString(u.GetUserType()),
		}
		mappedUsers = append(mappedUsers, mappedUser)
	}

	return mappedUsers, nil
}
