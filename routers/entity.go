package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type entityRouter struct{}

var EntityRouter *entityRouter

func newEntityRouter() *entityRouter {
	return &entityRouter{}
}

func init() {
	EntityRouter = newEntityRouter()
}

func (entity *entityRouter) Init(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	{
		entity.UrlrouterCheckAtHandler(group)
	}
	entity.routerCheckAtHandler(group)
	entity.routerNeedSystemSuper(group)
}

func (entity *entityRouter) routerNeedSystemSuper(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.CheckSystemSuper()))
	{
		group.POST("/", utils.Handler(api.EntityApi.CreateEntity))
		group.DELETE("/:entity_id", utils.Handler(api.EntityApi.DeleteEntity))
		group.GET("/list", utils.Handler(api.EntityApi.GetEntityList))
		group.GET("/:entity_id", utils.Handler(api.EntityApi.GetEntityByID))
		group.POST("/:entity_id/manager", utils.Handler(api.EntityApi.SetManager))
		group.DELETE("/:entity_id/manager/:user_id", utils.Handler(api.EntityApi.DeleteManager))
	}
}

func (entity *entityRouter) routerCheckAtHandler(group *gin.RouterGroup) {
	group.GET("/:entity_id/user/list", utils.Handler(api.EntityApi.UsersInEntity))
	group.GET("/:entity_id/department/list", utils.Handler(api.EntityApi.DepartmentsInEntity)) // change later
	group.PATCH("/:entity_id", utils.Handler(api.EntityApi.ModifyEntityInfo))
	group.GET("/:entity_id/department/sub", utils.Handler(api.EntityApi.GetEntitySubDepartments))

	group.POST("/:entity_id/department", utils.Handler(api.DepartmentApi.CreateDepartment))
	group.POST("/:entity_id/department/:department_id/department", utils.Handler(api.DepartmentApi.CreateDepartment))
	group.DELETE("/:entity_id/department/:department_id", utils.Handler(api.DepartmentApi.DeleteDepartment))
	group.GET("/:entity_id/department/:department_id", utils.Handler(api.DepartmentApi.GetDepartmentByID))
	group.GET("/:entity_id/department/:department_id/department/list", utils.Handler(api.DepartmentApi.GetSubDepartments))
	group.GET("/:entity_id/department/:department_id/user/list", utils.Handler(api.DepartmentApi.GetAllUsersUnderDepartment))
	group.POST("/:entity_id/department/:department_id/user", utils.Handler(api.DepartmentApi.CreateUserInDepartment))
	group.POST("/:entity_id/department/:department_id/manager", utils.Handler(api.DepartmentApi.SetManager))
	group.DELETE("/:entity_id/department/:department_id/manager/:user_id", utils.Handler(api.DepartmentApi.DeleteDepartmentManager))
	group.GET("/:entity_id/department/:department_id/manager", utils.Handler(api.DepartmentApi.GetDepartmentManager))
	group.GET("/:entity_id/department/tree", utils.Handler(api.DepartmentApi.GetDepartmentTree))
	group.GET("/:entity_id/department/:department_id/user/sub", utils.Handler(api.DepartmentApi.GetDepartmentSubUsers))
}

func (entity *entityRouter) UrlrouterCheckAtHandler(group *gin.RouterGroup) {
	group.GET("/:entity_id/url", utils.Handler(api.UrlApi.GetUrlsByEntity))
	group.POST("/:entity_id/url", utils.Handler(api.UrlApi.CreateUrl))
	group.PATCH("/:entity_id/url", utils.Handler(api.UrlApi.ModifyUrl))
	group.DELETE("/:entity_id/url", utils.Handler(api.UrlApi.DeleteUrl))
	group.GET("/:entity_id/url/list", utils.Handler(api.UrlApi.GetUrlList))
}
