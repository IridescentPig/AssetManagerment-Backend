package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type asyncRouter struct {
}

var AsyncRouter *asyncRouter

func newAsyncRouter() *asyncRouter {
	return &asyncRouter{}
}

func init() {
	AsyncRouter = newAsyncRouter()
}

func (asy *asyncRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	group.GET("/users/:user_id/async/list", utils.Handler(api.AsyncApi.GetUserAsyncTasks))
	group.POST("/users/:user_id/async", utils.Handler(api.AsyncApi.CreateAsyncTask))
	group.PATCH("/users/:user_id/async/:task_id", utils.Handler(api.AsyncApi.ModifyAsyncState))
}

// func (asy *asyncRouter) routerCheckAtHandler(group *gin.RouterGroup) {
// 	group.GET("/user/:user_id/async/list", utils.Handler(api.AsyncApi.GetUserAsyncTasks))
// 	group.POST("/user/:user_id/async", utils.Handler(api.AsyncApi.CreateAsyncTask))
// 	group.PATCH("/user/:user_id/async/:task_id", utils.Handler(api.AsyncApi.ModifyAsyncState))
// }
