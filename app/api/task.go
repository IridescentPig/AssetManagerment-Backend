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

func getTaskInfoRes(taskList []*model.Task) define.TaskListResponse {
	taskInfoList := funk.Map(taskList, func(thisTask *model.Task) define.TaskBasicInfo {
		taskInfo := define.TaskBasicInfo{
			ID:              thisTask.ID,
			TaskType:        thisTask.TaskType,
			TaskDescription: thisTask.TaskDescription,
			UserID:          thisTask.UserID,
			UserName:        thisTask.User.UserName,
			State:           thisTask.State,
		}
		return taskInfo
	}).([]define.TaskBasicInfo)

	taskListRes := define.TaskListResponse{
		TaskList: taskInfoList,
	}

	return taskListRes
}

func userTaskPrevilige(ctx *utils.Context) (*model.Task, bool) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return nil, false
	}
	taskID, err := service.EntityService.GetParamID(ctx, "task_id")
	if err != nil {
		return nil, false
	}
	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return nil, false
	}

	taskInfo, err := service.TaskService.GetTaskInfoByID(taskID)
	if err != nil {
		ctx.InternalError(err.Error())
		return nil, false
	} else if taskInfo == nil {
		ctx.NotFound(myerror.TASK_NOT_FOUND, myerror.TASK_NOT_FOUND_INFO)
		return nil, false
	}

	if taskInfo.UserID != userID {
		ctx.BadRequest(myerror.TASK_NOT_BELONG_TO_USER, myerror.TASK_NOT_BELONG_TO_USER_INFO)
		return nil, false
	}

	return taskInfo, true
}

func departmentTaskPrevillige(ctx *utils.Context) (*model.Task, *define.UserBasicInfo, bool) {
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return nil, nil, false
	}
	taskID, err := service.EntityService.GetParamID(ctx, "task_id")
	if err != nil {
		return nil, nil, false
	}
	thisUser := UserApi.GetOperatorInfo(ctx)
	if !thisUser.DepartmentSuper || thisUser.DepartmentID != departmentID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return nil, nil, false
	}

	taskInfo, err := service.TaskService.GetTaskInfoByID(taskID)
	if err != nil {
		ctx.InternalError(err.Error())
		return nil, nil, false
	} else if taskInfo == nil {
		ctx.NotFound(myerror.TASK_NOT_FOUND, myerror.TASK_NOT_FOUND_INFO)
		return nil, nil, false
	}

	if taskInfo.DepartmentID != departmentID {
		ctx.BadRequest(myerror.TASK_NOT_IN_DEPARTMENT, myerror.TASK_NOT_IN_DEPARTMENT_INFO)
		return nil, nil, false
	}

	return taskInfo, thisUser, true
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
		assetList, err = service.AssetService.GetUserAssetsByIDs(assetIdList, thisUser.UserID)
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
		} else if targetUser.DepartmentID == 0 {
			ctx.BadRequest(myerror.TARGET_NOT_IN_DEPARTMENT, myerror.TARGET_NOT_IN_DEPARTMENT_INFO)
			return
		}

		assetList, err = service.AssetService.GetUserAssetsByIDs(assetIdList, thisUser.UserID)
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
func (task *taskApi) GetUserTaskList(ctx *utils.Context) {
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

	// taskInfoList := funk.Map(taskList, func(thisTask *model.Task) define.TaskBasicInfo {
	// 	taskInfo := define.TaskBasicInfo{
	// 		ID:              thisTask.ID,
	// 		TaskType:        thisTask.TaskType,
	// 		TaskDescription: thisTask.TaskDescription,
	// 		UserID:          thisTask.UserID,
	// 		UserName:        thisTask.User.UserName,
	// 		State:           thisTask.State,
	// 	}
	// 	return taskInfo
	// }).([]define.TaskBasicInfo)

	// taskListRes := define.TaskListResponse{
	// 	TaskList: taskInfoList,
	// }

	taskListRes := getTaskInfoRes(taskList)

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

	taskListRes := getTaskInfoRes(taskList)

	ctx.Success(taskListRes)
}

/*
Handle func for GET /departments/:department_id/assets/tasks/:task_id
*/
func (task *taskApi) GetDepartmentTaskInfo(ctx *utils.Context) {
	// departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	// if err != nil {
	// 	return
	// }
	// task_id, err := service.EntityService.GetParamID(ctx, "task_id")
	// if err != nil {
	// 	return
	// }
	// thisUser := UserApi.GetOperatorInfo(ctx)
	// if !thisUser.DepartmentSuper || thisUser.DepartmentID != departmentID {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }

	// taskInfo, err := service.TaskService.GetTaskInfoByID(task_id)
	// if err != nil {
	// 	ctx.InternalError(err.Error())
	// 	return
	// } else if taskInfo == nil {
	// 	ctx.NotFound(myerror.TASK_NOT_FOUND, myerror.TASK_NOT_FOUND_INFO)
	// 	return
	// }

	// if taskInfo.DepartmentID != departmentID {
	// 	ctx.BadRequest(myerror.TASK_NOT_IN_DEPARTMENT, myerror.TASK_NOT_IN_DEPARTMENT_INFO)
	// 	return
	// }

	taskInfo, _, isOK := departmentTaskPrevillige(ctx)
	if !isOK {
		return
	}

	taskInfoRes := define.TaskInfo{
		ID:              taskInfo.ID,
		TaskType:        taskInfo.TaskType,
		TaskDescription: taskInfo.TaskDescription,
		UserID:          taskInfo.UserID,
		UserName:        taskInfo.User.UserName,
		TargetID:        taskInfo.TargetID,
		TargetName:      taskInfo.Target.UserName,
		DepartmentID:    taskInfo.DepartmentID,
		DepartmentName:  taskInfo.Department.Name,
		AssetList:       taskInfo.AssetList,
		State:           taskInfo.State,
	}

	ctx.Success(taskInfoRes)
}

/*
Handle func for GET /users/:user_id/assets/tasks/:task_id
*/
func (task *taskApi) GetUserTaskInfo(ctx *utils.Context) {
	// userID, err := service.EntityService.GetParamID(ctx, "user_id")
	// if err != nil {
	// 	return
	// }
	// task_id, err := service.EntityService.GetParamID(ctx, "task_id")
	// if err != nil {
	// 	return
	// }
	// thisUser := UserApi.GetOperatorInfo(ctx)
	// if thisUser.UserID != userID {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }

	// taskInfo, err := service.TaskService.GetTaskInfoByID(task_id)
	// if err != nil {
	// 	ctx.InternalError(err.Error())
	// 	return
	// } else if taskInfo == nil {
	// 	ctx.NotFound(myerror.TASK_NOT_FOUND, myerror.TASK_NOT_FOUND_INFO)
	// 	return
	// }

	// if taskInfo.UserID != userID {
	// 	ctx.BadRequest(myerror.TASK_NOT_BELONG_TO_USER, myerror.TASK_NOT_BELONG_TO_USER_INFO)
	// 	return
	// }

	taskInfo, isOK := userTaskPrevilige(ctx)
	if !isOK {
		return
	}

	taskInfoRes := define.TaskInfo{
		ID:              taskInfo.ID,
		TaskType:        taskInfo.TaskType,
		TaskDescription: taskInfo.TaskDescription,
		UserID:          taskInfo.UserID,
		UserName:        taskInfo.User.UserName,
		TargetID:        taskInfo.TargetID,
		TargetName:      taskInfo.Target.UserName,
		DepartmentID:    taskInfo.DepartmentID,
		DepartmentName:  taskInfo.Department.Name,
		AssetList:       taskInfo.AssetList,
		State:           taskInfo.State,
	}

	ctx.Success(taskInfoRes)
}

/*
Handle func for POST /departments/:department_id/assets/tasks/:task_id
*/
func (task *taskApi) ApproveTask(ctx *utils.Context) {
	taskInfo, thisUser, isOK := departmentTaskPrevillige(ctx)
	if !isOK {
		return
	}

	if taskInfo.User.DepartmentID != taskInfo.DepartmentID {
		ctx.BadRequest(myerror.USER_NOT_IN_DEPARTMENT, myerror.USER_NOT_IN_DEPARTMENT_INFO)
		return
	}

	if taskInfo.State != 0 {
		ctx.BadRequest(myerror.TASK_NOT_PENDING, myerror.TASK_NOT_PENDING_INFO)
		return
	}

	assetIDs := funk.Map(taskInfo.AssetList, func(thisAsset *model.Asset) uint {
		return thisAsset.ID
	}).([]uint)

	if taskInfo.TaskType == 0 {
		assetList, err := service.AssetService.GetDepartmentIdleAssets(assetIDs, taskInfo.DepartmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if len(assetList) != len(assetIDs) {
			ctx.BadRequest(myerror.ASSET_LIST_INVALID, myerror.ASSET_LIST_INVALID_INFO)
			return
		}

		err = service.AssetService.AcquireAssets(assetIDs, taskInfo.UserID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	} else if taskInfo.TaskType == 1 {
		assetList, err := service.AssetService.GetUserAssetsByIDs(assetIDs, taskInfo.UserID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if len(assetList) != len(assetIDs) {
			ctx.BadRequest(myerror.ASSET_LIST_INVALID, myerror.ASSET_LIST_INVALID_INFO)
			return
		}

		err = service.AssetService.CancelAssets(assetIDs, thisUser.UserID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	} else {
		assetList, err := service.AssetService.GetUserAssetsByIDs(assetIDs, taskInfo.UserID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if len(assetList) != len(assetIDs) {
			ctx.BadRequest(myerror.ASSET_LIST_INVALID, myerror.ASSET_LIST_INVALID_INFO)
			return
		}

		targetUser, err := service.UserService.GetUserByID(taskInfo.TargetID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		} else if targetUser == nil {
			ctx.BadRequest(myerror.TARGET_USER_NOT_FOUND, myerror.TARGET_USER_NOT_FOUND_INFO)
			return
		}

		if targetUser.EntityID != thisUser.EntityID {
			ctx.BadRequest(myerror.NOT_IN_SAME_ENTITY, myerror.NOT_IN_SAME_ENTITY_INFO)
			return
		} else if targetUser.DepartmentID == 0 {
			ctx.BadRequest(myerror.TARGET_NOT_IN_DEPARTMENT, myerror.TARGET_NOT_IN_DEPARTMENT_INFO)
			return
		}
		if taskInfo.TaskType == 2 {
			err = service.AssetService.ModifyAssetMaintainerAndState(assetIDs, taskInfo.TargetID)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			}
		} else {
			err = service.AssetService.TransferAssets(assetIDs, taskInfo.TargetID, taskInfo.Target.DepartmentID, taskInfo.DepartmentID)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			}
		}
	}

	err := service.TaskService.ModifyTaskState(taskInfo.ID, 1)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for DELETE /departments/:department_id/assets/tasks/:task_id
*/
func (task *taskApi) RejectTask(ctx *utils.Context) {
	taskInfo, _, isOK := departmentTaskPrevillige(ctx)
	if !isOK {
		return
	}

	if taskInfo.State != 0 {
		ctx.BadRequest(myerror.TASK_NOT_PENDING, myerror.TASK_NOT_PENDING_INFO)
		return
	}

	err := service.TaskService.ModifyTaskState(taskInfo.ID, 2)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for DELETE /users/:user_id/assets/tasks/:task_id
*/
func (task *taskApi) CancelTasks(ctx *utils.Context) {
	taskInfo, isOK := userTaskPrevilige(ctx)
	if !isOK {
		return
	}

	if taskInfo.State != 0 {
		ctx.BadRequest(myerror.TASK_NOT_PENDING, myerror.TASK_NOT_PENDING_INFO)
		return
	}

	err := service.TaskService.ModifyTaskState(taskInfo.ID, 3)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}
