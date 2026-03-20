package dto

type BindTokenRequest struct {
	AgentId   string `json:"agent_id"`
	TokenName string `json:"token_name"`
}

type TokenRequest struct {
	Token string `json:"token"`
}
