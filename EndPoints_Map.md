# Admin API Client Cheat Sheet

This document maps API endpoints to their corresponding Go client methods in `TeamToolBoxHelper/AdminApi.go` and explains how to handle/print the responses.

## Usage Pattern

```go
import "github.com/saricon83-sudo/Tyr365AdminCli/TeamToolBoxHelper"

// 1. Create client
client, _ := teamToolboxHelper.CreateAdminAPI()

// 2. Call method
resp, err := client.GetDashboard()

// 3. Print result (using receiver or helper function)
resp.PrintTable()
```

## 1. Dashboard & Statistics

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/dashboard` | `GetDashboard()` | `*AdminDashboardStats` | `resp.PrintTable()` |
| `/stats/requests-by-status` | `GetRequestsByStatus()` | `map[string]int` | *(Iterate map manually)* |
| `/stats/requests-by-tool` | `GetRequestsByTool()` | `[]ToolRequestCount` | `teamToolboxHelper.PrintToolRequestCountTable(resp)` |
| `/stats/requests-by-day` | `GetRequestsByDay(days)` | `[]DailyRequestCount` | `teamToolboxHelper.PrintDailyRequestCountTable(resp)` |
| `/stats/tool-adoption` | `GetToolAdoption()` | `[]ToolAdoptionStats` | `teamToolboxHelper.PrintToolAdoptionStatsTable(resp)` |
| `/stats/storage-released` | `GetStorageReleased()` | `*StorageReleasedSummary` | `resp.PrintTable()` |
| `/stats/archive-jobs` | `GetArchiveJobStats()` | `*ArchiveJobStats` | `resp.PrintTable()` |
| `/stats/pending` | `GetPendingCounts()` | `*PendingCountsSummary` | `resp.PrintTable()` |

## 2. Error & Stuck Request Management

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/errors` | `GetErrors(includeHidden)` | `[]ViewErrorRequest` | `teamToolboxHelper.PrintViewErrorRequestTable(resp)` |
| `/queued` | `GetQueued()` | `[]ViewQueuedRequest` | `teamToolboxHelper.PrintViewQueuedRequestTable(resp)` |
| `/running` | `GetRunning()` | `[]ViewRunningRequest` | `teamToolboxHelper.PrintViewRunningRequestTable(resp)` |
| `/stuck` | `GetStuck(hours)` | `[]Request` | `teamToolboxHelper.PrintRequestTable(resp)` |
| `/high-retry` | `GetHighRetry(minRetries)` | `[]Request` | `teamToolboxHelper.PrintRequestTable(resp)` |
| `/slowest` | `GetSlowest(count)` | `[]RequestDurationInfo` | `teamToolboxHelper.PrintRequestDurationInfoTable(resp)` |

## 3. Request Management

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/requests/{id}` | `GetRequest(id)` | `*Request` | `resp.PrintTable()` |
| `/requests/by-group` | `GetRequestsByGroup(groupId)` | `[]Request` | `teamToolboxHelper.PrintRequestTable(resp)` |
| `/requests/by-initiator` | `GetRequestsByInitiator(email)` | `[]Request` | `teamToolboxHelper.PrintRequestTable(resp)` |
| `/requests/by-endpoint` | `GetRequestsByEndpoint(name)` | `[]Request` | `teamToolboxHelper.PrintRequestTable(resp)` |
| `/requests/{id}/steps` | `GetRequestSteps(id)` | `[]RequestStep` | `teamToolboxHelper.PrintRequestStepTable(resp)` |
| `/requests/{id}/status` | `UpdateRequestStatus(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/requests/{id}/priority` | `UpdateRequestPriority(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/requests/{id}/retry` | `RetryRequest(id)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/requests/{id}/hidden` | `UpdateRequestHidden(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/requests/bulk-retry` | `BulkRetryRequests(ids)` | `*BulkOperationResult` | `resp.PrintTable()` |
| `/requests/bulk-hide` | `BulkHideRequests(...)` | `*BulkOperationResult` | `resp.PrintTable()` |
| `/requests/{id}/steps` | `AddRequestStep(...)` | `*RequestStep` | *(No table printer available)* |

## 4. Tool Requests

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/tool-requests` | `GetToolRequests(status, group)` | `[]TblToolRequest` | `teamToolboxHelper.PrintTblToolRequestTable(resp)` |
| `/tool-requests/{id}` | `GetToolRequest(id)` | `*TblToolRequest` | *(Wrap in slice or print manually)* |
| `/tool-requests/{id}/status` | `UpdateToolRequestStatus(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/tool-requests/count` | `GetToolRequestCountByStatus()` | `map[string]int` | *(Iterate map manually)* |

## 5. Tool Management

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/tools` | `GetTools()` | `[]TblTool` | `teamToolboxHelper.TblTools(resp).PrintTable()` |
| `/tools/{id}` | `GetToolFullDetails(id)` | `*ToolFullDetails` | `resp.PrintTable()` |
| `/tools/{id}/enabled` | `UpdateToolEnabled(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/tools/{id}/topic` | `UpdateToolTopic(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/tools/{id}/template` | `UpdateToolTemplate(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/tools/{id}` | `UpdateTool(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/tools/{id}/instances` | `GetToolInstances(toolId)` | `[]TblToolInstance` | `teamToolboxHelper.PrintTblToolInstanceTable(resp)` |
| `/tools/{id}/requests` | `GetToolRequestsForTool(id)` | `[]TblToolRequest` | `teamToolboxHelper.PrintTblToolRequestTable(resp)` |
| `/tools/unused` | `GetUnusedTools(days)` | `[]TblTool` | `teamToolboxHelper.TblTools(resp).PrintTable()` |
| `/tools` | `AddTool(tool)` | `*TblTool` | `resp.PrintTable()` |

## 6. Tool Instances

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/instances` | `GetInstances(toolId, group)` | `[]TblToolInstance` | `teamToolboxHelper.PrintTblToolInstanceTable(resp)` |
| `/instances/{id}` | `GetInstance(id)` | `*TblToolInstance` | *(Wrap in slice or print manually)* |
| `/instances/{id}/metadata` | `GetInstanceMetadata(id)` | `[]TblToolMetaDatum` | `teamToolboxHelper.PrintTblToolMetaDatumTable(resp)` |
| `/instances/{id}` | `DeleteInstance(id)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/instances/orphaned` | `GetOrphanedInstances()` | `[]TblToolInstance` | `teamToolboxHelper.PrintTblToolInstanceTable(resp)` |
| `/instances/outdated` | `GetOutdatedInstances()` | `[]TblToolInstance` | `teamToolboxHelper.PrintTblToolInstanceTable(resp)` |
| `/instances` | `AddInstance(instance)` | `*TblToolInstance` | *(Wrap in slice or print manually)* |

## 7. Team/Group Management

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/teams` | `GetTeams(status, origin)` | `[]ManagedTeam` | `teamToolboxHelper.PrintManagedTeamTable(resp)` |
| `/teams/{groupId}` | `GetTeam(groupId)` | `*ManagedTeam` | *(Wrap in slice or print manually)* |
| `/teams/{groupId}/details` | `GetTeamFullDetails(group)` | `*TeamFullDetails` | `resp.PrintTable()` |
| `/teams/{groupId}/instances`| `GetTeamInstances(group)` | `[]TblToolInstance` | `teamToolboxHelper.PrintTblToolInstanceTable(resp)` |
| `/teams/{groupId}/requests` | `GetTeamRequests(group)` | `[]Request` | `teamToolboxHelper.PrintRequestTable(resp)` |
| `/teams/{groupId}/status` | `UpdateTeamStatus(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/teams/without-tools` | `GetTeamsWithoutTools()` | `[]ManagedTeam` | `teamToolboxHelper.PrintManagedTeamTable(resp)` |
| `/teams/search` | `SearchTeams(query)` | `[]TeamSearchResult` | `teamToolboxHelper.PrintTeamSearchResultTable(resp)` |
| `/groups/{groupId}` | `GetGroup(groupId)` | `*TblGroup` | *(No printer available)* |
| `/groups/view` | `GetGroupsView()` | `[]ViewGroup` | *(No printer available)* |

## 8. Archive & Export Jobs

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/archive-jobs` | `GetArchiveJobs(status)` | `[]ArchiveJob` | `teamToolboxHelper.PrintArchiveJobTable(resp)` |
| `/archive-jobs/{id}` | `GetArchiveJob(id)` | `*ArchiveJob` | `resp.PrintTable()` |
| `/archive-jobs/{id}/status` | `UpdateArchiveJobStatus(...)` | `*MessageResponse` | `fmt.Println(resp.Message)` |
| `/archive-jobs/required` | `GetRequiredArchiveJobs()` | `[]RequiredArchiveJob` | `teamToolboxHelper.PrintRequiredArchiveJobTable(resp)` |
| `/export-jobs` | `GetExportJobs(groupId)` | `[]ExportDataJob` | `teamToolboxHelper.PrintExportDataJobTable(resp)` |
| `/clear-site-jobs` | `GetClearSiteJobs(status)` | `[]ClearSiteJob` | `teamToolboxHelper.PrintClearSiteJobTable(resp)` |
| `/clear-site-jobs/summary` | `GetClearSiteJobsSummary()` | `*ClearSiteJobsSummary` | `resp.PrintTable()` |

## 9. Logging

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/logs` | `GetLogs(...)` | `[]TblToolBoxLogger` | `teamToolboxHelper.PrintTblToolBoxLoggerTable(resp)` |
| `/logs/recent` | `GetRecentLogs(count)` | `[]TblToolBoxLogger` | `teamToolboxHelper.PrintTblToolBoxLoggerTable(resp)` |
| `/logs/by-subject` | `GetLogsBySubject(subject)` | `[]TblToolBoxLogger` | `teamToolboxHelper.PrintTblToolBoxLoggerTable(resp)` |
| `/logs` | `AddLog(entry)` | `*TblToolBoxLogger` | *(Wrap in slice or print manually)* |

## 10. Rules & Logic

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/rules` | `GetRules()` | `[]TblToolRule` | `teamToolboxHelper.PrintTblToolRuleTable(resp)` |
| `/rules/{id}` | `GetRule(id)` | `*TblToolRule` | *(Wrap in slice)* |
| `/rules/unused` | `GetUnusedRules()` | `[]TblToolRule` | `teamToolboxHelper.PrintTblToolRuleTable(resp)` |
| `/rule-logics` | `GetRuleLogics()` | `[]TblToolRuleLogic` | `teamToolboxHelper.PrintTblToolRuleLogicTable(resp)` |
| `/rule-logics/by-tool` | `GetRuleLogicsByTool(id)` | `[]TblToolRuleLogic` | `teamToolboxHelper.PrintTblToolRuleLogicTable(resp)` |

## 11. Extended Requirements

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/requirements/by-tool` | `GetExtendedRequirementsByTool`| `[]TblToolExtendedRequirement` | `teamToolboxHelper.PrintTblToolExtendedRequirementTable(resp)` |

## 12. Background Tasks

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/tasks` | `GetTasks(...)` | `[]TblTask` | `teamToolboxHelper.PrintTblTaskTable(resp)` |
| `/tasks/pending` | `GetPendingTasks()` | `[]TblTask` | `teamToolboxHelper.PrintTblTaskTable(resp)` |
| `/tasks/failed` | `GetFailedTasks()` | `[]TblTask` | `teamToolboxHelper.PrintTblTaskTable(resp)` |
| `/tasks/by-project` | `GetTasksByProject(no)` | `[]TblTask` | `teamToolboxHelper.PrintTblTaskTable(resp)` |

## 13. GeoBIM

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/geobim` | `GetGeoBIM()` | `[]ViewToolGeoBim` | `teamToolboxHelper.PrintViewToolGeoBimTable(resp)` |
| `/geobim/{groupId}` | `GetGeoBIMByGroup(group)` | `*ViewToolGeoBim` | *(Wrap in slice)* |

## 14. Health & Cleanup

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/health/db` | `GetDatabaseHealth()` | `*HealthStatus` | `resp.PrintTable()` |
| `/cleanup/orphaned` | `GetOrphanedRecords()` | `*OrphanedRecordsSummary` | `resp.PrintTable()` |

## 15. Bulk Operations

| Endpoint | Go Method | Return Type | Printer / Next Step |
| :--- | :--- | :--- | :--- |
| `/bulk/update-status` | `BulkUpdateStatus(...)` | `*BulkOperationResult` | `resp.PrintTable()` |
| `/bulk/enable-tools` | `BulkEnableTools(...)` | `*BulkOperationResult` | `resp.PrintTable()` |
| `/bulk/disable-tools` | `BulkDisableTools(...)` | `*BulkOperationResult` | `resp.PrintTable()` |
