package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Student struct {
	ID   uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name;unique;not null" json:"name"`
}

func TestError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Student{})
	result := db.Model(&Student{}).Create(&Student{
		Name: "test",
	})
	err = DBError(result)
	assert.Equal(t, nil, err)

	result = db.Model(&Student{}).Create(&Student{
		Name: "test",
	})
	err = DBError(result)
	assert.Equal(t, false, err == nil)

	code := 1
	info := "Invalid request body"
	code, info = ServiceError(code, info)
	assert.Equal(t, 1, code)
	assert.Equal(t, "Invalid request body", info)
}
