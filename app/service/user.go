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

func (user *userService) CreateUser(username, password string) error {
	password = utils.CreateMD5(password)
	return dao.UserDao.Create(model.User{
		UserName:     username,
		Password:     password,
		EntityID:     nil,
		DepartmentID: nil,
		Ban:          false,
	})
}

func (user *userService) VerifyPasswordAndGetUser(username, password string) (string, *model.User, error) {
	password = utils.CreateMD5(password)
	this_user, err := dao.UserDao.GetUserByName(username)
	if err != nil {
		return "", nil, err
	}
	if this_user == nil || this_user.Password != password {
		return "", nil, nil
	}
	userInfo := define.UserBasicInfo{
		UserID:          this_user.ID,
		UserName:        this_user.UserName,
		EntitySuper:     this_user.EntitySuper,
		DepartmentSuper: this_user.DepartmentSuper,
		SystemSuper:     this_user.SystemSuper,
	}
	token, err := utils.CreateToken(userInfo)
	if err != nil {
		return "", nil, err
	}
	return token, this_user, nil
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
