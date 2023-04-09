package model

type Entity struct {
	ID          uint      `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"entity_id"`
	Name        string    `gorm:"column:name;unique;not null" json:"entity_name"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   ModelTime `gorm:"column:created_at" json:"created_at"`
}

func (Entity) TableName() string {
	return "entity"
}
