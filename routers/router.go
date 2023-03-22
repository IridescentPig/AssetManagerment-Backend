package routers

import (
	"asset-management/utils"
	"net/http"

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
	r.NoRoute(routeNotFound)
	r.NoMethod(methodNotFound)

	UserRouter.Init(r.Group("/user"))

	return r
}

func routeNotFound(context *gin.Context) {
	utils.NewResponseJson(context).Error(http.StatusNotFound, 1, "Router not found.", nil)
}

func methodNotFound(context *gin.Context) {
	utils.NewResponseJson(context).Error(http.StatusNotFound, 1, "Method not found.", nil)
}
