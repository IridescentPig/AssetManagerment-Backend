package model

import "github.com/shopspring/decimal"

type Stat struct {
	ID           uint            `gorm:"primaryKey;column:id;AUTO_INCREMENT"`
	DepartmentID uint            `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Department   Department      `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"-"`
	Total        decimal.Decimal `gorm:"type:decimal(60,2);column:total" json:"total"`
	Time         ModelTime       `gorm:"time" json:"time"`
}
