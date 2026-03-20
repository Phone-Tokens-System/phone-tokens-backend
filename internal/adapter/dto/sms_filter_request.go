package dto

type SmsFilterRequest struct {
	ServiceName string        `json:"service_name"`
	Certificate string        `json:"certificate"`
	AgentID     string        `json:"agent_id,omitempty" gorm:"-"`
	Text        string        `json:"text"`
	Filters     FilterRequest `json:"filters"`
}
