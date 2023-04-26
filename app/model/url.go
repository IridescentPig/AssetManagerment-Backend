package model

type Url struct {
	Name     string `gorm:"column:name;unique;not null" json:"name"`
	Url      string `gorm:"column:url" json:"url"`
	EntityID uint   `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity   Entity `gorm:"foreignKey:EntityID;references:ID;default:null" json:"-"`
}
