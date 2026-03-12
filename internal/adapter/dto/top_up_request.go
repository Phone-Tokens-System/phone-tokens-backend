package dto

type TopUpRequest struct {
	AgentID string  `json:"agent_id"`
	Amount  float64 `json:"amount"`
}
