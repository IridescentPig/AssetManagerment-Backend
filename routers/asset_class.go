package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type assetClassRouter struct{}

var AssetClassRouter *assetClassRouter

func newAssetClassRouter() *assetClassRouter {
	return &assetClassRouter{}
}

func init() {
	AssetClassRouter = newAssetClassRouter()
}

func (assetClass *assetClassRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	assetClass.routerCheckAtHandler(group)
}

func (assetClass *assetClassRouter) routerCheckAtHandler(group *gin.RouterGroup) {
	group.GET("/:department_id/asset_class/tree", utils.Handler(api.AssetClassApi.GetAssetClassTree))
	group.POST("/:department_id/asset_class", utils.Handler(api.AssetClassApi.CreateAssetClass))
	group.DELETE("/:department_id/asset_class/:class_id", utils.Handler(api.AssetClassApi.DeleteAssetClass))
	group.PATCH("/:department_id/asset_class/:class_id", utils.Handler(api.AssetClassApi.ModifyAssetClassInfo))
	group.GET("/:department_id/asset_class/:class_id", utils.Handler(api.AssetClassApi.GetSubAssetClass))
}
