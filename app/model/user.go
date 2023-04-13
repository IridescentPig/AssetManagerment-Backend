package model

type User struct {
	ID              uint        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	UserName        string      `gorm:"column:username;unique;not null" json:"username"`
	Password        string      `gorm:"column:password;not null" json:"-"`
	EntityID        uint        `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity          *Entity     `gorm:"foreignKey:EntityID;references:ID;default:null" json:"entity"`
	EntitySuper     bool        `gorm:"column:entity_super;default:false" json:"entity_super"`
	DepartmentID    uint        `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department      *Department `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	DepartmentSuper bool        `gorm:"column:department_super;default:false" json:"department_super"`
	SystemSuper     bool        `gorm:"column:system_super;default:false" json:"system_super"`
	IsEmployee      bool        `gorm:"column:is_employee;default:true" json:"is_employee"`
	Ban             bool        `gorm:"column:ban;default:false" json:"-"`
}
