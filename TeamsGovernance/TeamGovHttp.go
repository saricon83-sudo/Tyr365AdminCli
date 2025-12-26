package teamGovHttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/saricon83-sudo/Tyr365AdminCli/internal/auth"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
	logging "github.com/saricon83-sudo/Tyr365AdminCli/logger"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// AuthGovernanceApi returns a token for the Teams Governance API.
// Deprecated: Use auth.GetGovernanceToken() directly instead.
func AuthGovernanceApi() (string, error) {
	return auth.GetGovernanceToken()
}

// AuthGraphApi returns a token for Microsoft Graph API.
// Deprecated: Use auth.GetGraphToken() directly instead.
func AuthGraphApi() (string, error) {
	return auth.GetGraphToken()
}

func RetrieveAuthToken() (string, error) {
	return "", errors.New("this function is deprecated")
}

func PrintToken() {
	fmt.Println("this function is deprecated")
}

func Get(endpoint string, queryParams ...map[string]string) ([]byte, error) {
	cfg := config.Get()

	token, err := AuthGovernanceApi()
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication token: %w", err)
	}
	if token == "" {
		return nil, errors.New("failed to get authentication token")
	}

	apiURL := cfg.GetString("resource") + "/api/teams/" + endpoint
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams[0] {
			query.Add(key, value)
		}
		apiURL += "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		trimmedBody := strings.TrimSpace(string(body))
		if trimmedBody == "" {
			return body, fmt.Errorf("unexpected response status: %s", resp.Status)
		}
		return body, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, trimmedBody)
	}

	return body, nil
}

func GetQuery(targetEndpoint string, queryParams map[string]string) ([]byte, error) {
	cfg := config.Get()

	// Retrieve the authentication token.
	token, err := AuthGovernanceApi()
	if err != nil || token == "" {
		return nil, fmt.Errorf("failed to get authentication token: %v", err)
	}

	// Construct the API URL.
	apiURL := cfg.GetString("resource") + "/api/teams/" + targetEndpoint
	query := url.Values{}

	// Iterate over the map of query parameters and add them to the query string.
	for key, value := range queryParams {
		query.Set(key, value)
	}

	// Append the query string to the API URL.
	if len(query) > 0 {
		apiURL += "?" + query.Encode()
	}
	// Create the HTTP GET request.
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read and return the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

func QueryManaged(groupId, teamName, status, origin, retention, fields string) ([]ManagedTeam, error) {
	cfg := config.Get()

	// Retrieve the authentication token
	token, err := AuthGovernanceApi()
	if err != nil || token == "" {
		return nil, fmt.Errorf("failed to get authentication token: %v", err)
	}

	// Construct the API URL and query parameters
	apiURL := cfg.GetString("resource") + "/api/teams/query"
	query := url.Values{}
	if groupId != "" {
		query.Set("groupId", groupId)
	}
	if teamName != "" {
		query.Set("teamName", teamName)
	}
	if status != "" {
		query.Set("status", status)
	}
	if origin != "" {
		query.Set("origin", origin)
	}
	if retention != "" {
		query.Set("retention", retention)
	}
	if fields != "" {
		query.Set("fields", fields)
	}

	// Append query string to the API URL
	if len(query) > 0 {
		apiURL += "?" + query.Encode()
	}

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read and unmarshal the response body
	var teams []ManagedTeam
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	err = json.Unmarshal(body, &teams)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return teams, nil
}

func Post(endpoint string, queryParams map[string]string) ([]byte, error) {
	cfg := config.Get()

	token, err := AuthGovernanceApi()
	if err != nil || token == "" {
		return nil, errors.New("failed to get authentication token")
	}

	apiURL := cfg.GetString("resource") + "/api/teams/" + endpoint
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(key, value)
		}
		apiURL += "?" + query.Encode()
	}

	req, err := http.NewRequest("POST", apiURL, nil) // No body is sent
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	// No need to set Content-Type for a request without a body

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return response, nil
}

// postSharePointUrl makes an HTTP POST request to a predefined URL with a JSON body containing the SharePoint URL.
func PostSharePointUrl(sharePointUrl string) error {
	// Define the URL to which the POST request will be sent.
	logger := logging.GetLogger()
	url := "https://github.com/saricon83-sudo/Tyr365AdminCli/security/secret-scanning/unblock-secret/2igdI04NqFYw9qo0dFVP8uiIGNj"

	// Create a map to hold the JSON payload.
	payload := map[string]string{
		"sharePointUrl": sharePointUrl,
	}

	// Marshal the map into a JSON object.
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/AzureFunction/removeRetention",
			"method": "Post",
			"status": "Error",
		}).Error(err)
		return err
	}

	// Create a new HTTP request with the JSON data.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/AzureFunction/removeRetention",
			"method": "Post",
			"status": "Error",
		}).Error(err)
		return err
	}

	// Set the Content-Type header to indicate JSON payload.
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{
			"url":    "/AzureFunction/removeRetention",
			"method": "Post",
			"status": "Error",
		}).Error(err)
		return err
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful.
	if resp.StatusCode != http.StatusOK {
		logger.WithFields(log.Fields{
			"url":    "/AzureFunction/removeRetention",
			"method": "Post",
			"status": "Error",
		}).Errorf("Received non-OK HTTP status: %s", resp.Status)

		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Optionally, handle the response data.
	return nil
}

func PostWithBody(endpoint string, queryParams map[string]string, body interface{}) ([]byte, error) {
	cfg := config.Get()

	token, err := AuthGovernanceApi()
	if err != nil || token == "" {
		return nil, errors.New("failed to get authentication token")
	}

	apiURL := cfg.GetString("resource") + "/api/teams/" + endpoint
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(key, value)
		}
		apiURL += "?" + query.Encode()
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return response, nil
}

func UnmarshalRequests(body *[]byte) (RequestSlice, error) {
	var requests RequestSlice
	err := json.Unmarshal(*body, &requests)
	if err == nil {
		return requests, nil
	}

	var request Request
	err = json.Unmarshal(*body, &request)
	if err == nil {
		return RequestSlice{request}, nil
	}

	return nil, fmt.Errorf("error unmarshalling to Request or []Request: %w", err)
}
func UnmarshalGroups(body *[]byte) (UnifiedGroupSlice, error) {
	var groups []UnifiedGroup
	err := json.Unmarshal(*body, &groups)
	if err == nil {
		return groups, nil
	}

	var group UnifiedGroup
	err = json.Unmarshal(*body, &group)
	if err == nil {
		return []UnifiedGroup{group}, nil
	}

	return nil, fmt.Errorf("error unmarshalling to UnifiedGroup or []UnifiedGroup: %w", err)
}

func UnmarshalGroup(body *[]byte) (UnifiedGroup, error) {
	var group UnifiedGroup
	err := json.Unmarshal(*body, &group)

	if err != nil {
		return group, err
	}
	return group, nil
}

func UnmarshalManagedTeams(body *[]byte) (ManagedTeamSlice, error) {
	var managedTeam []ManagedTeam
	err := json.Unmarshal(*body, &managedTeam)
	if err == nil {
		return managedTeam, nil
	}

	var request ManagedTeam
	err = json.Unmarshal(*body, &request)
	if err == nil {
		return []ManagedTeam{request}, nil
	}

	return nil, fmt.Errorf("error unmarshalling to Request or []Request: %w", err)
}
func UnmarshalInteger(body *[]byte) (int, error) {
	var value int
	err := json.Unmarshal(*body, &value)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling to integer: %w", err)
	}
	return value, nil
}

func PrintJSONResponseString(body *[]byte) error {
	var responseString string
	fmt.Println(*body)
	err := json.Unmarshal(*body, &responseString)
	if err != nil {
		return fmt.Errorf("error unmarshalling response body to string: %w", err)
	}
	fmt.Println(&responseString)
	return nil
}

func GetTaskETag(taskID string) (string, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/planner/tasks/%s/details", taskID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	accessToken, err := AuthGraphApi()
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	// ETag is found in the "ETag" response

	etag := resp.Header.Get("ETag")
	if etag == "" {
		return "", fmt.Errorf("ETag header not found in response")
	}
	return etag, nil
}

type RequestSlice []Request
type Printer interface {
	PrintTable()
}

func (r *RequestSlice) PrintTable() {

	// Create a table to display the response data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Created", "GroupID", "TeamName", "Endpoint", "CallerID", "Status", "ProvisioningStep", "Message", "InitiatedBy", "Modified", "RetryCount", "QueuePriority"}) // Customize the table header as needed
	for _, req := range *r {
		row := []string{
			fmt.Sprintf("%d", req.ID),
			req.Created,
			req.GroupID,
			req.TeamName,
			req.Endpoint,
			req.CallerID,
			req.Status,
			req.ProvisioningStep,
			req.Message,
			req.InitiatedBy,
			req.Modified,
			fmt.Sprintf("%v", req.RetryCount),
			fmt.Sprintf("%d", req.QueuePriority),
		}
		table.Append(row)
	}
	table.Render()
}
