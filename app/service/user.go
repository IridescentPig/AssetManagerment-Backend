package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
)

type userService struct{}

var UserService *userService

func newUserService() *userService {
	return &userService{}
}

func init() {
	UserService = newUserService()
}

func (user *userService) CreateUser(username, password string) error {
	return dao.UserDao.Create(model.User{
		UserName: username,
		Password: password,
		Ban:      false,
	})
}

func (user *userService) VerifyPasswordAndGetUser(username, password string) (*model.User, error) {
	this_user, err := dao.UserDao.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	if this_user == nil || this_user.Password != password {
		return nil, nil
	}
	return this_user, nil
}

func (user *userService) ExistsUser(username string) (bool, error) {
	this_user, err := dao.UserDao.GetUserByName(username)
	if err != nil || this_user == nil {
		return false, err
	}
	return true, nil
}
