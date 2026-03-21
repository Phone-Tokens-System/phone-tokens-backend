package model

type CsrRequest struct {
	ID          int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AgentID     string `json:"agent_id;omitempty" gorm:"column:agent_id;type:uuid"`
	ServiceName string `json:"service_name" gorm:"NOT NULL;"`
	Email       string `json:"email"`
	CSR         []byte `json:"csr" gorm:"size:65535"`
	Status      string `json:"status" gorm:"size:255"`
}

func (CsrRequest) TableName() string {
	return "certificate_requests"
}
