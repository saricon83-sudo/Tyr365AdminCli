package teamToolboxHelper

import "time"

// ============================================================================
// Dashboard & Statistics Models
// ============================================================================

// AdminDashboardStats represents the comprehensive dashboard statistics
type AdminDashboardStats struct {
	TotalTools           int            `json:"totalTools"`
	EnabledTools         int            `json:"enabledTools"`
	DisabledTools        int            `json:"disabledTools"`
	TotalToolInstances   int            `json:"totalToolInstances"`
	TotalManagedTeams    int            `json:"totalManagedTeams"`
	TotalRequests        int            `json:"totalRequests"`
	TotalToolRequests    int            `json:"totalToolRequests"`
	PendingArchiveJobs   int            `json:"pendingArchiveJobs"`
	ErrorRequests        int            `json:"errorRequests"`
	StuckRequests        int            `json:"stuckRequests"`
	RequestsByStatus     map[string]int `json:"requestsByStatus"`
	ToolRequestsByStatus map[string]int `json:"toolRequestsByStatus"`
}

// ToolRequestCount represents request counts per tool
type ToolRequestCount struct {
	ToolId         int    `json:"toolId"`
	ToolName       string `json:"toolName"`
	RequestCount   int    `json:"requestCount"`
	CompletedCount int    `json:"completedCount"`
	ErrorCount     int    `json:"errorCount"`
}

// DailyRequestCount represents daily request counts
type DailyRequestCount struct {
	Date      string `json:"date"`
	Count     int    `json:"count"`
	Completed int    `json:"completed"`
	Errors    int    `json:"errors"`
}

// ToolAdoptionStats represents tool adoption statistics
type ToolAdoptionStats struct {
	ToolId        int       `json:"toolId"`
	ToolName      string    `json:"toolName"`
	InstanceCount int       `json:"instanceCount"`
	UniqueTeams   int       `json:"uniqueTeams"`
	FirstAdoption time.Time `json:"firstAdoption"`
	LastAdoption  time.Time `json:"lastAdoption"`
}

// StorageReleasedPeriod represents storage released in a period
type StorageReleasedPeriod struct {
	Grain           string `json:"grain"`
	StorageReleased int64  `json:"storageReleased"`
	Period          string `json:"period"`
}

// StorageReleasedSummary represents storage released summary
type StorageReleasedSummary struct {
	TotalStorageReleased int64                   `json:"totalStorageReleased"`
	TotalFilesDeleted    int                     `json:"totalFilesDeleted"`
	TotalJobs            int                     `json:"totalJobs"`
	CompletedJobs        int                     `json:"completedJobs"`
	ByPeriod             []StorageReleasedPeriod `json:"byPeriod"`
}

// ArchiveJobStats represents archive job statistics
type ArchiveJobStats struct {
	Total            int `json:"total"`
	Pending          int `json:"pending"`
	Running          int `json:"running"`
	Completed        int `json:"completed"`
	Failed           int `json:"failed"`
	SubJobsTotal     int `json:"subJobsTotal"`
	SubJobsPending   int `json:"subJobsPending"`
	SubJobsCompleted int `json:"subJobsCompleted"`
}

// PendingCountsSummary represents pending counts across job types
type PendingCountsSummary struct {
	PendingRequests      int `json:"pendingRequests"`
	RunningRequests      int `json:"runningRequests"`
	PendingToolRequests  int `json:"pendingToolRequests"`
	PendingArchiveJobs   int `json:"pendingArchiveJobs"`
	PendingExportJobs    int `json:"pendingExportJobs"`
	PendingClearSiteJobs int `json:"pendingClearSiteJobs"`
	PendingTasks         int `json:"pendingTasks"`
}

// ============================================================================
// Request Models
// ============================================================================

// Request represents a request record
type Request struct {
	Id             int             `json:"id"`
	GroupId        string          `json:"groupId"`
	Endpoint       string          `json:"endpoint"`
	Status         string          `json:"status"`
	Priority       int             `json:"priority"`
	RetryCount     int             `json:"retryCount"`
	Hidden         bool            `json:"hidden"`
	Message        string          `json:"message"`
	InitiatedBy    string          `json:"initiatedBy"`
	Created        time.Time       `json:"created"`
	Modified       time.Time       `json:"modified"`
	RequestSteps   []RequestStep   `json:"requestSteps,omitempty"`
	ExportDataJobs []ExportDataJob `json:"exportDataJobs,omitempty"`
}

// RequestStep represents a step in a request
type RequestStep struct {
	Id        int       `json:"id"`
	RequestId int       `json:"requestId"`
	Step      string    `json:"step"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Created   time.Time `json:"created"`
}

// RequestDurationInfo represents request duration information
type RequestDurationInfo struct {
	RequestId       int       `json:"requestId"`
	GroupId         string    `json:"groupId"`
	Endpoint        string    `json:"endpoint"`
	Status          string    `json:"status"`
	DurationMinutes float64   `json:"durationMinutes"`
	Created         time.Time `json:"created"`
	Modified        time.Time `json:"modified"`
}

// ViewErrorRequest represents an error request from the view
type ViewErrorRequest struct {
	Id          int       `json:"id"`
	GroupId     string    `json:"groupId"`
	Endpoint    string    `json:"endpoint"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	RetryCount  int       `json:"retryCount"`
	Hidden      bool      `json:"hidden"`
	InitiatedBy string    `json:"initiatedBy"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// ViewQueuedRequest represents a queued request from the view
type ViewQueuedRequest struct {
	Id          int       `json:"id"`
	GroupId     string    `json:"groupId"`
	Endpoint    string    `json:"endpoint"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	InitiatedBy string    `json:"initiatedBy"`
	Created     time.Time `json:"created"`
}

// ViewRunningRequest represents a running request from the view
type ViewRunningRequest struct {
	Id          int       `json:"id"`
	GroupId     string    `json:"groupId"`
	Endpoint    string    `json:"endpoint"`
	Status      string    `json:"status"`
	InitiatedBy string    `json:"initiatedBy"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// ============================================================================
// Tool Request Models
// ============================================================================

// TblToolRequest represents a tool request
type TblToolRequest struct {
	Id          int       `json:"id"`
	ToolId      int       `json:"toolId"`
	GroupId     string    `json:"groupId"`
	Status      string    `json:"status"`
	RequestData string    `json:"requestData"`
	Message     string    `json:"message"`
	InitiatedBy string    `json:"initiatedBy"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// ============================================================================
// Tool Models
// ============================================================================

// ToolFullDetails represents full tool details with rules and requirements
type ToolFullDetails struct {
	Id                   int                          `json:"id"`
	ToolName             string                       `json:"toolName"`
	ToolDescription      string                       `json:"toolDescription"`
	TopicName            string                       `json:"topicName"`
	InfoPageUrl          string                       `json:"infoPageUrl"`
	FormTemplate         string                       `json:"formTemplate"`
	CurrentTemplateId    int                          `json:"currentTemplateId"`
	Enabled              bool                         `json:"enabled"`
	RequiresArchiving    bool                         `json:"requiresArchiving"`
	InstanceCount        int                          `json:"instanceCount"`
	RequestCount         int                          `json:"requestCount"`
	Rules                []TblToolRuleLogic           `json:"rules"`
	ExtendedRequirements []TblToolExtendedRequirement `json:"extendedRequirements"`
}

// ToolUpdateDto represents the DTO for updating tool properties
type ToolUpdateDto struct {
	ToolName          string `json:"toolName,omitempty"`
	ToolDescription   string `json:"toolDescription,omitempty"`
	TopicName         string `json:"topicName,omitempty"`
	InfoPageUrl       string `json:"infoPageUrl,omitempty"`
	FormTemplate      string `json:"formTemplate,omitempty"`
	CurrentTemplateId int    `json:"currentTemplateId,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	RequiresArchiving bool   `json:"requiresArchiving,omitempty"`
}

// ============================================================================
// Tool Instance Models
// ============================================================================

// TblToolInstance represents a tool instance
type TblToolInstance struct {
	Id              int       `json:"id"`
	ToolId          int       `json:"toolId"`
	GroupId         string    `json:"groupId"`
	TemplateVersion int       `json:"templateVersion"`
	Status          string    `json:"status"`
	Created         time.Time `json:"created"`
	Modified        time.Time `json:"modified"`
}

// TblToolMetaDatum represents metadata for a tool instance
type TblToolMetaDatum struct {
	Id             int    `json:"id"`
	ToolInstanceId int    `json:"toolInstanceId"`
	Key            string `json:"key"`
	Value          string `json:"value"`
}

// ============================================================================
// Team/Group Models
// ============================================================================

// ManagedTeam represents a managed team
type ManagedTeam struct {
	GroupId     string    `json:"groupId"`
	TeamName    string    `json:"teamName"`
	ProjectNo   string    `json:"projectNo"`
	ProjectName string    `json:"projectName"`
	Status      string    `json:"status"`
	Origin      string    `json:"origin"`
	Retention   string    `json:"retention"`
	SiteId      string    `json:"siteId"`
	Url         string    `json:"url"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// TeamFullDetails represents full team details
type TeamFullDetails struct {
	GroupId            string            `json:"groupId"`
	TeamName           string            `json:"teamName"`
	ProjectNo          string            `json:"projectNo"`
	ProjectName        string            `json:"projectName"`
	Status             string            `json:"status"`
	Origin             string            `json:"origin"`
	Retention          string            `json:"retention"`
	SiteId             string            `json:"siteId"`
	Url                string            `json:"url"`
	ToolInstances      []TblToolInstance `json:"toolInstances"`
	RecentRequests     []Request         `json:"recentRequests"`
	RecentToolRequests []TblToolRequest  `json:"recentToolRequests"`
}

// TeamSearchResult represents a team search result
type TeamSearchResult struct {
	GroupId      string `json:"groupId"`
	TeamName     string `json:"teamName"`
	ProjectNo    string `json:"projectNo"`
	ProjectName  string `json:"projectName"`
	Status       string `json:"status"`
	Origin       string `json:"origin"`
	MatchedField string `json:"matchedField"`
}

// TblGroup represents a group record
type TblGroup struct {
	GroupId     string    `json:"groupId"`
	TeamName    string    `json:"teamName"`
	ProjectNo   string    `json:"projectNo"`
	ProjectName string    `json:"projectName"`
	Status      string    `json:"status"`
	Origin      string    `json:"origin"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// ViewGroup represents a group from ViewGroups
type ViewGroup struct {
	GroupId     string `json:"groupId"`
	TeamName    string `json:"teamName"`
	ProjectNo   string `json:"projectNo"`
	ProjectName string `json:"projectName"`
	Status      string `json:"status"`
	Origin      string `json:"origin"`
}

// ============================================================================
// Archive & Export Job Models
// ============================================================================

// ArchiveJob represents an archive job
type ArchiveJob struct {
	Id             int             `json:"id"`
	GroupId        string          `json:"groupId"`
	Status         string          `json:"status"`
	JobType        string          `json:"jobType"`
	Message        string          `json:"message"`
	Created        time.Time       `json:"created"`
	Modified       time.Time       `json:"modified"`
	ArchiveSubJobs []ArchiveSubJob `json:"archiveSubJobs,omitempty"`
}

// ArchiveSubJob represents a sub-job of an archive job
type ArchiveSubJob struct {
	Id           int       `json:"id"`
	ArchiveJobId int       `json:"archiveJobId"`
	SubJobType   string    `json:"subJobType"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}

// RequiredArchiveJob represents a required archive job
type RequiredArchiveJob struct {
	GroupId    string    `json:"groupId"`
	TeamName   string    `json:"teamName"`
	Retention  string    `json:"retention"`
	ExpiryDate time.Time `json:"expiryDate"`
	Status     string    `json:"status"`
}

// ExportDataJob represents an export data job
type ExportDataJob struct {
	Id        int       `json:"id"`
	RequestId int       `json:"requestId"`
	GroupId   string    `json:"groupId"`
	Status    string    `json:"status"`
	FilePath  string    `json:"filePath"`
	FileSize  int64     `json:"fileSize"`
	Message   string    `json:"message"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
}

// ClearSiteJob represents a clear site job
type ClearSiteJob struct {
	Id              int       `json:"id"`
	GroupId         string    `json:"groupId"`
	Status          string    `json:"status"`
	StorageReleased int64     `json:"storageReleased"`
	FilesDeleted    int       `json:"filesDeleted"`
	Message         string    `json:"message"`
	Created         time.Time `json:"created"`
	Modified        time.Time `json:"modified"`
}

// ClearSiteJobsSummary represents clear site jobs summary
type ClearSiteJobsSummary struct {
	TotalJobs            int   `json:"totalJobs"`
	CompletedJobs        int   `json:"completedJobs"`
	PendingJobs          int   `json:"pendingJobs"`
	FailedJobs           int   `json:"failedJobs"`
	TotalStorageReleased int64 `json:"totalStorageReleased"`
	TotalFilesDeleted    int   `json:"totalFilesDeleted"`
}

// ============================================================================
// Logging Models
// ============================================================================

// TblToolBoxLogger represents a log entry
type TblToolBoxLogger struct {
	Id          int       `json:"id"`
	Subject     string    `json:"subject"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	InitiatedBy string    `json:"initiatedBy"`
	Created     time.Time `json:"created"`
}

// LogEntry represents a new log entry to be created
type LogEntry struct {
	Subject     string `json:"subject"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	InitiatedBy string `json:"initiatedBy"`
}

// ============================================================================
// Rules & Logic Models
// ============================================================================

// TblToolRule represents a tool rule
type TblToolRule struct {
	Id       int    `json:"id"`
	RuleName string `json:"ruleName"`
}

// TblToolRuleLogic represents rule logic for a tool
type TblToolRuleLogic struct {
	Id       int    `json:"id"`
	ToolId   int    `json:"toolId"`
	RuleId   int    `json:"ruleId"`
	RuleName string `json:"ruleName,omitempty"`
	Logic    string `json:"logic"`
	Value    string `json:"value"`
}

// RuleLogicCreate represents the DTO for creating rule logic
type RuleLogicCreate struct {
	ToolId int    `json:"toolId"`
	RuleId int    `json:"ruleId"`
	Logic  string `json:"logic"`
	Value  string `json:"value"`
}

// ============================================================================
// Extended Requirements Models
// ============================================================================

// TblToolExtendedRequirement represents an extended requirement
type TblToolExtendedRequirement struct {
	Id               int    `json:"id"`
	ToolId           int    `json:"toolId"`
	RequirementName  string `json:"requirementName"`
	RequirementValue string `json:"requirementValue"`
}

// ExtendedRequirementCreate represents the DTO for creating an extended requirement
type ExtendedRequirementCreate struct {
	ToolId           int    `json:"toolId"`
	RequirementName  string `json:"requirementName"`
	RequirementValue string `json:"requirementValue"`
}

// ============================================================================
// Background Task Models
// ============================================================================

// TblTask represents a background task
type TblTask struct {
	Id        int       `json:"id"`
	ProjectNo string    `json:"projectNo"`
	JobType   int       `json:"jobType"`
	Status    int       `json:"status"`
	Data      string    `json:"data"`
	Message   string    `json:"message"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
}

// ============================================================================
// GeoBIM Models
// ============================================================================

// ViewToolGeoBim represents the GeoBIM view
type ViewToolGeoBim struct {
	GroupId     string `json:"groupId"`
	TeamName    string `json:"teamName"`
	ProjectNo   string `json:"projectNo"`
	ProjectName string `json:"projectName"`
	GeoBimData  string `json:"geoBimData"`
}

// ============================================================================
// Health & Cleanup Models
// ============================================================================

// HealthStatus represents database health status
type HealthStatus struct {
	Status string `json:"status"`
}

// OrphanedRecordsSummary represents orphaned records summary
type OrphanedRecordsSummary struct {
	OrphanedToolInstances int `json:"orphanedToolInstances"`
	OrphanedRequests      int `json:"orphanedRequests"`
	OrphanedToolRequests  int `json:"orphanedToolRequests"`
	OrphanedRequestSteps  int `json:"orphanedRequestSteps"`
	OrphanedToolMetadata  int `json:"orphanedToolMetadata"`
}

// ============================================================================
// Bulk Operation Models
// ============================================================================

// BulkOperationError represents an error in a bulk operation
type BulkOperationError struct {
	Id    int    `json:"id"`
	Error string `json:"error"`
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	TotalRequested int                  `json:"totalRequested"`
	Succeeded      int                  `json:"succeeded"`
	Failed         int                  `json:"failed"`
	Errors         []BulkOperationError `json:"errors"`
}

// ============================================================================
// Request/Response DTOs
// ============================================================================

// StatusUpdateRequest represents a status update request body
type StatusUpdateRequest struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// PriorityUpdateRequest represents a priority update request body
type PriorityUpdateRequest struct {
	Priority int `json:"priority"`
}

// HiddenUpdateRequest represents a hidden update request body
type HiddenUpdateRequest struct {
	Hidden bool `json:"hidden"`
}

// BulkRetryRequest represents a bulk retry request body
type BulkRetryRequest struct {
	Ids []int `json:"ids"`
}

// BulkHideRequest represents a bulk hide request body
type BulkHideRequest struct {
	Ids    []int `json:"ids"`
	Hidden bool  `json:"hidden"`
}

// AddStepRequest represents a request to add a step
type AddStepRequest struct {
	Step    string `json:"step"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// EnabledUpdateRequest represents an enabled update request body
type EnabledUpdateRequest struct {
	Enabled bool `json:"enabled"`
}

// TopicUpdateRequest represents a topic update request body
type TopicUpdateRequest struct {
	TopicName string `json:"topicName"`
}

// TemplateUpdateRequest represents a template update request body
type TemplateUpdateRequest struct {
	TemplateId int `json:"templateId"`
}

// BulkStatusUpdateRequest represents a bulk status update request body
type BulkStatusUpdateRequest struct {
	Ids    []int  `json:"ids"`
	Status string `json:"status"`
}

// BulkToolIdsRequest represents a bulk tool IDs request body
type BulkToolIdsRequest struct {
	Ids []int `json:"ids"`
}

// RuleCreate represents the DTO for creating a rule
type RuleCreate struct {
	RuleName string `json:"ruleName"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
