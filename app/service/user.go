package service

import (
	"backend/app/dao"
	"backend/app/define"
)

type userService struct{}

var UserService *userService

func (user *userService) Create(req *define.UserRegisterReq) {
	dao.UserDao.Create(req.UserName, req.Password)
}
