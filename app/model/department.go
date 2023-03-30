package model

type Department struct {
	ID uint `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
}
