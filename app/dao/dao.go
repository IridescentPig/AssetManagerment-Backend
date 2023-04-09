package dao

import (
	"asset-management/app/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func connect() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Entity{}, &model.Department{}, &model.User{})

	//tables := make([]string, 0)
	//db.Raw("SELECT name FROM sqlite_master WHERE type='table' order by name").Scan(&tables)
	//fmt.Println("table is", tables)

	if !db.Migrator().HasTable(&model.User{}) {
		log.Fatal("database error")
	}
}

func Initial() {
	var err error

	if gin.Mode() == gin.DebugMode {
		if os.Getenv("DEBUG") == "" {
			db, err = gorm.Open(mysql.Open("manager:BinaryAbstract@tcp(49.233.51.221:25000)/asset"), nil)
		} else {
			db, err = gorm.Open(mysql.Open("manager:BinaryAbstract@tcp(AssetManagement-Database-dev.BinaryAbstract.secoder.local:80)/asset"), nil)
		}
	} else {
		db, err = gorm.Open(mysql.Open("manager:BinaryAbstract@tcp(AssetManagement-Database.BinaryAbstract.secoder.local:3306)/asset"), nil)
	}

	if err != nil {
		log.Fatal(err)
	}
	connect()
}

func InitForTest() {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		log.Fatal(err)
	}

	connect()
}
