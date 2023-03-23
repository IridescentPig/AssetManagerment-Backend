package dao

import (
	"log"
)

type User struct {
	ID       uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	UserName string `gorm:"column:username;unique;not null" json:"username"`
	Password string `gorm:"column:password;not null" json:"password"`
}

type userDao struct {
}

var UserDao *userDao

func newUserDao() *userDao {
	return &userDao{}
}

func init() {
	UserDao = newUserDao()
}

func (user *userDao) Create(username string, password string) {
	newUser := User{UserName: username, Password: password}
	err := db.Create(&newUser)
	if err != nil {
		log.Println(err)
	}
}
