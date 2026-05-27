package dto

type BuyPackageRequest struct {
	PkgID string `json:"pkg_id"`
}

type BalanceRequest struct {
	AgentID string  `json:"agent_id"`
	Amount  float64 `json:"amount"`
}

type CreatePackageRequest struct {
	Name         string  `json:"name"`
	Service      string  `json:"service"`       // "SMS" или "Call"
	Units        int64   `json:"units"`         // количество единиц (SMS/минут)
	Price        float64 `json:"price"`         // цена в рублях
	DurationDays int     `json:"duration_days"` // срок действия (дней), 0 = 30
	Description  string  `json:"description"`
}
