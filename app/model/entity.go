package model

type Entity struct {
	ID   uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	//UserList []*User `gorm:"column:user_list;foreignKey:EntityID" json:"users"`
}

func (Entity) TableName() string {
	return "entity"
}
