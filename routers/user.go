package routers

import (
	"asset-management/app/api"

	"github.com/gin-gonic/gin"
)

type userRouter struct{}

var UserRouter *userRouter

func newUserRouter() *userRouter {
	return &userRouter{}
}

func init() {
	UserRouter = newUserRouter()
}

func (user *userRouter) Init(group *gin.RouterGroup) {
	user.routerNotNeedLogin(group)
}

func (user *userRouter) routerNotNeedLogin(group *gin.RouterGroup) {
	group.POST("/register", api.UserApi.UserRegister)
	group.POST("/login", api.UserApi.UserLogin)
}
