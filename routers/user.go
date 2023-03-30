package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
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
	user.routerNeedLogin(group)
}

func (user *userRouter) routerNotNeedLogin(group *gin.RouterGroup) {
	group.POST("/register", utils.Handler(api.UserApi.UserRegister))
	// group.OPTIONS("/register", utils.Handler(optionsHandler))
	group.POST("/login", utils.Handler(api.UserApi.UserLogin))
	// group.OPTIONS("/login", utils.Handler(optionsHandler))
}
func (user *userRouter) routerNeedLogin(group *gin.RouterGroup) {
	group.POST("/logout", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.UserLogout))
	// group.OPTIONS("/logout", utils.Handler(optionsHandler))
	group.POST("", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.UserCreate))
	// group.OPTIONS("", utils.Handler(optionsHandler))
	// TODO:
	group.PATCH("/:username", utils.Handler(middleware.JWTMiddleware()), utils.Handler((api.UserApi.ResetContent)))
	// group.OPTIONS("/:username", utils.Handler(optionsHandler))
	group.GET("/:username/lock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.LockUser))
	// group.OPTIONS("/:username/lock", utils.Handler(optionsHandler))
	group.GET("/:username/unlock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.UnlockUser))
	// group.OPTIONS("/:username/unlock", utils.Handler(optionsHandler))
}

// func optionsHandler(ctx *utils.Context) {
// 	ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
// 	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
// }
