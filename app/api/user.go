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

func (user *userApi) UserLogout(ctx *utils.Context) {
	// 使用中间件验证 token 是否正确
	ctx.Success(nil)
}

func (user *userApi) UserCreate(ctx *utils.Context) {
	// TODO: 暂时使用 register 的 req
	var req define.UserRegisterReq
	if err := ctx.MustBindWith(&req, binding.JSON); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	// 暂时按照只有超级用户可以创建来处理
	if !service.UserService.SystemSuper(ctx) {
		ctx.Forbidden(2, "Permission Denied.")
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

func (user *userApi) ResetContent(ctx *utils.Context) {
	var info define.UriInfo
	var req define.ResetReq
	errInfo := ctx.ShouldBindUri(&info)
	errReq := ctx.MustBindWith(&req, binding.JSON)
	if errInfo != nil || errReq != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	//查找用户是否存在
	exists, err := service.UserService.ExistsUser(info.UserName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(1, "User Not Found")
		return
	}
	if req.Method == 0 {
		// 修改身份
		// 暂时按照只有超级用户可以创建来处理
		if !service.UserService.SystemSuper(ctx) {
			ctx.Forbidden(2, "Permission Denied.")
			return
		} else {
			err = service.UserService.ModifyUserIdentity(info.UserName, req.Identity)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			}
		}
	} else if req.Method == 1 {
		// 修改密码
		// 超级用户和自己都应该可以修改密码
		// 自己修改密码需要验证是否为本人
		if !service.UserService.SystemSuper(ctx) {
			username, err := service.UserService.UserName(ctx)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			}
			if username != info.UserName {
				ctx.Forbidden(2, "Permission Denied.")
			}
		}
		err = service.UserService.ModifyUserPassword(info.UserName, req.Password)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	} else {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	ctx.Success(nil)
}
