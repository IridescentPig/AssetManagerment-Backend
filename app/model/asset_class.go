package model

type AssetClass struct {
	ID           uint        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name         string      `gorm:"column:name;not null" json:"name"`
	ParentID     uint        `gorm:"column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent       *AssetClass `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	DepartmentID uint        `gorm:"column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department   Department  `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	Type         int         `gorm:"column:type" json:"type"` // 1-条目型资产 2-数量型资产
}

func (AssetClass) TableName() string {
	return "assetclass"
}
