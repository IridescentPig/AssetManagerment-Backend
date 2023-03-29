package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
)

type userService struct{}

var UserService *userService

func newUserService() *userService {
	return &userService{}
}

func init() {
	UserService = newUserService()
}

// 为保证注册功能能用，暂时留着，后面必须改API，否则不传用户类型怎么设置用户权限？
func (user *userService) CreateUser(username string, password string) error {
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

func (user *userService) SystemSuper(ctx *utils.Context) bool {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			if userInfo.SystemSuper {
				return true
			}
		}
	}
	return false
}

func (user *userService) EntitySuper(ctx *utils.Context) bool {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			if userInfo.EntitySuper || userInfo.SystemSuper {
				return true
			}
		}
	}
	return false
}

func (user *userService) DepartmentSuper(ctx *utils.Context) bool {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			if userInfo.DepartmentSuper ||
				userInfo.EntitySuper ||
				userInfo.SystemSuper {
				return true
			}
		}
	}
	return false
}
