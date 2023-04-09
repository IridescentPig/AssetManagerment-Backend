package model

type Entity struct {
	ID   uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"entity_id"`
	Name string `gorm:"column:name;unique;not null" json:"entity_name"`
}

func (Entity) TableName() string {
	return "entity"
}
