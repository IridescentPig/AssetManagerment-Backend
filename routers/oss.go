package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type ossRouter struct {
}

var OssRouter *ossRouter

func newOssRouter() *ossRouter {
	return &ossRouter{}
}

func init() {
	OssRouter = newOssRouter()
}

func (oss *ossRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	group.GET("/oss/key", utils.Handler(api.OssApi.GetTempKey))
}
