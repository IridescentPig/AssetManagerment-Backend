package dao

import (
	"database/sql"
	_ "database/sql"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(mysql.Open("root@localhost:3306"), nil)
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
}

