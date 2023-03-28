package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type userDao struct {
}

var UserDao *userDao

func newUserDao() *userDao {
	return &userDao{}
}

func init() {
	UserDao = newUserDao()
}

func (user *userDao) Create(newUser model.User) error {
	result := db.Model(&model.User{}).Create(&newUser)
	return utils.DB_error(result)
}

func (user *userDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.User{}).Where("id = ?", id).Updates(data)
	return utils.DB_error(result)
}

func (user *userDao) Delete(id []uint) error {
	result := db.Model(&model.User{}).Where("id in (?)", id).Delete(&model.User{})
	return utils.DB_error(result)
}

func (user *userDao) AllUser() (list []model.User, err error) {
	result := db.Model(&model.User{}).Find(&list)
	err = utils.DB_error(result)
	return
}

func (user *userDao) GetUserByName(username string) (*model.User, error) {
	ret := &model.User{}
	result := db.Model(&model.User{}).Where("username = ?", username).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DB_error(result)
}

func (user *userDao) GetUsersByNames(username []string) (list []model.User, err error) {
	result := db.Model(&model.User{}).Where("username = ?", username).Find(&list)
	err = utils.DB_error(result)
	return
}

func (user *userDao) UserCount() (count int64, err error) {
	result := db.Model(&model.User{}).Count(&count)
	err = utils.DB_error(result)
	return
}
