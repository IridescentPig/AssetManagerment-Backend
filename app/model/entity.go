package model

type Entity struct {
	ID   uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}
