package archiver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/auth"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
)

// AuthArchiverApi gets an authentication token for the Archiver API.
// Deprecated: Use auth.GetArchiverToken() directly instead.
func AuthArchiverApi() (string, error) {
	return auth.GetArchiverToken()
}

// CustomTime handles the API's custom datetime format
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements custom unmarshalling for the API's datetime format
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	str := string(data[1 : len(data)-1])

	// Try different formats that the API might use
	formats := []string{
		"2006-01-02T15:04:05.999", // With milliseconds
		"2006-01-02T15:04:05.99",  // With 2-digit milliseconds
		"2006-01-02T15:04:05.9",   // With 1-digit milliseconds
		"2006-01-02T15:04:05",     // Without milliseconds
		"2006-01-02T15:04:05Z",    // With Z timezone
		"2006-01-02",              // Date only (YYYY-MM-DD)
		time.RFC3339,              // Standard RFC3339
		time.RFC3339Nano,          // RFC3339 with nanoseconds
		time.DateOnly,             // Go 1.20+ date only format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			ct.Time = t
			return nil
		}
	}

	return fmt.Errorf("unable to parse time: %s", str)
}

// MarshalJSON implements custom marshalling
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format("2006-01-02T15:04:05.999") + `"`), nil
}

// Data models based on OpenAPI schemas

// ArchiveJob represents an archive job
type ArchiveJob struct {
	ID                         int             `json:"id"`
	IDExportDataJobs           int             `json:"idExportDataJobs"`
	FilePath                   string          `json:"filePath"`
	GroupID                    string          `json:"groupId"`
	Status                     *string         `json:"status"`
	SharePointURL              *string         `json:"sharePointUrl"`
	Created                    *CustomTime     `json:"created"`
	Alias                      string          `json:"alias"`
	ArchiveSubJobs             []ArchiveSubJob `json:"archiveSubJobs"`
	IDExportDataJobsNavigation interface{}     `json:"idExportDataJobsNavigation"`
}

// ArchiveSubJob represents a sub-job within an archive job
type ArchiveSubJob struct {
	ID                      int         `json:"id"`
	IDArchiveJobs           int         `json:"idArchiveJobs"`
	Type                    string      `json:"type"`
	Status                  *string     `json:"status"`
	GroupID                 *string     `json:"groupId"`
	Created                 *CustomTime `json:"created"`
	IDArchiveJobsNavigation *ArchiveJob `json:"idArchiveJobsNavigation"`
}

// ExportDataJob represents an export data job
type ExportDataJob struct {
	ID          int          `json:"id"`
	GroupID     string       `json:"groupId"`
	Alias       string       `json:"alias"`
	Status      string       `json:"status"`
	FilePath    string       `json:"filePath"`
	SiteURL     string       `json:"siteUrl"`
	RequestID   int          `json:"requestId"`
	ArchiveJobs []ArchiveJob `json:"archiveJobs"`
	Request     *Request     `json:"request"`
}

// Request represents a request
type Request struct {
	ID               int             `json:"id"`
	Created          CustomTime      `json:"created"`
	GroupID          *string         `json:"groupId"`
	TeamName         *string         `json:"teamName"`
	Endpoint         *string         `json:"endpoint"`
	CallerID         *string         `json:"callerId"`
	Parameters       *string         `json:"parameters"`
	Status           *string         `json:"status"`
	ProvisioningStep *string         `json:"provisioningStep"`
	Message          *string         `json:"message"`
	InitiatedBy      *string         `json:"initiatedBy"`
	Modified         *CustomTime     `json:"modified"`
	RowVersion       []byte          `json:"rowVersion"`
	ClientTaskID     *int            `json:"clientTaskId"`
	LTPMessageSent   *bool           `json:"ltpmessageSent"`
	Hidden           *bool           `json:"hidden"`
	GroupInformation *string         `json:"groupInformation"`
	RetryCount       *int            `json:"retryCount"`
	QueuePriority    *int            `json:"queuePriority"`
	ExportDataJobs   []ExportDataJob `json:"exportDataJobs"`
	RequestSteps     []RequestStep   `json:"requestSteps"`
}

// RequestStep represents a step in a request
type RequestStep struct {
	ID        int         `json:"id"`
	RequestID int         `json:"requestId"`
	Step      *string     `json:"step"`
	Status    *string     `json:"status"`
	Message   *string     `json:"message"`
	Modified  *CustomTime `json:"modified"`
	Request   *Request    `json:"request"`
}

// Slice types for table printing
type ArchiveJobSlice []ArchiveJob
type ArchiveSubJobSlice []ArchiveSubJob
type ExportDataJobSlice []ExportDataJob
type RequestSlice []Request
type RequestStepSlice []RequestStep

// Printer interface for table printing functionality
type Printer interface {
	PrintTable()
}

// ArchiverClient provides methods to interact with the M365 Archiver API
type ArchiverClient struct {
	baseURL string
}

// NewArchiverClient creates a new archiver client
func NewArchiverClient() (*ArchiverClient, error) {
	cfg := config.Get()

	baseURL := cfg.GetString("archiverAddress")
	if baseURL == "" {
		return nil, errors.New("archiverAddress not found in configuration")
	}
	fmt.Printf("Debug: Archiver base URL: %s\n", baseURL) // Debug print
	return &ArchiverClient{
		baseURL: baseURL,
	}, nil
}

// makeAuthenticatedRequest makes an HTTP request with authentication
func (c *ArchiverClient) makeAuthenticatedRequest(method, endpoint string, body io.Reader, queryParams map[string]string) ([]byte, error) {
	// Get authentication token
	token, err := AuthArchiverApi()
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication token: %w", err)
	}

	// Construct URL
	apiURL := c.baseURL + endpoint
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Set(key, value)
		}
		apiURL += "?" + query.Encode()
	}

	// Debug: Print the URL being called
	fmt.Printf("Debug: Making %s request to: %s\n", method, apiURL)

	// Create request
	req, err := http.NewRequest(method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read response body for error details
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(errorBody))
	}

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return responseBody, nil
}

// GET endpoints

// CreateArchiveJobs calls GET /api/Archiver/CreateArchiveJobs
func (c *ArchiverClient) CreateArchiveJobs() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/CreateArchiveJobs", nil, nil)
}

// CreateArchiveSubJobs calls GET /api/Archiver/CreateArchiveSubJobs
func (c *ArchiverClient) CreateArchiveSubJobs() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/CreateArchiveSubJobs", nil, nil)
}

// GetJobs calls GET /api/Archiver/GetJobs
func (c *ArchiverClient) GetJobs() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/GetJobs", nil, nil)
}

// ProcessArchiveSubJobs calls GET /api/Archiver/ProcessArchiveSubJobs
func (c *ArchiverClient) ProcessArchiveSubJobs() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/ProcessArchiveSubJobs", nil, nil)
}

// ProcessStatusOfSubJobs calls GET /api/Archiver/ProcessStatusOfSubJobs
func (c *ArchiverClient) ProcessStatusOfSubJobs() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/ProcessStatusOfSubJobs", nil, nil)
}

// GetArchiveJobsByStatus calls GET /api/Archiver/GetArchiveJobsByStatus
func (c *ArchiverClient) GetArchiveJobsByStatus(status string) ([]byte, error) {
	params := map[string]string{"status": status}
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/GetArchiveJobsByStatus", nil, params)
}

// GetArchiveSubJobsByStatus calls GET /api/Archiver/GetArchiveSubJobsByStatus
func (c *ArchiverClient) GetArchiveSubJobsByStatus(status string) ([]byte, error) {
	params := map[string]string{"status": status}
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/GetArchiveSubJobsByStatus", nil, params)
}

// GetArchiveJobByID calls GET /api/Archiver/GetArchiveJobById
func (c *ArchiverClient) GetArchiveJobByID(id int) ([]byte, error) {
	params := map[string]string{"id": strconv.Itoa(id)}
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/GetArchiveJobById", nil, params)
}

// ReprocessSubJobs calls GET /api/Archiver/ReprocessSubJobs
func (c *ArchiverClient) ReprocessSubJobs() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/ReprocessSubJobs", nil, nil)
}

// GetArchiveJobDTO calls GET /api/Archiver/GetArchiveJobDTO
func (c *ArchiverClient) GetArchiveJobDTO(groupID string) ([]byte, error) {
	params := map[string]string{"groupId": groupID}
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/GetArchiveJobDTO", nil, params)
}

// FinishExport calls GET /api/Archiver/FinishExport
func (c *ArchiverClient) FinishExport() ([]byte, error) {
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/FinishExport", nil, nil)
}

// Test calls GET /api/Archiver/Test
func (c *ArchiverClient) Test(id int, groupID string) ([]byte, error) {
	params := map[string]string{
		"id":      strconv.Itoa(id),
		"groupId": groupID,
	}
	return c.makeAuthenticatedRequest("GET", "/api/Archiver/Test", nil, params)
}

// POST endpoints

// CreateArchiveJobbet calls POST /api/Archiver/CreateArchiveJobbet
func (c *ArchiverClient) CreateArchiveJobbet(archiveJob ArchiveJob) ([]byte, error) {
	jsonBody, err := json.Marshal(archiveJob)
	if err != nil {
		return nil, fmt.Errorf("error marshalling archive job: %w", err)
	}

	return c.makeAuthenticatedRequest("POST", "/api/Archiver/CreateArchiveJobbet", bytes.NewBuffer(jsonBody), nil)
}

// UpdateArchiveJob calls POST /api/Archiver/UpdateArchiveJob
func (c *ArchiverClient) UpdateArchiveJob(id int, status string) ([]byte, error) {
	params := map[string]string{
		"id":     strconv.Itoa(id),
		"status": status,
	}
	return c.makeAuthenticatedRequest("POST", "/api/Archiver/UpdateArchiveJob", nil, params)
}

// UpdateArchiveSubJob calls POST /api/Archiver/UpdateArchiveSubJob
func (c *ArchiverClient) UpdateArchiveSubJob(id int, status string) ([]byte, error) {
	params := map[string]string{
		"id":     strconv.Itoa(id),
		"status": status,
	}
	return c.makeAuthenticatedRequest("POST", "/api/Archiver/UpdateArchiveSubJob", nil, params)
}

// DELETE endpoints

// DeleteArchiveJob calls DELETE /api/Archiver/DeleteArchiveJob
func (c *ArchiverClient) DeleteArchiveJob(id int) ([]byte, error) {
	params := map[string]string{"id": strconv.Itoa(id)}
	return c.makeAuthenticatedRequest("DELETE", "/api/Archiver/DeleteArchiveJob", nil, params)
}

// DeleteArchiveSubJob calls DELETE /api/Archiver/DeleteArchiveSubJob
func (c *ArchiverClient) DeleteArchiveSubJob(id int) ([]byte, error) {
	params := map[string]string{"id": strconv.Itoa(id)}
	return c.makeAuthenticatedRequest("DELETE", "/api/Archiver/DeleteArchiveSubJob", nil, params)
}

// Typed wrapper methods that return structured data instead of []byte
// These methods combine API calls with unmarshalling for convenience

// GetJobsTyped returns all jobs as a typed slice ready for table printing
func (c *ArchiverClient) GetJobsTyped() (ArchiveJobSlice, error) {
	body, err := c.GetJobs()
	if err != nil {
		return nil, err
	}
	jobs, err := UnmarshalArchiveJobs(body)
	if err != nil {
		return nil, err
	}
	return ArchiveJobSlice(jobs), nil
}

// GetArchiveJobByIDTyped returns a single job as a typed struct ready for table printing
func (c *ArchiverClient) GetArchiveJobByIDTyped(id int) (*ArchiveJob, error) {
	body, err := c.GetArchiveJobByID(id)
	if err != nil {
		return nil, err
	}
	return UnmarshalArchiveJob(body)
}

// GetArchiveJobsByStatusTyped returns jobs by status as a typed slice ready for table printing
func (c *ArchiverClient) GetArchiveJobsByStatusTyped(status string) (ArchiveJobSlice, error) {
	body, err := c.GetArchiveJobsByStatus(status)
	if err != nil {
		return nil, err
	}
	jobs, err := UnmarshalArchiveJobs(body)
	if err != nil {
		return nil, err
	}
	return ArchiveJobSlice(jobs), nil
}

// GetArchiveSubJobsByStatusTyped returns sub jobs by status as a typed slice ready for table printing
func (c *ArchiverClient) GetArchiveSubJobsByStatusTyped(status string) (ArchiveSubJobSlice, error) {
	body, err := c.GetArchiveSubJobsByStatus(status)
	if err != nil {
		return nil, err
	}
	subJobs, err := UnmarshalArchiveSubJobs(body)
	if err != nil {
		return nil, err
	}
	return ArchiveSubJobSlice(subJobs), nil
}

// GetArchiveJobDTOTyped returns archive job DTO as a typed struct
func (c *ArchiverClient) GetArchiveJobDTOTyped(groupID string) (*ArchiveJob, error) {
	body, err := c.GetArchiveJobDTO(groupID)
	if err != nil {
		return nil, err
	}
	return UnmarshalArchiveJob(body)
}

// Helper functions for unmarshalling responses

// UnmarshalArchiveJob unmarshals a single archive job from response body
func UnmarshalArchiveJob(body []byte) (*ArchiveJob, error) {
	var job ArchiveJob
	err := json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling archive job: %w", err)
	}
	return &job, nil
}

// UnmarshalArchiveJobs unmarshals multiple archive jobs from response body
func UnmarshalArchiveJobs(body []byte) ([]ArchiveJob, error) {
	var jobs []ArchiveJob
	err := json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling archive jobs: %w", err)
	}
	return jobs, nil
}

// UnmarshalArchiveSubJob unmarshals a single archive sub job from response body
func UnmarshalArchiveSubJob(body []byte) (*ArchiveSubJob, error) {
	var subJob ArchiveSubJob
	err := json.Unmarshal(body, &subJob)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling archive sub job: %w", err)
	}
	return &subJob, nil
}

// UnmarshalArchiveSubJobs unmarshals multiple archive sub jobs from response body
func UnmarshalArchiveSubJobs(body []byte) ([]ArchiveSubJob, error) {
	var subJobs []ArchiveSubJob
	err := json.Unmarshal(body, &subJobs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling archive sub jobs: %w", err)
	}
	return subJobs, nil
}

// UnmarshalExportDataJob unmarshals a single export data job from response body
func UnmarshalExportDataJob(body []byte) (*ExportDataJob, error) {
	var job ExportDataJob
	err := json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling export data job: %w", err)
	}
	return &job, nil
}

// UnmarshalExportDataJobs unmarshals multiple export data jobs from response body
func UnmarshalExportDataJobs(body []byte) ([]ExportDataJob, error) {
	var jobs []ExportDataJob
	err := json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling export data jobs: %w", err)
	}
	return jobs, nil
}

// UnmarshalRequest unmarshals a single request from response body
func UnmarshalRequest(body []byte) (*Request, error) {
	var request Request
	err := json.Unmarshal(body, &request)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling request: %w", err)
	}
	return &request, nil
}

// UnmarshalRequests unmarshals multiple requests from response body
func UnmarshalRequests(body []byte) ([]Request, error) {
	var requests []Request
	err := json.Unmarshal(body, &requests)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling requests: %w", err)
	}
	return requests, nil
}

// UnmarshalRequestStep unmarshals a single request step from response body
func UnmarshalRequestStep(body []byte) (*RequestStep, error) {
	var step RequestStep
	err := json.Unmarshal(body, &step)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling request step: %w", err)
	}
	return &step, nil
}

// UnmarshalRequestSteps unmarshals multiple request steps from response body
func UnmarshalRequestSteps(body []byte) ([]RequestStep, error) {
	var steps []RequestStep
	err := json.Unmarshal(body, &steps)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling request steps: %w", err)
	}
	return steps, nil
}

// UnmarshalString unmarshals a simple string response
func UnmarshalString(body []byte) (string, error) {
	var result string
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling string response: %w", err)
	}
	return result, nil
}

// UnmarshalInteger unmarshals a simple integer response
func UnmarshalInteger(body []byte) (int, error) {
	var result int
	err := json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling integer response: %w", err)
	}
	return result, nil
}

// PrintTable methods for each slice type

// PrintTable displays ArchiveJob slice data in a formatted table
func (a *ArchiveJobSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "GroupID", "Status", "SharePointURL", "Created", "Alias", "FilePath"})

	for _, job := range *a {
		// Handle pointer fields safely
		status := ""
		if job.Status != nil {
			status = *job.Status
		}

		sharePointURL := ""
		if job.SharePointURL != nil {
			sharePointURL = *job.SharePointURL
		}

		created := ""
		if job.Created != nil {
			created = job.Created.Time.Format("2006-01-02 15:04:05")
		}

		row := []string{
			fmt.Sprintf("%d", job.ID),
			job.GroupID,
			status,
			sharePointURL,
			created,
			job.Alias,
			job.FilePath,
		}
		table.Append(row)
	}

	table.Render()
}

// PrintTable displays ArchiveSubJob slice data in a formatted table
func (a *ArchiveSubJobSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "IDArchiveJobs", "Type", "Status", "GroupID", "Created"})

	for _, subJob := range *a {
		// Handle pointer fields safely
		status := ""
		if subJob.Status != nil {
			status = *subJob.Status
		}

		groupID := ""
		if subJob.GroupID != nil {
			groupID = *subJob.GroupID
		}

		created := ""
		if subJob.Created != nil {
			created = subJob.Created.Time.Format("2006-01-02 15:04:05")
		}

		row := []string{
			fmt.Sprintf("%d", subJob.ID),
			fmt.Sprintf("%d", subJob.IDArchiveJobs),
			subJob.Type,
			status,
			groupID,
			created,
		}
		table.Append(row)
	}

	table.Render()
}

// PrintTable displays ExportDataJob slice data in a formatted table
func (e *ExportDataJobSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "GroupID", "Alias", "Status", "FilePath", "SiteURL", "RequestID"})

	for _, job := range *e {
		row := []string{
			fmt.Sprintf("%d", job.ID),
			job.GroupID,
			job.Alias,
			job.Status,
			job.FilePath,
			job.SiteURL,
			fmt.Sprintf("%d", job.RequestID),
		}
		table.Append(row)
	}

	table.Render()
}

// PrintTable displays Request slice data in a formatted table
func (r *RequestSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Created", "GroupID", "TeamName", "Status", "Message", "InitiatedBy", "RetryCount"})

	for _, req := range *r {
		// Handle pointer fields safely
		groupID := ""
		if req.GroupID != nil {
			groupID = *req.GroupID
		}

		teamName := ""
		if req.TeamName != nil {
			teamName = *req.TeamName
		}

		status := ""
		if req.Status != nil {
			status = *req.Status
		}

		message := ""
		if req.Message != nil {
			message = *req.Message
		}

		initiatedBy := ""
		if req.InitiatedBy != nil {
			initiatedBy = *req.InitiatedBy
		}

		retryCount := ""
		if req.RetryCount != nil {
			retryCount = fmt.Sprintf("%d", *req.RetryCount)
		}

		row := []string{
			fmt.Sprintf("%d", req.ID),
			req.Created.Time.Format("2006-01-02 15:04:05"),
			groupID,
			teamName,
			status,
			message,
			initiatedBy,
			retryCount,
		}
		table.Append(row)
	}

	table.Render()
}

// PrintTable displays RequestStep slice data in a formatted table
func (r *RequestStepSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "RequestID", "Step", "Status", "Message", "Modified"})

	for _, step := range *r {
		// Handle pointer fields safely
		stepName := ""
		if step.Step != nil {
			stepName = *step.Step
		}

		status := ""
		if step.Status != nil {
			status = *step.Status
		}

		message := ""
		if step.Message != nil {
			message = *step.Message
		}

		modified := ""
		if step.Modified != nil {
			modified = step.Modified.Time.Format("2006-01-02 15:04:05")
		}

		row := []string{
			fmt.Sprintf("%d", step.ID),
			fmt.Sprintf("%d", step.RequestID),
			stepName,
			status,
			message,
			modified,
		}
		table.Append(row)
	}

	table.Render()
}
