package routers

import (
	"asset-management/app/api"
	"asset-management/utils"

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
	group.POST("/register", utils.Handler(api.UserApi.UserRegister))
	group.POST("/login", utils.Handler(api.UserApi.UserLogin))
}
