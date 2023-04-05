package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
	"strconv"
)

type entityApi struct {
}

var EntityApi *entityApi

func newEntityApi() *entityApi {
	return &entityApi{}
}

func init() {
	EntityApi = newEntityApi()
}

func (entity *entityApi) GetEntityList(ctx *utils.Context) {
	entityList, err := service.EntityService.GetAllEntity()
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(entityList)
}

func (entity *entityApi) GetEntityByID(ctx *utils.Context) {
	param := ctx.Param("id")
	tempID, err := strconv.ParseUint(param, 10, 0)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_PARAM, myerror.INVALID_PARAM_INFO)
		return
	}
	entityID := uint(tempID)

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.USER_NOTFOUND, myerror.USER_NOTFOUND_INFO)
		return
	}
	entityInfo, err := service.EntityService.GetEntityInfoByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	entityInfoRes := define.EntityInfoResponse{
		ID:   entityInfo.ID,
		Name: entityInfo.Name,
	}

	ctx.Success(entityInfoRes)
}
