package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type taskRouter struct{}

var TaskRouter *taskRouter

func newTaskRouter() *taskRouter {
	return &taskRouter{}
}

func init() {
	TaskRouter = newTaskRouter()
}

func (task *taskRouter) Init(group *gin.RouterGroup) {
	task.routerUserTask(group.Group("/users"))
	task.routerDepartmentTask(group.Group("/department"))
}

func (task *taskRouter) routerUserTask(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	group.POST("/:user_id/assets/task", utils.Handler(api.TaskApi.CreateNewTask))
	group.GET("/:user_id/assets/tasks", utils.Handler(api.TaskApi.GetUserTaskList))
	group.GET("/:user_id/assets/tasks/:task_id", utils.Handler(api.TaskApi.GetUserTaskInfo))
	group.DELETE("/:user_id/assets/tasks/:task_id", utils.Handler(api.TaskApi.CancelTasks))
	group.GET("/:user_id/assets/maintain", utils.Handler(api.AssetApi.GetUserMaintainAssets))
	group.POST("/:user_id/assets/:asset_id/maintain", utils.Handler(api.AssetApi.FinishMaintenance))
}

func (task *taskRouter) routerDepartmentTask(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()), utils.Handler(middleware.LogMiddleware()))
	group.GET("/:department_id/assets/tasks", utils.Handler(api.TaskApi.GetDepartmentTaskList))
	group.GET("/:department_id/assets/tasks/:task_id", utils.Handler(api.TaskApi.GetDepartmentTaskInfo))
	group.POST("/:department_id/assets/tasks/:task_id", utils.Handler(api.TaskApi.ApproveTask))
	group.DELETE("/:department_id/assets/tasks/:task_id", utils.Handler(api.TaskApi.RejectTask))
}
