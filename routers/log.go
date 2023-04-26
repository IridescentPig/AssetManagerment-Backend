package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type logRouter struct{}

var LogRouter *logRouter

func newLogRouter() *logRouter {
	return &logRouter{}
}

func init() {
	LogRouter = newLogRouter()
}

func (mylog *logRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	group.GET("/:entity_id/login-logs", utils.Handler(api.LogApi.GetLoginLog))
	group.GET("/:entity_id/data-logs", utils.Handler(api.LogApi.GetDataLog))
}
