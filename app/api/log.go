package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/jinzhu/copier"
)

type logApi struct {
}

var LogApi *logApi

func newLogApi() *logApi {
	return &logApi{}
}

func init() {
	LogApi = newLogApi()
}

/*
Handle func for GET /entity/:entity_id/login-logs
*/
func (mylog *logApi) GetLoginLog(ctx *utils.Context) {
	isSysyemSuper := service.UserService.SystemSuper(ctx)
	isEntitySuper := service.UserService.EntitySuper(ctx)
	if !isEntitySuper && !isSysyemSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	userInfo := UserApi.GetOperatorInfo(ctx)
	if !isSysyemSuper && userInfo.EntityID != entityID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	logList, err := service.LogService.GetLoginLog(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	var logListRes []*define.LogInfo

	err = copier.Copy(&logListRes, logList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	logListResponse := define.LogListResponse{
		LogList: logListRes,
	}

	ctx.Success(logListResponse)
}

/*
Handle func for GET /entity/:entity_id/data-logs
*/
func (mylog *logApi) GetDataLog(ctx *utils.Context) {
	isSysyemSuper := service.UserService.SystemSuper(ctx)
	isEntitySuper := service.UserService.EntitySuper(ctx)
	if !isEntitySuper && !isSysyemSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	userInfo := UserApi.GetOperatorInfo(ctx)
	if !isSysyemSuper && userInfo.EntityID != entityID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	logList, err := service.LogService.GetDataLog(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	var logListRes []*define.LogInfo

	err = copier.Copy(&logListRes, logList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	logListResponse := define.LogListResponse{
		LogList: logListRes,
	}

	ctx.Success(logListResponse)
}
