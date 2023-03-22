package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
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

	if err := context.ShouldBindWith(&req, binding.Form); err != nil {
		utils.NewResponseJson(context).Error(http.StatusBadRequest, -1, "Invalid request body.", nil)
		return
	}
	service.UserService.Create(&req)

	utils.NewResponseJson(context).Success("Successfully register.", nil)
}

func (user *userApi) UserLogin(context *gin.Context) {
	var req define.UserLoginReq

	if err := context.ShouldBindWith(&req, binding.Form); err != nil {
		utils.NewResponseJson(context).Error(http.StatusBadRequest, -1, "Invalid request body.", nil)
		return
	}

	utils.NewResponseJson(context).Success("Successfully login.", nil)
}
