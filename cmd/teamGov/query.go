package teamGov

import (
	"encoding/json"
	"fmt"
	"sync"

	saveToFile "github.com/saricon83-sudo/Tyr365AdminCli/SaveToFile"
	teamGovHttp "github.com/saricon83-sudo/Tyr365AdminCli/TeamsGovernance"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Assuming these variables are declared at the package level to store flag values
var (
	endpoint        string
	created         string
	createdEnd      string
	callerId        string
	initiatedByUser string
	top             int // Assuming there's a sensible default or 0 indicates "use default"
	templateID      int
)

// newCmd represents the command for the new endpoint
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Querys the governance API for requests",
	Long: `Querys the governance API for requests based on the provided parameters.
    The results can be printed as a table or saved to an Excel file.
    If no parameters are provided, the command will return the 50 latest Create requests
    
    Available parameters:
    --endpoint: Comma-separated endpoints (e.g. "endpoint1,endpoint2") . Endpoints are "Create", "ApplySPTemplate", "ApplyTeamTemplate", "Group", "ArchiveTeam". If no endpoint is provided, default endpoint is Create
    --created: Start date (YYYY/MM/DD) (e.g. "2021/01/01"). If no date is provided, default date is 60 days ago.
    --createdEnd: End date (YYYY/MM/DD) (e.g. "2021/01/01"). If no date is provided, default date is today.
    --callerId: Comma-separated caller IDs (e.g. "callerId1,callerId2"). Default callerId is "Tyra".
    --initiatedBy: User who initiated the request (e.g. "user1@tyrens.se"). If no user is provided, default user is "sposervice@tyrens.onmicrosoft.com".
    --status: Comma-separated statuses (e.g. "status1,status2"). Default status is "Succeeded". Available statuses are "Succeeded", "Error", "Queued", "Processing".
    --top: Limit the number of results. Default is 50. Max is 1000.
    --templateID: Template ID to filter the requests -- Assuming there's a sensible default or 0 indicates "use default". This command should not be used with --print or --excel
        .`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		queryParams := make(map[string]string)
		if endpoint != "" {
			queryParams["endpoint"] = endpoint
		}
		if created != "" {
			queryParams["created"] = created
		}
		if createdEnd != "" {
			queryParams["createdEnd"] = createdEnd
		}
		if callerId != "" {
			queryParams["callerId"] = callerId
		}
		if initiatedByUser != "" {
			queryParams["initiatedByUser"] = initiatedByUser
		}
		if status != "" {
			queryParams["status"] = status
		}
		if top > 0 { // Assuming a non-zero value should be included
			queryParams["top"] = fmt.Sprintf("%d", top)
		}

		body, err := teamGovHttp.GetQuery("CliQuery", queryParams)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CliQuery",
				"method": "GET",
				"status": "Error",
				"query":  queryParams,
			}).Error(err)
			return
		}
		requests, err := teamGovHttp.UnmarshalRequests(&body)

		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CliQuery",
				"method": "GET",
				"status": "Error",
			}).Error(err)
			return
		}

		if cmd.Flag("templateID").Changed {
			if templateID == 0 {
				fmt.Println("Please provide a template ID")
				return
			}
			RunGORutine(requests, templateID)
		}

		if cmd.Flag("excel").Changed {
			savedToFile(&requests)
		}
		if cmd.Flag("json").Changed {
			var fileName string
			fmt.Println("Enter a name for the JSON file (without extension):")
			fmt.Scanln(&fileName)

			err := saveToFile.SaveDataToJSONFile(&requests, fileName+".json")
			if err != nil {
				logger.WithFields(log.Fields{
					"url":    "/api/teams/CliQuery",
					"method": "GET",
					"status": "Error",
					"query":  queryParams,
				}).Error(err)
				return
			}
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CliQuery",
				"method": "GET",
				"status": "Succeeded",
				"query":  queryParams,
			}).Info("Data successfully saved to JSON file:", fileName+".json")

		}
		if cmd.Flag("print").Changed {
			ViewTable(&requests)

		}
		if err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CliQuery",
				"method": "GET",
				"status": "Error",
				"query":  queryParams,
			}).Error(err)

			return
		}

	},
}

func init() {
	// Register flags for the newCmd
	queryCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "Comma-separated endpoints")
	queryCmd.Flags().StringVarP(&created, "created", "c", "", "Start date (YYYY/MM/DD)")
	queryCmd.Flags().StringVarP(&createdEnd, "createdEnd", "C", "", "End date (YYYY/MM/DD)")
	queryCmd.Flags().StringVarP(&callerId, "callerId", "i", "", "Comma-separated caller IDs")
	queryCmd.Flags().StringVarP(&initiatedByUser, "initiatedBy", "u", "", "User who initiated")
	queryCmd.Flags().StringVarP(&status, "status", "s", "", "Comma-separated statuses")
	queryCmd.Flags().IntVarP(&top, "top", "t", 0, "Limit the number of results")
	queryCmd.Flags().Bool("help", false, "Print help for the command")
	queryCmd.Flags().Bool("excel", false, "Save the response to an Excel file")
	queryCmd.Flags().Bool("json", false, "Save the response to a JSON file")
	queryCmd.Flags().Bool("print", false, "Print the response as a table")
	// queryCmd.Flags().StringVarP(&jqQuery, "jq", "j", "", "jq query to filter the JSON output")
	queryCmd.Flags().IntVarP(&templateID, "templateID", "T", 0, "Template ID to filter the requests")
	TeamGovCmd.AddCommand(queryCmd)
}

func RunGORutine(requests []teamGovHttp.Request, templateID int) {
	// Use goroutines to process requests concurrently
	// WaitGroup is like an advnaced clock that waits for all goroutines to finish
	var wg sync.WaitGroup
	// Assuming a buffer size, which might be tuned based on your application's requirements
	bufferSize := 100

	requestsChan := make(chan teamGovHttp.Request, bufferSize) // Buffered channel for requests
	resultsChan := make(chan teamGovHttp.Request, bufferSize)  // Buffered channel for results

	// Start worker goroutines
	numWorkers := 20 // This might be adjusted based on your system's capabilities
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, templateID, requestsChan, resultsChan)
	}

	// Collector goroutine to gather results
	go func() {
		wg.Wait()
		close(resultsChan) // Safely close results channel once all workers are done
	}()

	// Send requests to the workers
	for _, req := range requests {
		requestsChan <- req // This is safe as long as the total number of requests doesn't exceed the channel's capacity significantly
	}
	close(requestsChan) // Signal workers that no more requests are coming

	// Collect and process matching requests
	var matchedRequests teamGovHttp.RequestSlice
	for req := range resultsChan {
		matchedRequests = append(matchedRequests, req)
	}

	ViewTable(&matchedRequests)
}

func worker(wg *sync.WaitGroup, templateID int, requestsChan <-chan teamGovHttp.Request, resultsChan chan<- teamGovHttp.Request) {
	logger := logging.GetLogger()
	defer wg.Done()
	for req := range requestsChan {
		var params teamGovHttp.Parameters
		if err := json.Unmarshal([]byte(req.Parameters), &params); err != nil {
			logger.WithFields(log.Fields{
				"url":    "/api/teams/CliQuery",
				"method": "GET",
				"status": "Error",
			}).Errorf("Error unmarshaling Parameters for request ID %d: %v\n", req.ID, err)
			continue
		}

		if params.TemplateId == templateID { // Your filter criteria
			resultsChan <- req // Send matching requests to the results channel
		}
	}
}

func savedToFile(requests *teamGovHttp.RequestSlice) {
	var fileName string
	fmt.Println("Name your new excel file:")
	fmt.Scanln(&fileName)
	saveToFile.SaveToExcel(requests, fileName)
}
