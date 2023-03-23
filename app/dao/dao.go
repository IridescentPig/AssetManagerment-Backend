package dao

import (
	"database/sql"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initial() {
	var err error
	db, err = gorm.Open(mysql.Open("root:management123@tcp(localhost:3306)/management"), nil)
	if err != nil {
		log.Fatal(err)
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	if !db.Migrator().HasTable(&User{}) {
		log.Fatal("database error")
	}
}
