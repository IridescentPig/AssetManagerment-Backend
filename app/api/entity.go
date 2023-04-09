package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
	"strconv"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
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

/*
Handle func for POST /entity
*/
func (entity *entityApi) CreateEntity(ctx *utils.Context) {
	isSystemSuper := service.UserService.SystemSuper(ctx)
	if !isSystemSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var createReq define.CreateEntityReq
	err := ctx.MustBindWith(&createReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}
	isExist, err := service.EntityService.ExistsEntityByName(createReq.EntityName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if isExist {
		ctx.BadRequest(myerror.DUPLICATED_NAME, myerror.DUPLICATED_NAME_INFO)
		return
	}

	err = service.EntityService.CreateEntity(createReq.EntityName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

/*
Handle func for GET /entity/list
*/
func (entity *entityApi) GetEntityList(ctx *utils.Context) {
	entityList, err := service.EntityService.GetAllEntity()
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	entityListRes := []define.EntityBasicInfo{}
	copier.Copy(&entityListRes, entityList)
	entityListResponse := define.EntityListResponse{
		EntityList: entityListRes,
	}
	ctx.Success(entityListResponse)
	// ctx.Success(entityList)
}

/*
Handle func for GET /entity/:entity_id
*/
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
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}
	entityInfo, err := service.EntityService.GetEntityInfoByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	entityInfoRes := define.EntityInfoResponse{
		EntityID:   entityInfo.ID,
		EntityName: entityInfo.Name,
	}

	ctx.Success(entityInfoRes)
}

/*
Handle func for GET /entity/{entity_id}/user/list
*/
func (entity *entityApi) UsersInEntity(ctx *utils.Context) {
	systemSuper := service.UserService.SystemSuper(ctx)
	entitySuper := service.UserService.EntitySuper(ctx)
	if !systemSuper && !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	ctx.Success(nil)
}

/*
只获取最高一级的部门
是否可以与 GetEntity 的响应合并？
*/
func (entity *entityApi) DepartmentsInEntity(ctx *utils.Context) {

}
