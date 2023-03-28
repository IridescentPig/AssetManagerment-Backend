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
	r.NoRoute(utils.Handler(routeNotFound))
	r.NoMethod(utils.Handler(methodNotFound))

	UserRouter.Init(r.Group("/user"))

	return r
}

func routeNotFound(ctx *utils.Context) {
	ctx.NotFound(1, "Router not found.")
}

func methodNotFound(ctx *utils.Context) {
	ctx.NotFound(1, "Method not found.")
}
