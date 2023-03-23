package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
)

type userService struct{}

var UserService *userService

func newUserService() *userService {
	return &userService{}
}

func init() {
	UserService = newUserService()
}

func (user *userService) Create(req *define.UserRegisterReq) {
	dao.UserDao.Create(req.UserName, req.Password)
}
