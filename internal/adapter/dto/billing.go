package dto

type BuyPackageRequest struct {
	PkgID string `json:"pkg_id"`
}

type BalanceRequest struct {
	AgentID string  `json:"agent_id"`
	Amount  float64 `json:"amount"`
}
