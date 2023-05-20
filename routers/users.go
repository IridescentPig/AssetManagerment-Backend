package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type usersRouter struct{}

var UsersRouter *usersRouter

func newUsersRouter() *usersRouter {
	return &usersRouter{}
}

func init() {
	UsersRouter = newUsersRouter()
}

func (users *usersRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	users.routerCheckAtHandler(group)
}

func (users *usersRouter) routerCheckAtHandler(group *gin.RouterGroup) {
	group.GET("/:user_id/assets", utils.Handler(api.UserApi.GetUserUsedAssets))
}
