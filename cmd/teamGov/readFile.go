/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package teamGov

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ApiResponse struct {
	StatusCode int `json:"statusCode"`
}

var fileName string

// readFileCmd represents the readFile command
var readFileCmd = &cobra.Command{
	Use:   "readFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var readGroups interface{}
		saveToFile.ReadDataFromJSONFile(fileName, &readGroups)

		if cmd.Flag("testPrint").Changed {
			outData, _ := json.Marshal(readGroups)
			fmt.Println(string(outData))
		}
		if cmd.Flag("updateCT").Changed {
			updateContenTypes(readGroups)
		}
		if cmd.Flag("getSharePointUrls").Changed {
			unifiedGroups, err := GetUnifiedGroupsFromRequests(readGroups)
			if err != nil {
				fmt.Print("FML")
			}
			saveToFile.SaveDataToJSONFile(unifiedGroups, "requestToUnifiedgroup.json")
		}
		if cmd.Flag("removeRetention").Changed {
			removeRetention(readGroups)
		}
		if cmd.Flag("removeRetentionValuesOnArchive").Changed {
			RemoveRetentionValuesOnArchive(readGroups)
		}
		if cmd.Flag("clearVersionHistoryOnSite").Changed {
			ClearVersionHistoryOnSite(readGroups)
		}
	},
}

func init() {
	readFileCmd.Flags().Bool("testPrint", false, "Ensure file can be read")
	readFileCmd.Flags().StringVarP(&fileName, "file", "f", "", "The name of the file you want to read from")
	readFileCmd.Flags().Bool("updateCT", false, "Call updateCTs in TeamGOv API")
	readFileCmd.Flags().Bool("getSharePointUrls", false, "Gets unifiedGroups model from requests objects")
	readFileCmd.Flags().Bool("clearVersionHistoryOnSite", false, "clearVersionHistoryOnSite")
	readFileCmd.Flags().Bool("removeRetentionValuesOnArchive", false, "Sets the retention Propertybag and extensionvalues to archived")
	readFileCmd.Flags().Bool("removeRetention", false, "Calls Exhange powershell (IPPS Session) to remove retention")

	TeamGovCmd.AddCommand(readFileCmd)
}

func removeRetention(readGroups interface{}) {
	logger := logging.GetLogger()
	outData, _ := json.Marshal(readGroups)
	groups, err := teamGovHttp.UnmarshalGroups(&outData)
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/AzureFunction/removeRetention",
			"method": "Post",
			"status": "Error",
		}).Error(err)
		return
	}

	// Semaphore channel to control the number of go routines
	semaphore := make(chan struct{}, 2)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failedGroups teamGovHttp.UnifiedGroupSlice

	for _, group := range groups {
		wg.Add(1)
		// Block if there are already two go routines running
		semaphore <- struct{}{}

		go func(group teamGovHttp.UnifiedGroup) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release a spot in the semaphore

			fmt.Printf("Clearing the following group of retention: %s\n", group.SharePointUrl)
			err := teamGovHttp.PostSharePointUrl(group.SharePointUrl)
			if err != nil {
				log.Printf("Error processing %s: %v\n", group.SharePointUrl, err)
				mu.Lock()
				failedGroups = append(failedGroups, group)
				mu.Unlock()
			} else {
				logger.WithFields(log.Fields{
					"url":    "/AzureFunction/removeRetention",
					"method": "Post",
					"status": "Success",
				}).Infof("Cleared site %s of retention", group.SharePointUrl)
			}
		}(group)
	}

	wg.Wait() // Wait for all go routines to finish

	if len(failedGroups) > 0 {
		saveToFile.SaveDataToJSONFile(failedGroups, "failedRetentionRemovals.json")
	}
}

func ClearVersionHistoryOnSite(readGroups interface{}) {
	logger := logging.GetLogger()
	outData, _ := json.Marshal(readGroups)
	requests, err := teamGovHttp.UnmarshalRequests(&outData)
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/api/teams/ReadFile/UpdateContentTypes",
			"method": "GET",
			"status": "Error",
		}).Error(err)
		return
	}

	// Create a channel to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 1)

	var wg sync.WaitGroup
	for _, group := range requests {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a token

		// Start a new goroutine for each group
		go func(group teamGovHttp.Request) {
			defer wg.Done()
			defer func() { <-semaphore }()

			queryParams := make(map[string]string)
			queryParams["groupId"] = group.GroupID
			_, err := teamGovHttp.Get("ClearVersionHistoryOnSite", queryParams)

			if err != nil {
				logger.WithFields(log.Fields{
					"url":         "/api/teams/ReadFile/ClearVersionHistoryOnSite",
					"method":      "GET",
					"status":      "Error",
					"queryParams": queryParams,
				}).Error(err)
				return
			}

			logger.WithFields(log.Fields{
				"url":         "/api/teams/ReadFile/ClearVersionHistoryOnSite",
				"method":      "GET",
				"status":      "Success",
				"queryParams": queryParams,
			}).Infof("Successfully processed group %s\n", group.GroupID)

		}(group)
	}

	// Wait for all goroutines to complete
	wg.Wait()

}

func updateContenTypes(readGroups interface{}) {
	logger := logging.GetLogger()
	outData, _ := json.Marshal(readGroups)
	requests, err := teamGovHttp.UnmarshalRequests(&outData)
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/api/teams/ReadFile/UpdateContentTypes",
			"method": "GET",
			"status": "Error",
		}).Error(err)
		return
	}

	// Create a channel to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 5)

	var wg sync.WaitGroup
	for _, group := range requests {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a token

		// Start a new goroutine for each group
		go func(group teamGovHttp.Request) {
			defer wg.Done()
			defer func() { <-semaphore }()

			queryParams := make(map[string]string)
			queryParams["groupId"] = group.GroupID
			_, err := teamGovHttp.Get("SetContentTypesToEditOnSite", queryParams)

			if err != nil {
				logger.WithFields(log.Fields{
					"url":         "/api/teams/ReadFile/SetContentTypesToEditOnSite",
					"method":      "GET",
					"status":      "Error",
					"queryParams": queryParams,
				}).Error(err)
				return
			}

			logger.WithFields(log.Fields{
				"url":         "/api/teams/ReadFile/SetContentTypesToEditOnSite",
				"method":      "GET",
				"status":      "Success",
				"queryParams": queryParams,
			}).Infof("Successfully processed group %s\n", group.GroupID)

		}(group)
	}

	// Wait for all goroutines to complete
	wg.Wait()

}

func RemoveRetentionValuesOnArchive(readGroups interface{}) {
	logger := logging.GetLogger()
	outData, _ := json.Marshal(readGroups)
	groups, err := teamGovHttp.UnmarshalGroups(&outData)
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/api/teams/ReadFile/UpdateContentTypes",
			"method": "GET",
			"status": "Error",
		}).Error(err)
		return
	}

	// Create a channel to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 3)

	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a token

		// Start a new goroutine for each group
		go func(group teamGovHttp.UnifiedGroup) {
			defer wg.Done()
			defer func() { <-semaphore }()

			queryParams := make(map[string]string)
			queryParams["groupId"] = group.GroupId
			queryParams["alias"] = group.Alias
			_, err := teamGovHttp.Get("RemoveRetentionValuesOnArchive", queryParams)

			if err != nil {
				logger.WithFields(log.Fields{
					"url":         "/api/teams/ReadFile/RemoveRetentionValuesOnArchive",
					"method":      "GET",
					"status":      "Error",
					"queryParams": queryParams,
				}).Error(err)
				return
			}

			logger.WithFields(log.Fields{
				"url":         "/api/teams/ReadFile/RemoveRetentionValuesOnArchive",
				"method":      "GET",
				"status":      "Success",
				"queryParams": queryParams,
			}).Infof("Successfully processed group %s\n", group.DisplayName)

		}(group)
	}

	// Wait for all goroutines to complete
	wg.Wait()

}

func GetUnifiedGroupsFromRequests(readGroups interface{}) (teamGovHttp.UnifiedGroupSlice, error) {
	logger := logging.GetLogger()

	// Marshal the readGroups into JSON
	outData, err := json.Marshal(readGroups)
	if err != nil {
		logger.Panic("Could not marshal read groups, exiting program")
	}

	// Unmarshal JSON into your request struct
	requests, err := teamGovHttp.UnmarshalRequests(&outData) // Ensure correct data handling
	if err != nil {
		logger.Panic("Could not unmarshal requests json, exiting program")
	}

	// Set up concurrency control and a slice to hold results
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var unifiedGroups teamGovHttp.UnifiedGroupSlice

	concurrency := 10
	sem := make(chan bool, concurrency)

	for _, request := range requests {
		wg.Add(1)
		sem <- true

		go func(request teamGovHttp.Request) {
			defer wg.Done()
			defer func() { <-sem }()

			// Fetch group data
			queryParams := map[string]string{"groupId": request.GroupID}
			group, err := teamGovHttp.Get("GetGroupDetails", queryParams)
			if err != nil {
				logger.WithFields(log.Fields{
					"url":          "/api/teams/ReadFile/GetGroup",
					"method":       "GET",
					"status":       "Error",
					"queryParams":  queryParams,
					"ErrorMessage": err,
				}).Warnf("Error with group with id: %s, skipping", request.GroupID)
				return
			}

			// Process group data
			m365Group, err := teamGovHttp.UnmarshalGroup(&group)
			if err != nil {
				logger.WithFields(log.Fields{
					"function":     "UnmarshalGroup",
					"status":       "Error",
					"ErrorMessage": err,
				}).Error("Failed to unmarshal group data")
				return
			}
			logger.WithFields(log.Fields{
				"url":         "/api/teams/ReadFile/GetGroup",
				"method":      "GET",
				"status":      "Success",
				"queryParams": queryParams,
				"group":       m365Group.DisplayName,
			}).Infof("Retrivered Unified group with id: %s", m365Group.GroupId)
			time.Sleep(5 * time.Second)

			// Aggregate the result into the slice
			mutex.Lock()
			unifiedGroups = append(unifiedGroups, m365Group)
			mutex.Unlock()

		}(request)
	}

	wg.Wait() // Wait for all goroutines to finish

	return unifiedGroups, nil // Return the collected results
}
