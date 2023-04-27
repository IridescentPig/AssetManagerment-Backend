package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type assetRouter struct{}

var AssetRouter *assetRouter

func newAssetRouter() *assetRouter {
	return &assetRouter{}
}

func init() {
	AssetRouter = newAssetRouter()
}

func (asset *assetRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	asset.routerCheckAtHandler(group)
}

func (asset *assetRouter) routerCheckAtHandler(group *gin.RouterGroup) {
	group.GET("/:department_id/asset/list", utils.Handler(api.AssetApi.GetAssetList))
	group.PATCH("/:department_id/asset/:asset_id", utils.Handler(api.AssetApi.ModifyAssetInfo))
	group.POST("/:department_id/asset", utils.Handler(api.AssetApi.CreateAssets))
	group.PATCH("/:department_id/asset/expire", utils.Handler(api.AssetApi.ExpireAsset))
	group.POST("/:department_id/asset/transfer", utils.Handler(api.AssetApi.TransferAssets))
}
