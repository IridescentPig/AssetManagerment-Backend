package model

import (
	"github.com/shopspring/decimal"
)

type Asset struct {
	ID          uint            `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name        string          `gorm:"column:name;not null" json:"username"`
	ParentID    uint            `gorm:"column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent      *Asset          `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	UserID      uint            `gorm:"column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	User        User            `gorm:"foreignKey:EntityID;references:ID;default:null" json:"entity"`
	Price       decimal.Decimal `gorm:"type:decimal(8,2);column:price" json:"price"`
	Description string          `gorm:"column:description" json:"description"`
	Position    string          `gorm:"column:position" json:"position"`
	Expire      bool            `gorm:"column:expire" json:"expire"`
}

func (Asset) TableName() string {
	return "asset"
}
