package routers

import (
	"asset-management/app/api"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type entityRouter struct{}

var EntityRouter *entityRouter

func newEntityRouter() *entityRouter {
	return &entityRouter{}
}

func init() {
	EntityRouter = newEntityRouter()
}

func (entity *entityRouter) Init(group *gin.RouterGroup) {
	group.GET("/list", utils.Handler(api.EntityApi.GetEntityList))
	group.GET("/user/list", utils.Handler(api.EntityApi.UsersInEntity))
}
