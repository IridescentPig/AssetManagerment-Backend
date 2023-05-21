package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/thoas/go-funk"
)

type asyncApi struct {
}

var AsyncApi *asyncApi

func newAsyncApi() *asyncApi {
	return &asyncApi{}
}

func init() {
	AsyncApi = newAsyncApi()
}

/*
Handle func for GET /user/:user_id/async/list
*/
func (asy *asyncApi) GetUserAsyncTasks(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	userInfo, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if userInfo == nil {
		ctx.NotFound(myerror.USER_NOT_FOUND, myerror.USER_HAS_EXISTED_INFO)
		return
	}

	taskList, err := service.AsyncService.GetUserAsyncTasks(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	taskListRes := funk.Map(taskList, func(task *model.AsyncTask) define.AsyncTaskInfo {
		return define.AsyncTaskInfo{
			ID:           task.ID,
			Type:         task.Type,
			UserID:       task.UserID,
			Username:     task.User.UserName,
			State:        task.State,
			DownloadLink: task.DownloadLink,
			Message:      task.Message,
			LogType:      task.LogType,
		}
	}).([]define.AsyncTaskInfo)

	taskListResponse := define.AsyncTaskListResponse{
		AsyncList: taskListRes,
	}

	ctx.Success(taskListResponse)
}

/*
Handle func for POST /user/:user_id/async
*/
func (asy *asyncApi) CreateAsyncTask(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	userInfo, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	var req define.CreateAsyncTaskReq
	err = ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		log.Println(err.Error())
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	if req.Type == 0 {
		if !userInfo.DepartmentSuper {
			ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
			return
		} else if userInfo.DepartmentID != req.DepartmentID {
			ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
			return
		}

		req.EntityID = userInfo.EntityID

		err = service.AsyncService.CreateAsyncTask(userID, &req)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	} else {
		if !userInfo.EntitySuper {
			ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
			return
		} else if userInfo.EntityID != req.EntityID {
			ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
			return
		}

		req.DepartmentID = userInfo.DepartmentID

		err = service.AsyncService.CreateAsyncTask(userID, &req)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	}

	ctx.Success(nil)
}

/*
Handle func for PATCH /user/:user_id/async/:task_id
*/
func (asy *asyncApi) ModifyAsyncState(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}
	taskID, err := service.EntityService.GetParamID(ctx, "task_id")
	if err != nil {
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	taskInfo, err := service.AsyncService.GetAsyncTaskByID(taskID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if taskInfo.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var req define.ModifyAsyncTaskStateReq
	err = ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	if taskInfo.Type == 0 {
		if !thisUser.DepartmentSuper || thisUser.DepartmentID != taskInfo.DepartmentID {
			ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
			return
		}
	} else {
		if !thisUser.EntitySuper || thisUser.EntityID != taskInfo.EntityID {
			ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
			return
		}
	}

	if req.State == 4 && taskInfo.State != 0 {
		ctx.BadRequest(myerror.TASK_NOT_PENDING, myerror.TASK_NOT_PENDING_INFO)
		return
	} else if req.State == 0 && taskInfo.State == 1 {
		ctx.BadRequest(myerror.RUNNING_TASK_CANNOT_BE_RESTART, myerror.RUNNING_TASK_CANNOT_BE_RESTART_INFO)
		return
	}

	err = service.AsyncService.ModifyAsyncTaskState(taskID, req.State)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}
