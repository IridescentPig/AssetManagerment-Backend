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
	group.POST("/:department_id/asset/:asset_id/property", utils.Handler(api.AssetApi.CreateAssetProperty))
	group.PATCH("/:department_id/asset/:asset_id/property", utils.Handler(api.AssetApi.ModifyAssetProperty))
	group.DELETE("/:department_id/asset/:asset_id/property", utils.Handler(api.AssetApi.DeleteAssetProperty))
	group.GET("/:department_id/asset/:asset_id/history", utils.Handler(api.AssetApi.GetAssetHistory))
	group.POST("/:department_id/asset/search", utils.Handler(api.AssetApi.SearchAssets))
	group.POST("/:department_id/asset/search/spare", utils.Handler(api.AssetApi.SearchSpareAssets))
	group.GET("/:department_id/asset/stat/total", utils.Handler(api.StatApi.GetDepartmentStatTotal))
	group.GET("/:department_id/asset/stat/distribution", utils.Handler(api.StatApi.GetDepartmentStatDistribution))
	group.GET("/:department_id/asset/stat/sub", utils.Handler(api.StatApi.GetSubDepartmentsAssetDistribution))
	group.GET("/:department_id/asset/:asset_id", utils.Handler(api.AssetApi.GetAssetInfo))
	group.POST("/:department_id/template", utils.Handler(api.DepartmentApi.DefineDepartmentAssetTemplate))
	group.GET("/:department_id/template", utils.Handler(api.DepartmentApi.GetDepartmentTemplate))
	group.POST("/:department_id/warn", utils.Handler(api.DepartmentApi.SetDepartmentWarnStrategy))
	group.GET("/:department_id/warn", utils.Handler(api.DepartmentApi.GetDepartmentAssetWarnInfo))
}
