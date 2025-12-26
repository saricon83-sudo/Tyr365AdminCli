package GraphHelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/saricon83-sudo/Tyr365AdminCli/internal/auth"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
	"github.com/google/uuid"
	bmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type NewAssignment struct {
	ODataType string `json:"@odata.type"`
	OrderHint string `json:"orderHint"`
}

// Constructor function for NewAssignment
func NewNewAssignment() *NewAssignment {
	return &NewAssignment{
		ODataType: "#microsoft.graph.plannerAssignment",
		OrderHint: " !",
	}
}

type NewTaskPayload struct {
	PlanId      string                            `json:"planId"`
	BucketId    string                            `json:"bucketId"`
	Title       string                            `json:"title"`
	Assignments map[string]map[string]interface{} `json:"assignments"`
}
type PlannerTask struct {
	PlanId      string                    `json:"planId"`
	BucketId    string                    `json:"bucketId"`
	Title       string                    `json:"title"`
	Assignments map[string]*NewAssignment `json:"assignments"`
}
type ChecklistItem struct {
	ODataType string `json:"@odata.type"`
	Title     string `json:"title"`
	IsChecked bool   `json:"isChecked"`
}

type TaskDetailsUpdate struct {
	Checklist map[string]interface{} `json:"checklist"`
}

type User struct {
	ID                              string   `json:"id,omitempty"`
	DeletedDateTime                 string   `json:"deletedDateTime,omitempty"`
	AccountEnabled                  bool     `json:"accountEnabled,omitempty"`
	AgeGroup                        string   `json:"ageGroup,omitempty"`
	BusinessPhones                  []string `json:"businessPhones,omitempty"`
	City                            string   `json:"city,omitempty"`
	CreatedDateTime                 string   `json:"createdDateTime,omitempty"`
	CreationType                    *string  `json:"creationType,omitempty"`
	CompanyName                     string   `json:"companyName,omitempty"`
	ConsentProvidedForMinor         *string  `json:"consentProvidedForMinor,omitempty"`
	Country                         string   `json:"country,omitempty"`
	Department                      string   `json:"department,omitempty"`
	DisplayName                     string   `json:"displayName,omitempty"`
	EmployeeId                      *string  `json:"employeeId,omitempty"`
	EmployeeHireDate                *string  `json:"employeeHireDate,omitempty"`
	EmployeeLeaveDateTime           *string  `json:"employeeLeaveDateTime,omitempty"`
	EmployeeType                    *string  `json:"employeeType,omitempty"`
	FaxNumber                       *string  `json:"faxNumber,omitempty"`
	GivenName                       string   `json:"givenName,omitempty"`
	ImAddresses                     []string `json:"imAddresses,omitempty"`
	InfoCatalogs                    []string `json:"infoCatalogs,omitempty"`
	IsLicenseReconciliationNeeded   bool     `json:"isLicenseReconciliationNeeded,omitempty"`
	IsManagementRestricted          *bool    `json:"isManagementRestricted,omitempty"`
	IsResourceAccount               *bool    `json:"isResourceAccount,omitempty"`
	JobTitle                        string   `json:"jobTitle,omitempty"`
	LegalAgeGroupClassification     string   `json:"legalAgeGroupClassification,omitempty"`
	Mail                            string   `json:"mail,omitempty"`
	MailNickname                    string   `json:"mailNickname,omitempty"`
	MobilePhone                     string   `json:"mobilePhone,omitempty"`
	OnPremisesDistinguishedName     *string  `json:"onPremisesDistinguishedName,omitempty"`
	OfficeLocation                  string   `json:"officeLocation,omitempty"`
	OnPremisesDomainName            *string  `json:"onPremisesDomainName,omitempty"`
	OnPremisesImmutableId           *string  `json:"onPremisesImmutableId,omitempty"`
	OnPremisesLastSyncDateTime      *string  `json:"onPremisesLastSyncDateTime,omitempty"`
	OnPremisesObjectIdentifier      *string  `json:"onPremisesObjectIdentifier,omitempty"`
	OnPremisesSecurityIdentifier    *string  `json:"onPremisesSecurityIdentifier,omitempty"`
	OnPremisesSamAccountName        *string  `json:"onPremisesSamAccountName,omitempty"`
	OnPremisesSyncEnabled           *bool    `json:"onPremisesSyncEnabled,omitempty"`
	OnPremisesUserPrincipalName     *string  `json:"onPremisesUserPrincipalName,omitempty"`
	OtherMails                      []string `json:"otherMails,omitempty"`
	PasswordPolicies                *string  `json:"passwordPolicies,omitempty"`
	PostalCode                      string   `json:"postalCode,omitempty"`
	PreferredDataLocation           *string  `json:"preferredDataLocation,omitempty"`
	PreferredLanguage               string   `json:"preferredLanguage,omitempty"`
	ProxyAddresses                  []string `json:"proxyAddresses,omitempty"`
	RefreshTokensValidFromDateTime  string   `json:"refreshTokensValidFromDateTime,omitempty"`
	SecurityIdentifier              string   `json:"securityIdentifier,omitempty"`
	ShowInAddressList               *bool    `json:"showInAddressList,omitempty"`
	SignInSessionsValidFromDateTime string   `json:"signInSessionsValidFromDateTime,omitempty"`
	State                           string   `json:"state,omitempty"`
	StreetAddress                   string   `json:"streetAddress,omitempty"`
	Surname                         string   `json:"surname,omitempty"`
	UsageLocation                   string   `json:"usageLocation,omitempty"`
	UserPrincipalName               string   `json:"userPrincipalName,omitempty"`
	ExternalUserConvertedOn         *string  `json:"externalUserConvertedOn,omitempty"`
	ExternalUserState               *string  `json:"externalUserState,omitempty"`
	ExternalUserStateChangeDateTime *string  `json:"externalUserStateChangeDateTime,omitempty"`
	UserType                        string   `json:"userType,omitempty"`
	EmployeeOrgData                 *string  `json:"employeeOrgData,omitempty"`
}

func (g *GraphHelper) GetGroupById(groupId string) (models.Groupable, error) {
	group, err := g.appClient.Groups().ByGroupId(groupId).Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *GraphHelper) GetDeletedGroups() ([]bmodels.Groupable, error) {
	graphGroups, err := g.betaClient.Directory().DeletedItems().GraphGroup().Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	groups := graphGroups.GetValue()
	return groups, nil
}

func (g *GraphHelper) GetUsers(selectProperties []string, amount *int32, filter string) (models.UserCollectionResponseable, error) {
	var topValue int32
	if amount == nil {
		topValue = 25 // Default value if amount is not provided
	} else {
		topValue = *amount
	}

	query := users.UsersRequestBuilderGetQueryParameters{
		Select: selectProperties,
		Top:    &topValue,
		// Orderby: []string{"displayName"},
		Filter: &filter,
	}

	return g.appClient.Users().
		Get(context.Background(), &users.UsersRequestBuilderGetRequestConfiguration{
			QueryParameters: &query,
		})
}

func (g *GraphHelper) CreateTask(taskTitle string) (models.PlannerTaskable, error) {
	cfg := config.Get()
	requestBody := models.NewPlannerTask()

	planId := cfg.GetString("planId")
	bucketId := cfg.GetString("bucketId")

	requestBody.SetPlanId(&planId)
	requestBody.SetBucketId(&bucketId)
	requestBody.SetTitle(&taskTitle)
	result, err := g.appClient.Planner().Tasks().Post(context.Background(), requestBody, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *GraphHelper) GetAllTasks() ([]string, error) {
	cfg := config.Get()

	planId := cfg.GetString("planId")
	tasks, err := g.appClient.Planner().Plans().ByPlannerPlanId(planId).Tasks().Get(context.Background(), nil)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	var taskTitles []string
	taskValues := tasks.GetValue()
	for _, task := range taskValues {
		title := task.GetTitle()
		rdy := task.GetCompletedBy()
		if rdy == nil {

			taskTitles = append(taskTitles, *title)
		}
		if rdy != nil {
		}
	}
	return taskTitles, nil
}

func (g *GraphHelper) CreateTaskWithChecklist(title, checklistStr string) (string, error) {
	cfg := config.Get()

	planId := cfg.GetString("planId")
	bucketId := cfg.GetString("bucketId")
	accessToken, _ := AuthGraphApi()

	assignees := make(map[string]*NewAssignment)
	// Assign the task to a specific user by their ID
	assignees["fe429714-d600-4948-a412-b9983986356e"] = &NewAssignment{
		ODataType: "#microsoft.graph.plannerAssignment",
		OrderHint: " !",
	}

	// Prepare the task payload
	newTask := PlannerTask{
		PlanId:      planId,
		BucketId:    bucketId,
		Title:       title,
		Assignments: assignees,
	}

	taskBytes, err := json.Marshal(newTask)
	if err != nil {
		return "", err
	}

	taskID, err := createPlannerTask(taskBytes, accessToken)
	if err != nil {
		return "", err
	}

	err = UpdateTaskWithChecklistItems(taskID, checklistStr)
	if err != nil {
		return "", err
	}

	return taskID, nil
}

// Helper function to create a Planner task
func createPlannerTask(taskPayload []byte, accessToken string) (string, error) {
	req, err := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/planner/tasks", bytes.NewBuffer(taskPayload))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create task, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Assuming the task ID is available in the response
	taskID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("task ID not found in response")
	}

	return taskID, nil
}

func (g *GraphHelper) GetTeamById(teamId string) (models.Teamable, error) {
	team, nil := g.appClient.Teams().ByTeamId(teamId).Get(context.Background(), nil)
	return team, nil
}

func (g *GraphHelper) GetAllChannels(teamId string) (models.ChannelCollectionResponseable, error) {
	channels, nil := g.appClient.Teams().ByTeamId(teamId).AllChannels().Get(context.Background(), nil)
	return channels, nil
}

func (g *GraphHelper) GetChannelById(teamId string, channelId string) (models.Channelable, error) {
	channel, nil := g.appClient.Teams().ByTeamId(teamId).Channels().ByChannelId(channelId).Get(context.Background(), nil)
	return channel, nil
}

func (g *GraphHelper) EnsureFilesFolder(teamId string, channelId string) (models.DriveItemable, error) {
	drive, nil := g.appClient.Teams().ByTeamId(teamId).Channels().ByChannelId(channelId).FilesFolder().Get(context.Background(), nil)
	return drive, nil
}
func (g *GraphHelper) EnsureOneNote(teamId string, channelId string) (models.NotebookCollectionResponseable, error) {
	Onenote, nil := g.appClient.Groups().ByGroupId(teamId).Onenote().Notebooks().Get(context.Background(), nil)
	return Onenote, nil
}
func (g *GraphHelper) GetTabs(teamId string, channelId string) (models.TeamsTabCollectionResponseable, error) {
	teamTabs, nil := g.appClient.Teams().ByTeamId(teamId).Channels().ByChannelId(channelId).Tabs().Get(context.Background(), nil)
	return teamTabs, nil
}

func (g *GraphHelper) GetTeamsInstalledApps (teamId string) (models.TeamsAppInstallationCollectionResponseable, error){
 installedApps, err := g.appClient.Teams().ByTeamId(teamId).InstalledApps().Get(context.Background(), nil)
 if err != nil {
  return nil, err
 }

 return installedApps, nil
}

func UpdateTaskWithChecklistItems(taskID, checklistStr string) error {
	// Split the checklistStr into individual titles
	titles := strings.Split(checklistStr, ",")

	// Initialize the checklist map
	checklist := make(map[string]interface{})
	for _, title := range titles {
			checklistItemId := uuid.New().String()
		checklist[checklistItemId] = ChecklistItem{
			ODataType: "microsoft.graph.plannerChecklistItem",
			Title:     strings.TrimSpace(title),
			IsChecked: false,
		}
	}

	// Prepare the update payload
	updatePayload := TaskDetailsUpdate{
		Checklist: checklist,
	}

	updateBytes, err := json.Marshal(updatePayload)
	if err != nil {
		return fmt.Errorf("error marshalling update payload: %v", err)
	}

	// Prepare the PATCH request
	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://graph.microsoft.com/v1.0/planner/tasks/%s/details", taskID), bytes.NewBuffer(updateBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	accessToken, err := AuthGraphApi()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	eTag, err := GetTaskETag(taskID)
	fmt.Println(eTag)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("If-Match", eTag)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	fmt.Println("Checklist items added successfully.")
	return nil
}

type TokenCached struct {
	Token string
}

var TokenCache string

func makePOSTRequest(postUrl string, bodyValues []byte) (*http.Response, error) {
	// Encode the body values into a URL-encoded format
	body := bytes.NewBuffer(bodyValues)

	// Create the request
	req, err := http.NewRequest("POST", postUrl, body)
	if err != nil {
		return nil, err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AuthGraphApi gets an authentication token for the Graph API.
// Deprecated: Use auth.GetGraphToken() directly instead.
func AuthGraphApi() (string, error) {
	return auth.GetGraphToken()
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

// GetUserOwnedObjects retrieves all objects owned by a user (groups, applications, etc.)
func (g *GraphHelper) GetUserOwnedObjects(userId string) ([]models.DirectoryObjectable, error) {
	ownedObjects, err := g.appClient.Users().ByUserId(userId).OwnedObjects().Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned objects: %w", err)
	}
	return ownedObjects.GetValue(), nil
}
