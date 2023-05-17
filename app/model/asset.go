package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Asset struct {
	ID           uint                        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name         string                      `gorm:"column:name;not null" json:"name"`
	ParentID     uint                        `gorm:"default:null;column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent       *Asset                      `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	UserID       uint                        `gorm:"default:null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	User         User                        `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	DepartmentID uint                        `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department   Department                  `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	Price        decimal.Decimal             `gorm:"type:decimal(10,2);column:price" json:"price"`
	Description  string                      `gorm:"column:description" json:"description"`
	Position     string                      `gorm:"column:position" json:"position"`
	Expire       uint                        `gorm:"column:expire;default:0" json:"expire"`
	ClassID      uint                        `gorm:"default:null;column:class_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"class_id"`
	Class        AssetClass                  `gorm:"foreignKey:ClassID;references:ID;default:null" json:"class"`
	Number       int                         `gorm:"column:number" json:"count"`
	Type         int                         `gorm:"column:type" json:"type"`   // 1-条目型资产 2-数量型资产
	State        uint                        `gorm:"column:state" json:"state"` // 0idle;1in_use;2in_maintain;3retired;4deleted
	MaintainerID uint                        `gorm:"default:null;column:maintainer_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"maintainer_id"`
	Maintainer   User                        `gorm:"foreignKey:MaintainerID;references:ID;default:null" json:"maintainer"`
	Property     datatypes.JSON              `gorm:"column:property;" json:"property"`
	TaskList     []*Task                     `gorm:"many2many:task_assets;" json:"task_list"`
	CreatedAt    *ModelTime                  `gorm:"column:created_at" json:"created_at"`
	NetWorth     decimal.Decimal             `gorm:"type:decimal(10,2);column:net_worth" json:"net_worth"`
	ImgList      datatypes.JSONSlice[string] `gorm:"column:img_list" json:"img_list"`
	Warn         bool                        `gorm:"default:false;colimn:warn" json:"warn"`
	Threshold    uint                        `gorm:"default:0;column:threshold" json:"threshold"`
}
