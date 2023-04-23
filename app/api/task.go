package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/thoas/go-funk"
)

type taskApi struct {
}

var TaskApi *taskApi

func newTaskApi() *taskApi {
	return &taskApi{}
}

func init() {
	TaskApi = newTaskApi()
}

/*
Handle func for /users/:user_id/assets/task
*/
func (task *taskApi) CreateNewTask(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	existUser, err := service.UserService.ExistsUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !existUser {
		ctx.NotFound(myerror.USER_NOT_FOUND, myerror.USER_HAS_EXISTED_INFO)
		return
	}

	if thisUser.DepartmentID == 0 {
		ctx.Forbidden(myerror.USER_NOT_IN_DEPARTMENT, myerror.USER_NOT_IN_DEPARTMENT_INFO)
		return
	}

	var req define.CreateTaskReq
	err = ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	assetIdList := []uint{}
	var assetList []*model.Asset

	for _, asset := range req.AssetList {
		assetIdList = append(assetIdList, asset.AssetID)
	}

	assetIdList = funk.UniqUInt(assetIdList)

	if req.TaskType == 0 {
		assetList, err = service.AssetService.GetDepartmentAssetsByIDs(assetIdList, thisUser.DepartmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if len(assetList) != len(assetIdList) {
			ctx.BadRequest(myerror.ASSET_LIST_INVALID, myerror.ASSET_LIST_INVALID_INFO)
			return
		}
		req.TargetID = 0
	} else if req.TaskType == 1 {
		assetList, err = service.AssetService.GetDepartmentAssetsByIDs(assetIdList, thisUser.DepartmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if len(assetList) != len(assetIdList) {
			ctx.BadRequest(myerror.ASSET_LIST_INVALID, myerror.ASSET_LIST_INVALID_INFO)
			return
		}
		req.TargetID = 0
	} else {
		if req.TargetID == 0 {
			ctx.BadRequest(myerror.TARGET_EMPTY, myerror.TARGET_EMPTY_INFO)
			return
		}
		targetUser, err := service.UserService.GetUserByID(userID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if targetUser == nil {
			ctx.BadRequest(myerror.TARGET_USER_NOT_FOUND, myerror.TARGET_USER_NOT_FOUND_INFO)
			return
		}

		if thisUser.EntityID != targetUser.EntityID {
			ctx.BadRequest(myerror.NOT_IN_SAME_ENTITY, myerror.NOT_IN_SAME_ENTITY_INFO)
			return
		}

		assetList, err = service.AssetService.GetDepartmentAssetsByIDs(assetIdList, thisUser.DepartmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if len(assetList) != len(assetIdList) {
			ctx.BadRequest(myerror.ASSET_LIST_INVALID, myerror.ASSET_LIST_INVALID_INFO)
			return
		}
	}

	err = service.TaskService.CreateTask(req, thisUser.UserID, thisUser.DepartmentID, assetList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

/*
Handle func for GET /user/:user_id/assets/tasks
*/
func (task *taskApi) GetUserTask(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	existUser, err := service.UserService.ExistsUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !existUser {
		ctx.NotFound(myerror.USER_NOT_FOUND, myerror.USER_HAS_EXISTED_INFO)
		return
	}

	taskList, err := service.TaskService.GetTasksByUserID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	taskInfoList := funk.Map(taskList, func(thisTask *model.Task) define.TaskInfo {
		taskInfo := define.TaskInfo{
			ID:              thisTask.ID,
			TaskType:        thisTask.TaskType,
			TaskDescription: thisTask.TaskDescription,
			UserID:          thisTask.UserID,
			UserName:        thisTask.User.UserName,
			DepartmentID:    thisTask.DepartmentID,
			DepartmentName:  thisTask.Department.Name,
			AssetList:       thisTask.AssetList,
			State:           thisTask.State,
		}
		if thisTask.TargetID != 0 {
			taskInfo.TargetID = thisTask.TargetID
			taskInfo.TargetName = thisTask.Target.UserName
		}
		return taskInfo
	}).([]define.TaskInfo)

	taskListRes := define.TaskListResponse{
		TaskList: taskInfoList,
	}

	ctx.Success(taskListRes)
}

/*
Handle func for GET /departments/:department_id/assets/tasks
*/
func (task *taskApi) GetDepartmentTaskList(ctx *utils.Context) {
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return
	}

	departmentSuper := service.UserService.DepartmentSuper(ctx)
	thisUser := UserApi.GetOperatorInfo(ctx)
	if !departmentSuper || thisUser.DepartmentID != departmentID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	existDepartment, err := service.DepartmentService.ExistsDepartmentByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !existDepartment {
		ctx.NotFound(myerror.DEPARTMENT_NOT_FOUND, myerror.DEPARTMENT_NOT_FOUND_INFO)
		return
	}

	taskList, err := service.TaskService.GetTasksByDepartmentID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	taskInfoList := funk.Map(taskList, func(thisTask *model.Task) define.TaskInfo {
		taskInfo := define.TaskInfo{
			ID:              thisTask.ID,
			TaskType:        thisTask.TaskType,
			TaskDescription: thisTask.TaskDescription,
			UserID:          thisTask.UserID,
			UserName:        thisTask.User.UserName,
			DepartmentID:    thisTask.DepartmentID,
			DepartmentName:  thisTask.Department.Name,
			AssetList:       thisTask.AssetList,
			State:           thisTask.State,
		}
		if thisTask.TargetID != 0 {
			taskInfo.TargetID = thisTask.TargetID
			taskInfo.TargetName = thisTask.Target.UserName
		}
		return taskInfo
	}).([]define.TaskInfo)

	taskListRes := define.TaskListResponse{
		TaskList: taskInfoList,
	}

	ctx.Success(taskListRes)
}
