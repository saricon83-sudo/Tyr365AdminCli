package teamToolboxHelper

type APIClient struct {
	AuthProvider *TokenProvider
	BaseURL      string
}

var tokenProvider *TokenProvider

type TblTool struct {
	Id               int32  `json:"Id"`
	ToolName         string `json:"toolName"`
	CurrentTempateId int32  `json:"currentTemplateId"`
	TopicName        string `json:"topicName"`
}

type TblTools []TblTool

type ToolStatusMessage struct {
	Id     int32
	Status string
}

type RulesandLogic struct {
	Id       int32  `json:"Id"`
	RuleName string `json:"ruleName"`
	ToolId   int32  `json:"toolId"`
	RuleId   int32  `json:"ruleId"`
	Value    string `json:"value"`
	Logic    string `json:"logic"`
}

type RulesandLogics []RulesandLogic
