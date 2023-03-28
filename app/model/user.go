package model

type User struct {
	ID              uint       `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	UserName        string     `gorm:"column:username;unique;not null" json:"username"`
	Password        string     `gorm:"column:password;not null" json:"password"`
	EntityID        uint       `gorm:"column:entity_id" json:"entity_id"`
	Entity          Entity     `gorm:"foreignKey:EntityID;references:ID"`
	EntitySuper     bool       `gorm:"column:entity_super;default:false"`
	DepartmentID    uint       `gorm:"column:department_id" json:"department_id"`
	Department      Department `gorm:"foreignKey:DepartmentID;references:ID"`
	DepartmentSuper bool       `gorm:"column:department_super;default:false"`
	SystemSuper     bool       `gorm:"column:system_super;default:false"`
	Ban             bool       `gorm:"column:ban;default:false"`
}
