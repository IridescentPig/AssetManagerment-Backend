package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
	"errors"
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
	thisUser, err := dao.UserDao.GetUserByName(username)
	if err != nil {
		return "", nil, err
	}
	if thisUser == nil || thisUser.Password != password {
		return "", nil, nil
	}
	userInfo := define.UserBasicInfo{
		UserID:          thisUser.ID,
		UserName:        thisUser.UserName,
		EntitySuper:     thisUser.EntitySuper,
		DepartmentSuper: thisUser.DepartmentSuper,
		SystemSuper:     thisUser.SystemSuper,
	}
	token, err := utils.CreateToken(userInfo)
	if err != nil {
		return "", nil, err
	}
	return token, thisUser, nil
}

func (user *userService) ExistsUser(username string) (bool, error) {
	thisUser, err := dao.UserDao.GetUserByName(username)
	if err != nil || thisUser == nil {
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

func (user *userService) UserName(ctx *utils.Context) (string, error) {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			return userInfo.UserName, nil
		}
	}
	return "", errors.New("no user vertification info")
}

func (user *userService) ModifyUserIdentity(username string, identity int) error {
	return dao.UserDao.ModifyUserIdentity(username, identity)
}

func (user *userService) ModifyUserPassword(username string, password string) error {
	password = utils.CreateMD5(password)
	return dao.UserDao.ModifyUserPassword(username, password)
}

func (user *userService) ModifyUserBanstate(username string, ban bool) error {
	return dao.UserDao.ModifyUserBanstate(username, ban)
}
