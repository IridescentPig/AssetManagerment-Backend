package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
)

type userApi struct {
}

var UserApi *userApi

func newUserApi() *userApi {
	return &userApi{}
}

func init() {
	UserApi = newUserApi()
}

func (user *userApi) UserRegister(ctx *utils.Context) {
	var req define.UserRegisterReq

	if err := ctx.MustBindWith(&req, binding.Form); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	exists, err := service.UserService.ExistsUser(req.UserName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if exists {
		ctx.BadRequest(1, "Duplicated Name")
		return
	}
	err = service.UserService.CreateUser(req.UserName, req.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

func (user *userApi) UserLogin(ctx *utils.Context) {
	var req define.UserLoginReq
	var this_user *model.User

	if err := ctx.MustBindWith(&req, binding.Form); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	var err error
	var token string
	token, this_user, err = service.UserService.VerifyPasswordAndGetUser(req.UserName, req.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if this_user == nil {
		ctx.BadRequest(1, "Wrong Username Or Password")
		return
	} else if this_user.Ban {
		ctx.BadRequest(3, "User Banned")
	}

	data := struct {
		Token string      `json:"token"`
		User  *model.User `json:"user"`
	}{
		Token: token,
		User:  this_user,
	}

	ctx.Success(data)
}

func (user *userApi) Logout(ctx *utils.Context) {
	//clear session
}
