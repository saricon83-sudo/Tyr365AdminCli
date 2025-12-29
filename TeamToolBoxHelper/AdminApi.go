package teamToolboxHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// ============================================================================
// Admin API Client - All endpoints targeting /api/admin
// ============================================================================

// AdminAPI provides methods to interact with the Admin API endpoints
type AdminAPI struct {
	client *APIClient
}

// NewAdminAPI creates a new AdminAPI instance
func NewAdminAPI(client *APIClient) *AdminAPI {
	return &AdminAPI{client: client}
}

// CreateAdminAPI creates a new AdminAPI with default client
func CreateAdminAPI() (*AdminAPI, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}
	return NewAdminAPI(client), nil
}

// ============================================================================
// Helper Methods
// ============================================================================

func (api *AdminAPI) doGet(endpoint string, result interface{}) error {
	httpClient, err := api.client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("failed to get authenticated client: %w", err)
	}

	address := fmt.Sprintf("%s/admincli%s", api.client.BaseURL, endpoint)
	resp, err := httpClient.Get(address)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (api *AdminAPI) doPost(endpoint string, requestBody interface{}, result interface{}) error {
	httpClient, err := api.client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("failed to get authenticated client: %w", err)
	}

	address := fmt.Sprintf("%s/admincli%s", api.client.BaseURL, endpoint)

	var bodyReader io.Reader
	if requestBody != nil {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	resp, err := httpClient.Post(address, "application/json", bodyReader)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (api *AdminAPI) doPatch(endpoint string, requestBody interface{}, result interface{}) error {
	httpClient, err := api.client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("failed to get authenticated client: %w", err)
	}

	address := fmt.Sprintf("%s/admincli%s", api.client.BaseURL, endpoint)

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, address, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create PATCH request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make PATCH request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (api *AdminAPI) doDelete(endpoint string, result interface{}) error {
	httpClient, err := api.client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("failed to get authenticated client: %w", err)
	}

	address := fmt.Sprintf("%s/admincli%s", api.client.BaseURL, endpoint)

	req, err := http.NewRequest(http.MethodDelete, address, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// ============================================================================
// 1. Dashboard & Statistics
// ============================================================================

// GetDashboard retrieves comprehensive dashboard statistics
func (api *AdminAPI) GetDashboard() (*AdminDashboardStats, error) {
	var result AdminDashboardStats
	err := api.doGet("/dashboard", &result)
	return &result, err
}

// GetRequestsByStatus retrieves requests grouped by status
func (api *AdminAPI) GetRequestsByStatus() (map[string]int, error) {
	var result map[string]int
	err := api.doGet("/stats/requests-by-status", &result)
	return result, err
}

// GetRequestsByTool retrieves request counts per tool
func (api *AdminAPI) GetRequestsByTool() ([]ToolRequestCount, error) {
	var result []ToolRequestCount
	err := api.doGet("/stats/requests-by-tool", &result)
	return result, err
}

// GetRequestsByDay retrieves daily request counts for the last N days
func (api *AdminAPI) GetRequestsByDay(days int) ([]DailyRequestCount, error) {
	var result []DailyRequestCount
	endpoint := fmt.Sprintf("/stats/requests-by-day?days=%d", days)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetToolAdoption retrieves tool adoption statistics
func (api *AdminAPI) GetToolAdoption() ([]ToolAdoptionStats, error) {
	var result []ToolAdoptionStats
	err := api.doGet("/stats/tool-adoption", &result)
	return result, err
}

// GetStorageReleased retrieves storage released summary from clear site jobs
func (api *AdminAPI) GetStorageReleased() (*StorageReleasedSummary, error) {
	var result StorageReleasedSummary
	err := api.doGet("/stats/storage-released", &result)
	return &result, err
}

// GetArchiveJobStats retrieves archive job statistics
func (api *AdminAPI) GetArchiveJobStats() (*ArchiveJobStats, error) {
	var result ArchiveJobStats
	err := api.doGet("/stats/archive-jobs", &result)
	return &result, err
}

// GetPendingCounts retrieves pending counts across all job types
func (api *AdminAPI) GetPendingCounts() (*PendingCountsSummary, error) {
	var result PendingCountsSummary
	err := api.doGet("/stats/pending", &result)
	return &result, err
}

// ============================================================================
// 2. Error & Stuck Request Management
// ============================================================================

// GetErrors retrieves all error requests
func (api *AdminAPI) GetErrors(includeHidden bool) ([]ViewErrorRequest, error) {
	var result []ViewErrorRequest
	endpoint := fmt.Sprintf("/errors?includeHidden=%t", includeHidden)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetQueued retrieves queued requests
func (api *AdminAPI) GetQueued() ([]ViewQueuedRequest, error) {
	var result []ViewQueuedRequest
	err := api.doGet("/queued", &result)
	return result, err
}

// GetRunning retrieves running requests
func (api *AdminAPI) GetRunning() ([]ViewRunningRequest, error) {
	var result []ViewRunningRequest
	err := api.doGet("/running", &result)
	return result, err
}

// GetStuck retrieves stuck requests (running for longer than specified hours)
func (api *AdminAPI) GetStuck(hours int) ([]Request, error) {
	var result []Request
	endpoint := fmt.Sprintf("/stuck?hours=%d", hours)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetHighRetry retrieves requests with high retry counts
func (api *AdminAPI) GetHighRetry(minRetries int) ([]Request, error) {
	var result []Request
	endpoint := fmt.Sprintf("/high-retry?minRetries=%d", minRetries)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetSlowest retrieves slowest requests by duration
func (api *AdminAPI) GetSlowest(count int) ([]RequestDurationInfo, error) {
	var result []RequestDurationInfo
	endpoint := fmt.Sprintf("/slowest?count=%d", count)
	err := api.doGet(endpoint, &result)
	return result, err
}

// ============================================================================
// 3. Request Management
// ============================================================================

// GetRequest retrieves a specific request by ID with full details
func (api *AdminAPI) GetRequest(id int) (*Request, error) {
	var result Request
	endpoint := fmt.Sprintf("/requests/%d", id)
	err := api.doGet(endpoint, &result)
	return &result, err
}

// GetRequestsByGroup retrieves all requests for a specific group
func (api *AdminAPI) GetRequestsByGroup(groupId string) ([]Request, error) {
	var result []Request
	endpoint := fmt.Sprintf("/requests/by-group/%s", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetRequestsByInitiator retrieves requests by initiator email
func (api *AdminAPI) GetRequestsByInitiator(email string) ([]Request, error) {
	var result []Request
	endpoint := fmt.Sprintf("/requests/by-initiator/%s", url.PathEscape(email))
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetRequestsByEndpoint retrieves requests by endpoint
func (api *AdminAPI) GetRequestsByEndpoint(endpointName string) ([]Request, error) {
	var result []Request
	endpoint := fmt.Sprintf("/requests/by-endpoint/%s", url.PathEscape(endpointName))
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetRequestSteps retrieves request steps for a request
func (api *AdminAPI) GetRequestSteps(requestId int) ([]RequestStep, error) {
	var result []RequestStep
	endpoint := fmt.Sprintf("/requests/%d/steps", requestId)
	err := api.doGet(endpoint, &result)
	return result, err
}

// UpdateRequestStatus updates request status
func (api *AdminAPI) UpdateRequestStatus(id int, status string, message string) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/requests/%d/status", id)
	body := StatusUpdateRequest{Status: status, Message: message}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// UpdateRequestPriority updates request priority
func (api *AdminAPI) UpdateRequestPriority(id int, priority int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/requests/%d/priority", id)
	body := PriorityUpdateRequest{Priority: priority}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// RetryRequest retries a failed request
func (api *AdminAPI) RetryRequest(id int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/requests/%d/retry", id)
	err := api.doPost(endpoint, nil, &result)
	return &result, err
}

// UpdateRequestHidden hides/unhides a request
func (api *AdminAPI) UpdateRequestHidden(id int, hidden bool) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/requests/%d/hidden", id)
	body := HiddenUpdateRequest{Hidden: hidden}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// BulkRetryRequests bulk retries multiple requests
func (api *AdminAPI) BulkRetryRequests(ids []int) (*BulkOperationResult, error) {
	var result BulkOperationResult
	body := BulkRetryRequest{Ids: ids}
	err := api.doPost("/requests/bulk-retry", body, &result)
	return &result, err
}

// BulkHideRequests bulk hides/unhides multiple requests
func (api *AdminAPI) BulkHideRequests(ids []int, hidden bool) (*BulkOperationResult, error) {
	var result BulkOperationResult
	body := BulkHideRequest{Ids: ids, Hidden: hidden}
	err := api.doPost("/requests/bulk-hide", body, &result)
	return &result, err
}

// AddRequestStep adds a step to a request
func (api *AdminAPI) AddRequestStep(requestId int, step, status, message string) (*RequestStep, error) {
	var result RequestStep
	endpoint := fmt.Sprintf("/requests/%d/steps", requestId)
	body := AddStepRequest{Step: step, Status: status, Message: message}
	err := api.doPost(endpoint, body, &result)
	return &result, err
}

// ============================================================================
// 4. Tool Requests (tblToolRequest)
// ============================================================================

// GetToolRequests retrieves tool requests with optional filtering
func (api *AdminAPI) GetToolRequests(status, groupId string) ([]TblToolRequest, error) {
	var result []TblToolRequest
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if groupId != "" {
		params.Set("groupId", groupId)
	}
	endpoint := "/tool-requests"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetToolRequest retrieves a specific tool request
func (api *AdminAPI) GetToolRequest(id int) (*TblToolRequest, error) {
	var result TblToolRequest
	endpoint := fmt.Sprintf("/tool-requests/%d", id)
	err := api.doGet(endpoint, &result)
	return &result, err
}

// UpdateToolRequestStatus updates tool request status
func (api *AdminAPI) UpdateToolRequestStatus(id int, status string) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/tool-requests/%d/status", id)
	body := StatusUpdateRequest{Status: status}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// GetToolRequestCountByStatus retrieves tool request counts by status
func (api *AdminAPI) GetToolRequestCountByStatus() (map[string]int, error) {
	var result map[string]int
	err := api.doGet("/tool-requests/count-by-status", &result)
	return result, err
}

// ============================================================================
// 5. Tool Management
// ============================================================================

// GetTools retrieves all tools
func (api *AdminAPI) GetTools() ([]TblTool, error) {
	var result []TblTool
	err := api.doGet("/tools", &result)
	return result, err
}

// GetToolFullDetails retrieves tool with full details (rules, requirements, instances)
func (api *AdminAPI) GetToolFullDetails(id int) (*ToolFullDetails, error) {
	var result ToolFullDetails
	endpoint := fmt.Sprintf("/tools/%d", id)
	err := api.doGet(endpoint, &result)
	return &result, err
}

// UpdateToolEnabled enables/disables a tool
func (api *AdminAPI) UpdateToolEnabled(id int, enabled bool) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/tools/%d/enabled", id)
	body := EnabledUpdateRequest{Enabled: enabled}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// UpdateToolTopic updates tool topic name
func (api *AdminAPI) UpdateToolTopic(id int, topicName string) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/tools/%d/topic", id)
	body := TopicUpdateRequest{TopicName: topicName}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// UpdateToolTemplate updates tool template version
func (api *AdminAPI) UpdateToolTemplate(id int, templateId int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/tools/%d/template", id)
	body := TemplateUpdateRequest{TemplateId: templateId}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// UpdateTool updates tool properties
func (api *AdminAPI) UpdateTool(id int, update ToolUpdateDto) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/tools/%d", id)
	err := api.doPatch(endpoint, update, &result)
	return &result, err
}

// GetToolInstances retrieves instances of a specific tool
func (api *AdminAPI) GetToolInstances(toolId int) ([]TblToolInstance, error) {
	var result []TblToolInstance
	endpoint := fmt.Sprintf("/tools/%d/instances", toolId)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetToolRequestsForTool retrieves requests for a specific tool
func (api *AdminAPI) GetToolRequestsForTool(toolId int) ([]TblToolRequest, error) {
	var result []TblToolRequest
	endpoint := fmt.Sprintf("/tools/%d/requests", toolId)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetUnusedTools retrieves tools with no requests in the last N days
func (api *AdminAPI) GetUnusedTools(days int) ([]TblTool, error) {
	var result []TblTool
	endpoint := fmt.Sprintf("/tools/unused?days=%d", days)
	err := api.doGet(endpoint, &result)
	return result, err
}

// AddTool adds a new tool
func (api *AdminAPI) AddTool(tool TblTool) (*TblTool, error) {
	var result TblTool
	err := api.doPost("/tools", tool, &result)
	return &result, err
}

// ============================================================================
// 6. Tool Instances
// ============================================================================

// GetInstances retrieves all tool instances with optional filtering
func (api *AdminAPI) GetInstances(toolId int, groupId string) ([]TblToolInstance, error) {
	var result []TblToolInstance
	params := url.Values{}
	if toolId > 0 {
		params.Set("toolId", strconv.Itoa(toolId))
	}
	if groupId != "" {
		params.Set("groupId", groupId)
	}
	endpoint := "/instances"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetInstance retrieves a specific tool instance with metadata
func (api *AdminAPI) GetInstance(id int) (*TblToolInstance, error) {
	var result TblToolInstance
	endpoint := fmt.Sprintf("/instances/%d", id)
	err := api.doGet(endpoint, &result)
	return &result, err
}

// GetInstanceMetadata retrieves metadata for a tool instance
func (api *AdminAPI) GetInstanceMetadata(id int) ([]TblToolMetaDatum, error) {
	var result []TblToolMetaDatum
	endpoint := fmt.Sprintf("/instances/%d/metadata", id)
	err := api.doGet(endpoint, &result)
	return result, err
}

// DeleteInstance deletes a tool instance
func (api *AdminAPI) DeleteInstance(id int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/instances/%d", id)
	err := api.doDelete(endpoint, &result)
	return &result, err
}

// GetOrphanedInstances retrieves orphaned instances (group no longer exists)
func (api *AdminAPI) GetOrphanedInstances() ([]TblToolInstance, error) {
	var result []TblToolInstance
	err := api.doGet("/instances/orphaned", &result)
	return result, err
}

// GetOutdatedInstances retrieves instances with outdated template versions
func (api *AdminAPI) GetOutdatedInstances() ([]TblToolInstance, error) {
	var result []TblToolInstance
	err := api.doGet("/instances/outdated", &result)
	return result, err
}

// AddInstance adds a new tool instance
func (api *AdminAPI) AddInstance(instance TblToolInstance) (*TblToolInstance, error) {
	var result TblToolInstance
	err := api.doPost("/instances", instance, &result)
	return &result, err
}

// ============================================================================
// 7. Team/Group Management
// ============================================================================

// GetTeams retrieves managed teams with optional filtering
func (api *AdminAPI) GetTeams(status, origin string) ([]ManagedTeam, error) {
	var result []ManagedTeam
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if origin != "" {
		params.Set("origin", origin)
	}
	endpoint := "/teams"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetTeam retrieves a specific managed team
func (api *AdminAPI) GetTeam(groupId string) (*ManagedTeam, error) {
	var result ManagedTeam
	endpoint := fmt.Sprintf("/teams/%s", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return &result, err
}

// GetTeamFullDetails retrieves full team details (instances, requests, etc.)
func (api *AdminAPI) GetTeamFullDetails(groupId string) (*TeamFullDetails, error) {
	var result TeamFullDetails
	endpoint := fmt.Sprintf("/teams/%s/details", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return &result, err
}

// GetTeamInstances retrieves tool instances for a team
func (api *AdminAPI) GetTeamInstances(groupId string) ([]TblToolInstance, error) {
	var result []TblToolInstance
	endpoint := fmt.Sprintf("/teams/%s/instances", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetTeamRequests retrieves requests for a team
func (api *AdminAPI) GetTeamRequests(groupId string) ([]Request, error) {
	var result []Request
	endpoint := fmt.Sprintf("/teams/%s/requests", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return result, err
}

// UpdateTeamStatus updates managed team status
func (api *AdminAPI) UpdateTeamStatus(groupId string, status string) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/teams/%s/status", url.PathEscape(groupId))
	body := StatusUpdateRequest{Status: status}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// GetTeamsWithoutTools retrieves teams without any tool instances
func (api *AdminAPI) GetTeamsWithoutTools() ([]ManagedTeam, error) {
	var result []ManagedTeam
	err := api.doGet("/teams/without-tools", &result)
	return result, err
}

// SearchTeams searches teams by name, groupId, projectNo, etc.
func (api *AdminAPI) SearchTeams(query string) ([]TeamSearchResult, error) {
	var result []TeamSearchResult
	endpoint := fmt.Sprintf("/teams/search?q=%s", url.QueryEscape(query))
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetGroup retrieves a group record
func (api *AdminAPI) GetGroup(groupId string) (*TblGroup, error) {
	var result TblGroup
	endpoint := fmt.Sprintf("/groups/%s", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return &result, err
}

// GetGroupsView retrieves all groups from ViewGroups
func (api *AdminAPI) GetGroupsView() ([]ViewGroup, error) {
	var result []ViewGroup
	err := api.doGet("/groups/view", &result)
	return result, err
}

// ============================================================================
// 8. Archive & Export Jobs
// ============================================================================

// GetArchiveJobs retrieves archive jobs with optional status filter
func (api *AdminAPI) GetArchiveJobs(status string) ([]ArchiveJob, error) {
	var result []ArchiveJob
	endpoint := "/archive-jobs"
	if status != "" {
		endpoint += "?status=" + url.QueryEscape(status)
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetArchiveJob retrieves archive job with sub-jobs
func (api *AdminAPI) GetArchiveJob(id int) (*ArchiveJob, error) {
	var result ArchiveJob
	endpoint := fmt.Sprintf("/archive-jobs/%d", id)
	err := api.doGet(endpoint, &result)
	return &result, err
}

// UpdateArchiveJobStatus updates archive job status
func (api *AdminAPI) UpdateArchiveJobStatus(id int, status string) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/archive-jobs/%d/status", id)
	body := StatusUpdateRequest{Status: status}
	err := api.doPatch(endpoint, body, &result)
	return &result, err
}

// GetRequiredArchiveJobs retrieves required archive jobs
func (api *AdminAPI) GetRequiredArchiveJobs() ([]RequiredArchiveJob, error) {
	var result []RequiredArchiveJob
	err := api.doGet("/archive-jobs/required", &result)
	return result, err
}

// GetExportJobs retrieves export data jobs
func (api *AdminAPI) GetExportJobs(groupId string) ([]ExportDataJob, error) {
	var result []ExportDataJob
	endpoint := "/export-jobs"
	if groupId != "" {
		endpoint += "?groupId=" + url.QueryEscape(groupId)
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetClearSiteJobs retrieves clear site jobs
func (api *AdminAPI) GetClearSiteJobs(status string) ([]ClearSiteJob, error) {
	var result []ClearSiteJob
	endpoint := "/clear-site-jobs"
	if status != "" {
		endpoint += "?status=" + url.QueryEscape(status)
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetClearSiteJobsSummary retrieves clear site jobs summary
func (api *AdminAPI) GetClearSiteJobsSummary() (*ClearSiteJobsSummary, error) {
	var result ClearSiteJobsSummary
	err := api.doGet("/clear-site-jobs/summary", &result)
	return &result, err
}

// ============================================================================
// 9. Logging
// ============================================================================

// GetLogs retrieves logs with optional filtering
func (api *AdminAPI) GetLogs(subject, status, from, to string, limit int) ([]TblToolBoxLogger, error) {
	var result []TblToolBoxLogger
	params := url.Values{}
	if subject != "" {
		params.Set("subject", subject)
	}
	if status != "" {
		params.Set("status", status)
	}
	if from != "" {
		params.Set("from", from)
	}
	if to != "" {
		params.Set("to", to)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	endpoint := "/logs"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetRecentLogs retrieves recent logs
func (api *AdminAPI) GetRecentLogs(count int) ([]TblToolBoxLogger, error) {
	var result []TblToolBoxLogger
	endpoint := fmt.Sprintf("/logs/recent?count=%d", count)
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetLogsBySubject retrieves logs by subject
func (api *AdminAPI) GetLogsBySubject(subject string) ([]TblToolBoxLogger, error) {
	var result []TblToolBoxLogger
	endpoint := fmt.Sprintf("/logs/by-subject/%s", url.PathEscape(subject))
	err := api.doGet(endpoint, &result)
	return result, err
}

// AddLog adds a log entry
func (api *AdminAPI) AddLog(entry LogEntry) (*TblToolBoxLogger, error) {
	var result TblToolBoxLogger
	err := api.doPost("/logs", entry, &result)
	return &result, err
}

// DeleteLogsBefore deletes logs older than a date
func (api *AdminAPI) DeleteLogsBefore(date string) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/logs/before/%s", url.PathEscape(date))
	err := api.doDelete(endpoint, &result)
	return &result, err
}

// ============================================================================
// 10. Rules & Logic
// ============================================================================

// GetRules retrieves all rules
func (api *AdminAPI) GetRules() ([]TblToolRule, error) {
	var result []TblToolRule
	err := api.doGet("/rules", &result)
	return result, err
}

// GetRule retrieves a specific rule
func (api *AdminAPI) GetRule(id int) (*TblToolRule, error) {
	var result TblToolRule
	endpoint := fmt.Sprintf("/rules/%d", id)
	err := api.doGet(endpoint, &result)
	return &result, err
}

// AddRule adds a new rule
func (api *AdminAPI) AddRule(ruleName string) (*TblToolRule, error) {
	var result TblToolRule
	body := RuleCreate{RuleName: ruleName}
	err := api.doPost("/rules", body, &result)
	return &result, err
}

// DeleteRule deletes a rule
func (api *AdminAPI) DeleteRule(id int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/rules/%d", id)
	err := api.doDelete(endpoint, &result)
	return &result, err
}

// GetUnusedRules retrieves unused rules (not assigned to any tool)
func (api *AdminAPI) GetUnusedRules() ([]TblToolRule, error) {
	var result []TblToolRule
	err := api.doGet("/rules/unused", &result)
	return result, err
}

// GetRuleLogics retrieves all rule logics
func (api *AdminAPI) GetRuleLogics() ([]TblToolRuleLogic, error) {
	var result []TblToolRuleLogic
	err := api.doGet("/rule-logics", &result)
	return result, err
}

// GetRuleLogicsByTool retrieves rule logics for a specific tool
func (api *AdminAPI) GetRuleLogicsByTool(toolId int) ([]TblToolRuleLogic, error) {
	var result []TblToolRuleLogic
	endpoint := fmt.Sprintf("/rule-logics/by-tool/%d", toolId)
	err := api.doGet(endpoint, &result)
	return result, err
}

// AddRuleLogic adds a rule logic to a tool
func (api *AdminAPI) AddRuleLogic(toolId, ruleId int, logic, value string) (*TblToolRuleLogic, error) {
	var result TblToolRuleLogic
	body := RuleLogicCreate{ToolId: toolId, RuleId: ruleId, Logic: logic, Value: value}
	err := api.doPost("/rule-logics", body, &result)
	return &result, err
}

// DeleteRuleLogic deletes a rule logic
func (api *AdminAPI) DeleteRuleLogic(id int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/rule-logics/%d", id)
	err := api.doDelete(endpoint, &result)
	return &result, err
}

// ============================================================================
// 11. Extended Requirements
// ============================================================================

// GetExtendedRequirementsByTool retrieves extended requirements for a tool
func (api *AdminAPI) GetExtendedRequirementsByTool(toolId int) ([]TblToolExtendedRequirement, error) {
	var result []TblToolExtendedRequirement
	endpoint := fmt.Sprintf("/requirements/by-tool/%d", toolId)
	err := api.doGet(endpoint, &result)
	return result, err
}

// AddExtendedRequirement adds an extended requirement to a tool
func (api *AdminAPI) AddExtendedRequirement(toolId int, requirementName, requirementValue string) (*TblToolExtendedRequirement, error) {
	var result TblToolExtendedRequirement
	body := ExtendedRequirementCreate{ToolId: toolId, RequirementName: requirementName, RequirementValue: requirementValue}
	err := api.doPost("/requirements", body, &result)
	return &result, err
}

// DeleteExtendedRequirement deletes an extended requirement
func (api *AdminAPI) DeleteExtendedRequirement(id int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/requirements/%d", id)
	err := api.doDelete(endpoint, &result)
	return &result, err
}

// ============================================================================
// 12. Background Tasks
// ============================================================================

// GetTasks retrieves tasks with optional filtering
func (api *AdminAPI) GetTasks(status, jobType int) ([]TblTask, error) {
	var result []TblTask
	params := url.Values{}
	if status != -999 { // Using -999 as "not set" since -1 is a valid status
		params.Set("status", strconv.Itoa(status))
	}
	if jobType >= 0 {
		params.Set("jobType", strconv.Itoa(jobType))
	}
	endpoint := "/tasks"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	err := api.doGet(endpoint, &result)
	return result, err
}

// GetPendingTasks retrieves pending tasks
func (api *AdminAPI) GetPendingTasks() ([]TblTask, error) {
	var result []TblTask
	err := api.doGet("/tasks/pending", &result)
	return result, err
}

// GetFailedTasks retrieves failed tasks
func (api *AdminAPI) GetFailedTasks() ([]TblTask, error) {
	var result []TblTask
	err := api.doGet("/tasks/failed", &result)
	return result, err
}

// RetryTask retries a failed task
func (api *AdminAPI) RetryTask(id int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/tasks/%d/retry", id)
	err := api.doPost(endpoint, nil, &result)
	return &result, err
}

// GetTasksByProject retrieves tasks by project number
func (api *AdminAPI) GetTasksByProject(projectNo string) ([]TblTask, error) {
	var result []TblTask
	endpoint := fmt.Sprintf("/tasks/by-project/%s", url.PathEscape(projectNo))
	err := api.doGet(endpoint, &result)
	return result, err
}

// ============================================================================
// 13. GeoBIM
// ============================================================================

// GetGeoBIM retrieves GeoBIM tools view
func (api *AdminAPI) GetGeoBIM() ([]ViewToolGeoBim, error) {
	var result []ViewToolGeoBim
	err := api.doGet("/geobim", &result)
	return result, err
}

// GetGeoBIMByGroup retrieves GeoBIM by group ID
func (api *AdminAPI) GetGeoBIMByGroup(groupId string) (*ViewToolGeoBim, error) {
	var result ViewToolGeoBim
	endpoint := fmt.Sprintf("/geobim/%s", url.PathEscape(groupId))
	err := api.doGet(endpoint, &result)
	return &result, err
}

// ============================================================================
// 14. Health & Cleanup
// ============================================================================

// GetDatabaseHealth checks database health
func (api *AdminAPI) GetDatabaseHealth() (*HealthStatus, error) {
	var result HealthStatus
	err := api.doGet("/health/db", &result)
	return &result, err
}

// GetOrphanedRecords finds orphaned records across all tables
func (api *AdminAPI) GetOrphanedRecords() (*OrphanedRecordsSummary, error) {
	var result OrphanedRecordsSummary
	err := api.doGet("/cleanup/orphaned", &result)
	return &result, err
}

// DeleteOldRequests deletes old completed requests
func (api *AdminAPI) DeleteOldRequests(daysOld int) (*MessageResponse, error) {
	var result MessageResponse
	endpoint := fmt.Sprintf("/cleanup/old-requests?daysOld=%d", daysOld)
	err := api.doDelete(endpoint, &result)
	return &result, err
}

// ============================================================================
// 15. Bulk Operations
// ============================================================================

// BulkUpdateStatus bulk updates request statuses
func (api *AdminAPI) BulkUpdateStatus(ids []int, status string) (*BulkOperationResult, error) {
	var result BulkOperationResult
	body := BulkStatusUpdateRequest{Ids: ids, Status: status}
	err := api.doPost("/bulk/update-status", body, &result)
	return &result, err
}

// BulkRetryFailedForTool retries all failed requests for a specific tool
func (api *AdminAPI) BulkRetryFailedForTool(toolId int) (*BulkOperationResult, error) {
	var result BulkOperationResult
	endpoint := fmt.Sprintf("/bulk/retry-failed-for-tool/%d", toolId)
	err := api.doPost(endpoint, nil, &result)
	return &result, err
}

// BulkEnableTools bulk enables tools
func (api *AdminAPI) BulkEnableTools(ids []int) (*BulkOperationResult, error) {
	var result BulkOperationResult
	body := BulkToolIdsRequest{Ids: ids}
	err := api.doPost("/bulk/enable-tools", body, &result)
	return &result, err
}

// BulkDisableTools bulk disables tools
func (api *AdminAPI) BulkDisableTools(ids []int) (*BulkOperationResult, error) {
	var result BulkOperationResult
	body := BulkToolIdsRequest{Ids: ids}
	err := api.doPost("/bulk/disable-tools", body, &result)
	return &result, err
}
