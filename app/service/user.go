package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
)

type userService struct{}

var UserService *userService

func (user *userService) Create(req *define.UserRegisterReq) {
	dao.UserDao.Create(req.UserName, req.Password)
}
