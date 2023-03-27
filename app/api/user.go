package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (user *userApi) UserRegister(context *gin.Context) {
	var req define.UserRegisterReq

	if err := context.MustBindWith(&req, binding.Form); err != nil {
		utils.NewResponseJson(context).Error(http.StatusBadRequest, -1, "Invalid request body.", nil)
		return
	}
	code, info := service.UserService.Register(&req)
	if code != 0 {
		utils.NewResponseJson(context).Error(http.StatusBadRequest, code, info, nil)
		return
	}

	utils.NewResponseJson(context).Success("Successfully register.", nil)
}

func (user *userApi) UserLogin(context *gin.Context) {
	var req define.UserLoginReq

	if err := context.MustBindWith(&req, binding.Form); err != nil {
		utils.NewResponseJson(context).Error(http.StatusBadRequest, -1, "Invalid request body.", nil)
		return
	}
	code, info := service.UserService.Login(&req)
	if code != 0 {
		utils.NewResponseJson(context).Error(http.StatusBadRequest, code, info, nil)
	}

	utils.NewResponseJson(context).Success("Successfully login.", nil)
}
