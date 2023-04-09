package model

type AssetClass struct {
	ID       uint        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name     string      `gorm:"column:name;not null" json:"name"`
	ParentID uint        `gorm:"column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent   *AssetClass `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	EntityID uint        `gorm:"column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity   Entity      `gorm:"foreignKey:EntityID;references:ID;default:null" json:"entity"`
	Type     int         `gorm:"column:type" json:"type"` // 1-条目型资产 2-数量型资产
}
