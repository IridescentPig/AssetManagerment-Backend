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
	// 是否需要加密传输？
	// 是否需要通过中间件处理？
	var req define.UserLoginReq
	var this_user *model.User

	if err := ctx.MustBindWith(&req, binding.Form); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	var err error
	this_user, err = service.UserService.VerifyPasswordAndGetUser(req.UserName, req.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if this_user == nil {
		ctx.BadRequest(1, "Wrong Username Or Password")
		return
	} else if this_user.Ban {
		ctx.BadRequest(3, "User Banned")
	}
	//todo
	ctx.Success(nil)
}

func (user *userApi) UserLogout(ctx *utils.Context) {
	// 使用中间件验证 token 是否正确
	ctx.Success(nil)
}
