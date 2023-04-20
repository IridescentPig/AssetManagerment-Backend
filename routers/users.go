package routers

import (
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
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	users.routerCheckAtHandler(group)
}

func (users *usersRouter) routerCheckAtHandler(group *gin.RouterGroup) {

}
