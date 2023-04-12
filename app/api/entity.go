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
检查与实体有关的查看信息权限，暂时认为只有超级管理员以及本实体的系统管理员有该权限
修改实体信息权限与此相同
*/
func (entity *entityApi) CheckViewIdentity(ctx *utils.Context) (bool, uint) {
	systemSuper := service.UserService.SystemSuper(ctx)
	entitySuper := service.UserService.EntitySuper(ctx)
	if !systemSuper && !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return false, 0
	}
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return false, 0
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false, 0
	}
	if !exists {
		ctx.NotFound(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return false, 0
	}

	if systemSuper {
		return true, entityID
	}

	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	if !isInEntity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return false, 0
	}
	return true, entityID
}

/*
Handle func for POST /entity
*/
func (entity *entityApi) CreateEntity(ctx *utils.Context) {
	// isSystemSuper := service.UserService.SystemSuper(ctx)
	// if !isSystemSuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }

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
Handle func for DELETE /entity/{entity_id}
*/
func (entity *entityApi) DeleteEntity(ctx *utils.Context) {
	// isSystemSuper := service.UserService.SystemSuper(ctx)
	// if !isSystemSuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }

	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	hasUsers, err := service.EntityService.EntityHasUser(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if hasUsers {
		ctx.BadRequest(myerror.ENTITY_HAS_USERS, myerror.ENTITY_HAS_USERS_INFO)
		return
	}

	err = service.EntityService.DeleteEntity(entityID)
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
	// isSystemSuper := service.UserService.SystemSuper(ctx)
	// if !isSystemSuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	entityList, err := service.EntityService.GetAllEntity()
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	entityListRes := []define.EntityBasicInfo{}
	err = copier.Copy(&entityListRes, &entityList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	entityListResponse := define.EntityListResponse{
		EntityList: entityListRes,
	}
	ctx.Success(entityListResponse)
}

/*
Handle func for GET /entity/:entity_id
*/
func (entity *entityApi) GetEntityByID(ctx *utils.Context) {
	// isSystemSuper := service.UserService.SystemSuper(ctx)
	// if !isSystemSuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	hasIdentity, entityID := entity.CheckViewIdentity(ctx)
	if !hasIdentity {
		return
	}

	thisEntity, err := service.EntityService.GetEntityInfoByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	// else if thisEntity == nil {
	// 	ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
	// 	return
	// }

	managerList, err := service.EntityService.GetEntityManagerList(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	managerListRes := []define.EntityManager{}
	err = copier.Copy(&managerListRes, managerList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	entityInfoRes := define.EntityInfoResponse{
		EntityID:    thisEntity.ID,
		EntityName:  thisEntity.Name,
		ManagerList: managerListRes,
		CreatedAt:   thisEntity.CreatedAt,
		Description: thisEntity.Description,
	}

	ctx.Success(entityInfoRes)
}

/*
Handle func for GET /entity/{entity_id}/user/list
*/
func (entity *entityApi) UsersInEntity(ctx *utils.Context) {
	// systemSuper := service.UserService.SystemSuper(ctx)
	// entitySuper := service.UserService.EntitySuper(ctx)
	// if !systemSuper && !entitySuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	// entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	// if err != nil {
	// 	return
	// }

	// exists, err := service.EntityService.ExistsEntityByID(entityID)
	// if err != nil {
	// 	ctx.InternalError(err.Error())
	// 	return
	// }
	// if !exists {
	// 	ctx.NotFound(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
	// 	return
	// }

	hasIdentity, entityID := entity.CheckViewIdentity(ctx)
	if !hasIdentity {
		return
	}

	userList, err := service.EntityService.GetUsersUnderEntity(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	userListRes := []define.EntityUserInfo{}
	err = copier.Copy(&userListRes, userList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	userListResponse := define.EntityUserListResponse{
		UserList: userListRes,
	}
	ctx.Success(userListResponse)
}

/*
Handler func for GET /entity/{entity_id}/department/list
*/
func (entity *entityApi) DepartmentsInEntity(ctx *utils.Context) {
	entitySuper := service.UserService.EntitySuper(ctx)
	departmentSuper := service.UserService.DepartmentSuper(ctx)
	if !entitySuper && !departmentSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	exists, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if !exists {
		ctx.NotFound(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	if !isInEntity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	departmentList, err := service.EntityService.GetAllDepartmentsUnderEntity(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	departmentListRes := []define.DepartmentBasicInfo{}
	err = copier.Copy(&departmentListRes, departmentList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	departmentListResponse := define.DepartmentListResponse{
		DepartmentList: departmentListRes,
	}
	ctx.Success(departmentListResponse)
}

/*
Handle func for POST /entity/{entity_id}/manager
*/
func (entity *entityApi) SetManager(ctx *utils.Context) {
	// isSystemSuper := service.UserService.SystemSuper(ctx)
	// if !isSystemSuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	var setManagerReq define.ManagerReq
	err = ctx.MustBindWith(&setManagerReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}
	thisUser, err := service.UserService.GetUserByName(setManagerReq.Username)
	if err != nil {
		ctx.InternalError(err.Error())
	}
	// Update user to manager
	if setManagerReq.Password == nil {
		if thisUser == nil {
			ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
			return
		}
		if thisUser.EntityID != entityID {
			ctx.BadRequest(myerror.USER_NOT_IN_ENTITY, myerror.USER_NOT_IN_ENTITY_INFO)
			return
		}
		err = service.EntityService.SetManager(setManagerReq.Username)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	} else { // Create a manager
		if thisUser != nil {
			ctx.BadRequest(myerror.USER_HAS_EXISTED, myerror.USER_HAS_EXISTED_INFO)
			return
		}
		err = service.EntityService.CreateManager(setManagerReq.Username, *setManagerReq.Password, entityID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	}
	ctx.Success(nil)
}

/*
Handle func for DELETE /entity/{entity_id}/manager/{user_id}
*/
func (entity *entityApi) DeleteManager(ctx *utils.Context) {
	// isSystemSuper := service.UserService.SystemSuper(ctx)
	// if !isSystemSuper {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}
	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	} else if thisUser.EntityID != entityID {
		ctx.BadRequest(myerror.USER_NOT_IN_ENTITY, myerror.USER_NOT_IN_ENTITY_INFO)
		return
	}

	err = service.EntityService.DeleteManager(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for PATCH /entity/{entity_id}
*/
func (entity *entityApi) ModifyEntityInfo(ctx *utils.Context) {
	hasIdentity, entityID := entity.CheckViewIdentity(ctx)
	if !hasIdentity {
		return
	}

	var modifyEntityInfoReq define.ModifyEntityInfoReq
	err := ctx.MustBindWith(&modifyEntityInfoReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}
	if modifyEntityInfoReq.EntityName != nil && *modifyEntityInfoReq.EntityName == "" {
		ctx.BadRequest(myerror.NAME_CANNOT_EMPTY, myerror.NAME_CANNOT_EMPTY_INFO)
		return
	}

	err = service.EntityService.ModifyEntity(entityID, modifyEntityInfoReq)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}
