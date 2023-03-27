package dao

import (
	"asset-management/utils"
)

type User struct {
	ID               uint       `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	UserName         string     `gorm:"column:username;unique;not null" json:"username"`
	Password         string     `gorm:"column:password;not null" json:"password"`
	EntityID         uint       `gorm:"column:entity_id" json:"entity_id"`
	Entity           Entity     `gorm:"foreignKey:EntityID;references:ID"`
	Entity_super     bool       `gorm:"column:entity_super;default:false"`
	DepartmentID     uint       `gorm:"column:department_id" json:"department_id"`
	Department       Department `gorm:"foreignKey:DepartmentID;references:ID"`
	Department_super bool       `gorm:"column:department_super;default:false"`
	System_super     bool       `gorm:"column:system_super;default:false"`
	HasLogin         bool       `gorm:"column:has_login;default:false"`
	Ban              bool       `gorm:"column:ban;default:false"`
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

func (user *userDao) Create(username string, password string) error {
	newUser := User{UserName: username, Password: password}
	result := db.Model(&User{}).Create(&newUser)
	return utils.DB_error(result)
}

func (user *userDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&User{}).Where("id = ?", id).Updates(data)
	return utils.DB_error(result)
}

func (user *userDao) Delete(id []uint) error {
	result := db.Model(&User{}).Where("id in (?)", id).Delete(&User{})
	return utils.DB_error(result)
}

func (user *userDao) AllUser() (list []User, err error) {
	result := db.Model(&User{}).Find(&list)
	err = utils.DB_error(result)
	return
}

func (user *userDao) AllUserWhere(query interface{}, args ...interface{}) (list []User, err error) {
	result := db.Model(&User{}).Where(query, args...).Find(&list)
	err = utils.DB_error(result)
	return
}

func (user *userDao) OneUserWhere(query interface{}, args ...interface{}) (record User, err error) {
	result := db.Model(&User{}).Where(query, args...).First(&record)
	err = utils.DB_error(result)
	return
}

func (user *userDao) UserCount() (count int64, err error) {
	result := db.Model(&User{}).Count(&count)
	err = utils.DB_error(result)
	return
}

func (user *userDao) UserCountWhere(query interface{}, args ...interface{}) (count int64, err error) {
	result := db.Model(&User{}).Where(query, args...).Count(&count)
	err = utils.DB_error(result)
	return
}
