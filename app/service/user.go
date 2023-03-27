package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
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

func (user *userService) Register(req *define.UserRegisterReq) (code int, info string) {
	err := dao.UserDao.Create(req.UserName, req.Password)
	if err != nil {
		return utils.Service_error(1, "Duplicated Name")
	}
	return
}

func (user *userService) Login(req *define.UserLoginReq) (code int, info string) {
	this_user, err := dao.UserDao.OneUserWhere("username = ? and password = ?", req.UserName, req.Password)
	if err != nil {
		return utils.Service_error(2, "Wrong UserName Or Password")
	}
	if this_user.HasLogin {
		return utils.Service_error(1, "User Has Logged In")
	}
	if this_user.Ban {
		return utils.Service_error(3, "User was banned")
	}
	return
}
