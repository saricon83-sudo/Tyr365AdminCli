package teamGovHttp

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type Request struct {
	ID               int          `json:"Id"`
	Created          string       `json:"Created"`
	GroupID          string       `json:"GroupId"`
	TeamName         string       `json:"TeamName"`
	Endpoint         string       `json:"Endpoint"`
	CallerID         string       `json:"CallerId"`
	Parameters       string       `json:"Parameters"`
	Status           string       `json:"Status"`
	ProvisioningStep string       `json:"ProvisioningStep"`
	Message          string       `json:"Message"`
	InitiatedBy      string       `json:"InitiatedBy"`
	Modified         string       `json:"Modified"`
	ClientTaskID     int          `json:"ClientTaskId"`
	LtpmessageSent   bool         `json:"LtpmessageSent"`
	Hidden           bool         `json:"Hidden"`
	RetryCount       int          `json:"RetryCount"`
	QueuePriority    int          `json:"QueuePriority"`
	GroupDetails     GroupDetails `json:"GroupDetails"`
}
type Parameters struct {
	GroupID        string `json:"groupId"`
	TemplateId     int    `json:"templateId"`
	Description    string `json:"description"`
	CallerId       string `json:"callerId"`
	InitiatedBy    string `json:"initiatedBy"`
	FlowParameters struct {
		TemplateID      string `json:"templateID"`
		ProjectNumber   string `json:"ProjectNumber"`
		TyrAProcessType string `json:"TyrAProcessType"`
	} `json:"flowParameters"`
	ClientTaskId int `json:"clientTaskId"`
	// Add other fields as needed
}
type UnifiedGroupSlice []UnifiedGroup

type UnifiedGroup struct {
	GroupId            string      `json:"groupId"`
	DisplayName        string      `json:"displayName"`
	Alias              string      `json:"alias"`
	Description        string      `json:"description"`
	CreatedDate        string      `json:"createdDate"`
	SharePointUrl      string      `json:"sharePointUrl"`
	Visibility         string      `json:"visibility"`
	Team               string      `json:"team"`
	Yammer             interface{} `json:"yammer"`
	Label              interface{} `json:"label"`
	ExpirationDateTime interface{} `json:"expirationDateTime"`
	ExchangeProperties interface{} `json:"exchangeProperties"`
}
type GroupDetails struct {
	GroupID            string      `json:"groupId"`
	DisplayName        string      `json:"displayName"`
	Alias              string      `json:"alias"`
	Description        string      `json:"description"`
	CreatedDate        string      `json:"createdDate"`
	SharePointURL      string      `json:"sharePointUrl"`
	Visibility         string      `json:"visibility"`
	Team               string      `json:"team"`
	Yammer             string      `json:"yammer"`
	Label              string      `json:"label"`
	ExpirationDateTime string      `json:"expirationDateTime"`
	ExchangeProperties interface{} `json:"exchangeProperties"`
}
type ManagedTeam struct {
	Id        int    `json:"Id"`
	GroupId   string `json:"groupId"`
	TeamName  string `json:"teamName"`
	Status    string `json:"status"`
	Origin    string `json:"origin"`
	Retention string `json:"retention"`
}
type ManagedTeamSlice []ManagedTeam

func (m *ManagedTeamSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"GroupId", "TeamName", "Status", "Origin", "Retention"}) // Customize the table header as needed

	// Populate the table with data from the response
	for _, req := range *m {
		row := []string{
			// fmt.Sprintf("%d", req.Id),
			req.GroupId,
			req.TeamName,
			req.Status,
			req.Origin,
			req.Retention,
		}
		table.Append(row)
	}

	// Render the table
	table.Render()

}

type TokenCached struct {
	Token string
}

func (u *UnifiedGroupSlice) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"GroupId", "DisplayName", "Alias", "Description", "CreatedDate", "SharePointUrl", "Visibility", "Team", "Yammer", "Label"}) // Customize the table header as needed

	// Populate the table with data from the response
	for _, req := range *u {
		row := []string{
			req.GroupId,
			req.DisplayName,
			req.Alias,
			req.Description,
			req.CreatedDate,
			req.SharePointUrl,
			req.Visibility,
			req.Team,
			fmt.Sprintf("%v", req.Yammer),
			fmt.Sprintf("%v", req.Label),
		}
		table.Append(row)
	}

	// Render the table
	table.Render()

}
