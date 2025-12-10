package model

type CsrRequest struct {
	ID     int    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Email  string `json:"email" gorm:"column:email"`
	CSR    []byte `json:"csr" gorm:"column:csr;size:65535"`
	Status string `json:"status" gorm:"column:status;size:255"`
}

func (CsrRequest) TableName() string {
	return "certificate_requests"
}
