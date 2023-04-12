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
	group.POST("/login", utils.Handler(api.UserApi.UserLogin))

}

func (user *userRouter) routerNeedLogin(group *gin.RouterGroup) {
	group.POST("/logout", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.UserLogout))
	group.POST("", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.UserCreate))
	group.PATCH("/:username", utils.Handler(middleware.JWTMiddleware()), utils.Handler((api.UserApi.ResetContent)))
	group.GET("/:username/lock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.LockUser))
	group.GET("/:username/unlock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.UnlockUser))
	group.GET("/info/:user_id", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.GetUserInfoByID))
	group.GET("/list", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.GetAllUsers))
	group.DELETE("/:user_id", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.DeleteUser))
	group.POST("/info/:user_id/password", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.ChangePassword))
	group.DELETE("/info/:user_id/entity", utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.CheckSystemSuper()), utils.Handler(api.UserApi.ChangeUserEntity))
	group.DELETE("/info/:user_id/department", utils.Handler(middleware.JWTMiddleware()), utils.Handler(api.UserApi.ChangeUserDepartment))
}
