package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
	"errors"

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

func (user *userDao) UpdateByName(username string, data map[string]interface{}) error {
	result := db.Model(&model.User{}).Where("username = ?", username).Updates(data)
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

func (user *userDao) ModifyUserIdentity(username string, identity int) error {
	this_user, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	if identity == 0 {
		err = user.Update(this_user.ID, map[string]interface{}{
			"system_super":     false,
			"entity_super":     false,
			"department_super": false,
		})
	} else if identity == 1 {
		err = user.Update(this_user.ID, map[string]interface{}{
			"system_super":     false,
			"entity_super":     false,
			"department_super": true,
		})
	} else if identity == 2 {
		err = user.Update(this_user.ID, map[string]interface{}{
			"system_super":     false,
			"entity_super":     true,
			"department_super": false,
		})
	} else if identity == 3 {
		err = user.Update(this_user.ID, map[string]interface{}{
			"system_super":     true,
			"entity_super":     false,
			"department_super": false,
		})
	} else {
		err = errors.New("invalid identity number")
	}
	return err
}

func (user *userDao) ModifyUserPassword(username string, password string) error {
	this_user, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	err = user.Update(this_user.ID, map[string]interface{}{
		"password": password,
	})
	return err
}
