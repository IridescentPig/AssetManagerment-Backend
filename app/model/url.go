package model

type Url struct {
	ID              uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name            string `gorm:"column:name;not null" json:"name"`
	Url             string `gorm:"column:url" json:"url"`
	EntityID        uint   `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity          Entity `gorm:"foreignKey:EntityID;references:ID;default:null" json:"-"`
	EntitySuper     bool   `gorm:"column:entity_super;default:false" json:"entity_super"`
	DepartmentSuper bool   `gorm:"column:department_super;default:false" json:"department_super"`
	SystemSuper     bool   `gorm:"column:system_super;default:false" json:"system_super"`
}
