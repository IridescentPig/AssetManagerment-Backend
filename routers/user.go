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
	user.routerNotNeedLogin(group.Group(""))
	user.routerNeedLogin(group.Group(""))
}

func (user *userRouter) routerNotNeedLogin(group *gin.RouterGroup) {
	group.POST("/feishu/callback", utils.Handler(api.FeishuApi.FeishuCallBack))
	group.POST("/feishu/login", utils.Handler(api.FeishuApi.FeishuLogin))
	group.Use(utils.Handler(middleware.LogMiddleware()))
	group.POST("/register", utils.Handler(api.UserApi.UserRegister))
	group.POST("/login", utils.Handler(api.UserApi.UserLogin))
}

func (user *userRouter) routerNeedLogin(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	group.POST("/logout", utils.Handler(api.UserApi.UserLogout))
	group.POST("", utils.Handler(api.UserApi.UserCreate))
	group.PATCH("/:username", utils.Handler((api.UserApi.ResetContent)))
	group.GET("/:username/lock", utils.Handler(api.UserApi.LockUser))
	group.GET("/:username/unlock", utils.Handler(api.UserApi.UnlockUser))
	group.GET("/info/:user_id", utils.Handler(api.UserApi.GetUserInfoByID))
	group.GET("/list", utils.Handler(api.UserApi.GetAllUsers))
	group.DELETE("/:user_id", utils.Handler(api.UserApi.DeleteUser))
	group.POST("/info/:user_id/password", utils.Handler(api.UserApi.ChangePassword))
	group.PATCH("/info/:user_id/identity", utils.Handler(api.UserApi.ModifyUserIdentity))
	group.DELETE("/info/:user_id/entity", utils.Handler(middleware.CheckSystemSuper()), utils.Handler(api.UserApi.ChangeUserEntity))
	group.DELETE("/info/:user_id/department", utils.Handler(api.UserApi.ChangeUserDepartment))
	group.POST("/info/:user_id/entity", utils.Handler(api.UserApi.ChangeUserEntity))
	group.POST("/info/:user_id/department", utils.Handler(api.UserApi.ChangeUserDepartment))
	group.POST("/feishu/bind", utils.Handler(api.FeishuApi.FeishuBind))
	group.DELETE("/feishu/bind", utils.Handler(api.FeishuApi.FeishuUnBind))
}
