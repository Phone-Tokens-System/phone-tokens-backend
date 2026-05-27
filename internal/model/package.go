package model

import "github.com/google/uuid"

type Package struct {
	Id           uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string      `json:"name" gorm:"type:varchar(255); unique"`
	Service      ServiceType `json:"service" gorm:"type:varchar(50)"`
	Units        int64       `json:"units" gorm:"type:int; not null"`
	Price        float64     `json:"price" gorm:"type:float; not null"`
	DurationDays int         `json:"duration_days" gorm:"type:int; not null; default:30"` // срок действия в днях
	Description  string      `json:"description" gorm:"type:text"`
}

func (Package) TableName() string {
	return "packages"
}
