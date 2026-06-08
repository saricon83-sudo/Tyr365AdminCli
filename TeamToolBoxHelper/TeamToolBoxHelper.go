package teamToolboxHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func (client *APIClient) GetTestPolicy() (string, error) {
	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Get(fmt.Sprintf("%s/Tools/TestPolicy", client.BaseURL))
	if err != nil {
		return "", fmt.Errorf("unable to call API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response: %w", err)
	}

	return string(body), nil
}

func (client *APIClient) GetToolById(id string) (*TblTool, error) {
	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}
	var adress string = fmt.Sprintf("%s/Admin/"+id, client.BaseURL)
	resp, err := httpClient.Get(adress)
	if err != nil {
		return nil, fmt.Errorf("unable to call public API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read public response: %w", err)
	}
	var tool TblTool
	err = json.Unmarshal(body, &tool)
	if err != nil {
		fmt.Println(err)
	}
	return &tool, nil
}

func (client *APIClient) GetRulesAndLogic() (*RulesandLogics, error) {
	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}
	var adress string = fmt.Sprintf("%s/Admin/RulesandLogic", client.BaseURL)
	resp, err := httpClient.Get(adress)
	if err != nil {
		return nil, fmt.Errorf("unable to call public API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read public response: %w", err)
	}
	var rulesandLogic RulesandLogics
	err = json.Unmarshal(body, &rulesandLogic)
	if err != nil {
		fmt.Println(err)
	}
	return &rulesandLogic, nil
}

func (client *APIClient) PostWithJSONBody(endpoint string, jsonBody []byte) ([]byte, error) {
	ttclient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated client: %w", err)
	}

	address := fmt.Sprintf("%s/Admin/%s", client.BaseURL, endpoint)

	resp, err := ttclient.Post(address, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
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

func (r *TblTool) PrintTable() {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "ToolName", "CurrentTemplateId", "TopicName"})

	row := []string{
		fmt.Sprintf("%d", r.Id),
		r.ToolName,
		fmt.Sprintf("%d", r.CurrentTempateId),
		r.TopicName,
	}
	table.Append(row)

	table.Render()
}

func (r *TblTools) PrintTable() {

	// Create a table to display the response data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "ToolName", "CurrentTemplateId", "TopicName"})
	for _, req := range *r {
		row := []string{
			fmt.Sprintf("%d", req.Id),
			req.ToolName,
			fmt.Sprintf("%d", req.CurrentTempateId),
			req.TopicName,
		}
		table.Append(row)
	}

	table.Render()
}

func (r *RulesandLogics) PrintTable() {

	// Create a table to display the response data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "RuleName", "ToolId", "RuleId", "Value", "Logic"})
	for _, req := range *r {
		row := []string{
			fmt.Sprintf("%d", req.Id),
			req.RuleName,
			fmt.Sprintf("%d", req.ToolId),
			fmt.Sprintf("%d", req.RuleId),
			req.Value,
			req.Logic,
		}
		table.Append(row)
	}

	table.Render()
}

func MarshalToJSON(body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}
	return jsonBody, nil
}

// GetAllTblToolsForTeam calls GET /Tools/GetAllTblToolsForTeam?groupId={guid}
func (client *APIClient) GetAllTblToolsForTeam(groupId string) ([]TblTool, error) {
	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s/Tools/GetAllTblToolsForTeam?groupId=%s", client.BaseURL, url.QueryEscape(groupId))
	resp, err := httpClient.Get(address)
	if err != nil {
		return nil, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var tools []TblTool
	if err := json.Unmarshal(body, &tools); err != nil {
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}
	return tools, nil
}

// UserIsOwnerInTeam calls GET /Tools/UserIsOwnerInTeam?groupId={guid}&userId={upn}
func (client *APIClient) UserIsOwnerInTeam(groupId, userId string) (bool, error) {
	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return false, err
	}
	address := fmt.Sprintf("%s/Tools/UserIsOwnerInTeam?groupId=%s&userId=%s", client.BaseURL, url.QueryEscape(groupId), url.QueryEscape(userId))
	resp, err := httpClient.Get(address)
	if err != nil {
		return false, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("unable to read response: %w", err)
	}

	var isOwner bool
	if err := json.Unmarshal(body, &isOwner); err != nil {
		val := strings.ToLower(strings.TrimSpace(string(body)))
		return val == "true", nil
	}
	return isOwner, nil
}

// AddRequestForTool calls POST /Tools/AddRequestForTool
func (client *APIClient) AddRequestForTool(toolReq *TblToolRequest) (*TblToolRequest, error) {
	jsonBody, err := json.Marshal(toolReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s/Tools/AddRequestForTool", client.BaseURL)
	resp, err := httpClient.Post(address, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var result TblToolRequest
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}
	return &result, nil
}

// GetRequestById calls GET /Tools/GetRequestById?id={int}
func (client *APIClient) GetRequestById(id int) (*TblToolRequest, error) {
	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s/Tools/GetRequestById?id=%d", client.BaseURL, id)
	resp, err := httpClient.Get(address)
	if err != nil {
		return nil, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var toolReq TblToolRequest
	if err := json.Unmarshal(body, &toolReq); err != nil {
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}
	return &toolReq, nil
}

// UpdateToolRequestStatus calls POST /Tools/UpdateToolRequestStatus
func (client *APIClient) UpdateToolRequestStatus(id int, status string) (*MessageResponse, error) {
	bodyData := map[string]interface{}{
		"id":     id,
		"status": status,
	}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s/Tools/UpdateToolRequestStatus", client.BaseURL)
	resp, err := httpClient.Post(address, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var msgResp MessageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return &MessageResponse{Message: string(body)}, nil
	}
	return &msgResp, nil
}

// AddInstanceOfToolToDb calls POST /Tools/AddInstanceOfToolToDb
func (client *APIClient) AddInstanceOfToolToDb(groupId string, toolId int) (*TblToolInstance, error) {
	bodyData := map[string]interface{}{
		"groupId": groupId,
		"toolId":  toolId,
	}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s/Tools/AddInstanceOfToolToDb", client.BaseURL)
	resp, err := httpClient.Post(address, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var instance TblToolInstance
	if err := json.Unmarshal(body, &instance); err != nil {
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}
	return &instance, nil
}

// LogMessage calls POST /Tools/LogMessage
func (client *APIClient) LogMessage(subject, message, status string) (*LogEntry, error) {
	entry := &LogEntry{
		Subject: subject,
		Message: message,
		Status:  status,
	}
	jsonBody, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	httpClient, err := client.AuthProvider.GetAuthenticatedClient()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s/Tools/LogMessage", client.BaseURL)
	resp, err := httpClient.Post(address, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to call Tools API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	var result LogEntry
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}
	return &result, nil
}
