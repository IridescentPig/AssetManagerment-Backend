package model

import (
	"github.com/shopspring/decimal"
)

type Asset struct {
	ID          uint            `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name        string          `gorm:"column:name;not null" json:"name"`
	ParentID    uint            `gorm:"column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent      *Asset          `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	UserID      uint            `gorm:"column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	User        User            `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	Price       decimal.Decimal `gorm:"type:decimal(8,2);column:price" json:"price"`
	Description string          `gorm:"column:description" json:"description"`
	Position    string          `gorm:"column:position" json:"position"`
	Expire      bool            `gorm:"column:expire;default:false" json:"expire"`
	ClassID     uint            `gorm:"column:class_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"class_id"`
	Class       AssetClass      `gorm:"foreignKey:ClassID;references:ID;default:null" json:"class"`
	Number      int             `gorm:"column:number" json:"number"`
	Type        int             `gorm:"column:type" json:"type"` // 1-条目型资产 2-数量型资产
}
