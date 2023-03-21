package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type userApi struct {
}

var UserApi *userApi

func (user *userApi) UserRegister(context *gin.Context) {
	var req define.UserRegisterReq

	if err := context.BindJSON(&req); err != nil {
		return
	}
	service.UserService.Create(&req)

	utils.NewResponseJson(context).Success("Successfully register.", nil)
}
