package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

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
	copier.Copy(&entityListRes, &entityList)
	entityListResponse := define.EntityListResponse{
		EntityList: entityListRes,
	}
	ctx.Success(entityListResponse)
}

/*
Handle func for GET /entity/:entity_id
*/
func (entity *entityApi) GetEntityByID(ctx *utils.Context) {
	entityID, err := service.EntityService.GetParamID(ctx)
	if err != nil {
		return
	}

	thisEntity, err := service.EntityService.GetEntityInfoByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisEntity == nil {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	entityInfoRes := define.EntityInfoResponse{
		EntityID:   thisEntity.ID,
		EntityName: thisEntity.Name,
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
	entityID, err := service.EntityService.GetParamID(ctx)
	if err != nil {
		return
	}
	userList, err := service.EntityService.GetUsersUnderEntity(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	userListRes := []define.EntityUserInfo{}
	copier.Copy(&userListRes, userList)
	userListResponse := define.EntityUserListResponse{
		UserList: userListRes,
	}
	ctx.Success(userListResponse)
}

/*
/entity/{entity_id}/department/list
*/
func (entity *entityApi) DepartmentsInEntity(ctx *utils.Context) {
	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	entityID, err := service.EntityService.GetParamID(ctx)
	if err != nil {
		return
	}
	departmentList, err := service.EntityService.GetAllDepartmentsUnderEntity(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	departmentListRes := []define.DepartmentBasicInfo{}
	copier.Copy(&departmentListRes, departmentList)
	departmentListResponse := define.DepartmentListResponse{
		DepartmentList: departmentListRes,
	}
	ctx.Success(departmentListResponse)
}
