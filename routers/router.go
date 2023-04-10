package routers

import (
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type router struct{}

var Router *router

func newRouter() *router {
	return &router{}
}

func init() {
	Router = newRouter()
}

func (router *router) Init(r *gin.Engine) *gin.Engine {
	r.NoRoute(utils.Handler(RouteNotFound))
	r.NoMethod(utils.Handler(MethodNotFound))

	UserRouter.Init(r.Group("/user"))
	EntityRouter.Init(r.Group("/entity"))

	return r
}

func RouteNotFound(ctx *utils.Context) {
	ctx.NotFound(1, "Router not found.")
}

func MethodNotFound(ctx *utils.Context) {
	ctx.NotFound(1, "Method not found.")
}
