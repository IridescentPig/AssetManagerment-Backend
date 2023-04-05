package model

type Department struct {
	ID       uint        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name     string      `gorm:"column:name;unique;not null" json:"name"`
	ParentID *uint       `gorm:"column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent   *Department `gorm:"foreignKey:EntityID;references:ID;default:null" json:"parent"`
	EntityID *uint       `gorm:"column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity   *Entity     `gorm:"foreignKey:EntityID;references:ID;default:null" json:"entity"`
	//UserList       []*User       `gorm:"column:user_list;foreignKey:DepartmentID" json:"users"`
	//DepartmentList []*Department `gorm:"column:department_list;foreignKey:ID" json:"departments"`
}

func (Department) TableName() string {
	return "department"
}
