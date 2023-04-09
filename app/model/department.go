package model

type Department struct {
	ID       uint        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name     string      `gorm:"column:name;unique;not null" json:"name"`
	ParentID uint        `gorm:"default:null;column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent   *Department `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	EntityID uint        `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity   Entity      `gorm:"foreignKey:EntityID;references:ID;default:null" json:"entity"`
}

func (Department) TableName() string {
	return "department"
}
