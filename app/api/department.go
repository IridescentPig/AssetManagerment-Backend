package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
)

type departmentApi struct {
}

var DepartmentApi *departmentApi

func newdepartmentApi() *departmentApi {
	return &departmentApi{}
}

func init() {
	DepartmentApi = newdepartmentApi()
}

func (department *departmentApi) CheckEntityDepartmentValid(ctx *utils.Context, entityID uint, departmentID uint) bool {
	existsEntity, err := service.EntityService.ExistsEntityByID(entityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false
	}
	if !existsEntity {
		ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return false
	}

	existsDepartment, err := service.DepartmentService.ExistsDepartmentByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false
	}
	if !existsDepartment {
		ctx.BadRequest(myerror.DEPARTMENT_NOT_FOUND, myerror.DEPARTMENT_NOT_FOUND_INFO)
		return false
	}

	flag, err := service.DepartmentService.CheckDepartmentInEntity(entityID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false
	}
	if !flag {
		ctx.BadRequest(myerror.DEPARTMENT_NOT_IN_ENTITY, myerror.DEPARTMENT_NOT_IN_ENTITY_INFO)
		return false
	}
	return true
}

func (department *departmentApi) GetTwoIDs(ctx *utils.Context) (uint, uint, error) {
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return 0, 0, err
	}
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return 0, 0, err
	}
	return entityID, departmentID, nil
}

/*
Handle func for POST /entity/{entity_id}/department and /entity/{entity_id}/department/{department_id}/department
*/
func (department *departmentApi) CreateDepartment(ctx *utils.Context) {
	isEntitySuper := service.UserService.EntitySuper(ctx)
	if !isEntitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}

	var createDepartmentReq define.CreateDepartmentReq
	param := ctx.Param("department_id")
	if param == "" {
		existsEntity, err := service.EntityService.ExistsEntityByID(entityID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if !existsEntity {
			ctx.BadRequest(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
			return
		}
		createDepartmentReq = define.CreateDepartmentReq{
			EntityID: entityID,
		}
	} else {
		departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
		if err != nil {
			return
		}
		isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
		if !isValid {
			return
		}
		createDepartmentReq = define.CreateDepartmentReq{
			EntityID:     entityID,
			DepartmentID: departmentID,
		}
	}

	err = service.DepartmentService.CreateDepartment(createDepartmentReq.EntityID, createDepartmentReq.DepartmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

/*
Handle func for GET /entity/{entity_id}/department/{department_id}
*/
func (department *departmentApi) GetDepartmentByID(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}
	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}
	identity, err := service.DepartmentService.CheckDepartmentIdentity(ctx, entityID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if !identity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	thisDepartment, err := service.DepartmentService.GetDepartmentInfoByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	departmentBasicInfo := define.DepartmentBasicInfo{
		ID:       thisDepartment.ID,
		Name:     thisDepartment.Name,
		ParentID: thisDepartment.ParentID,
	}
	entityBasicInfo := define.EntityBasicInfo{
		ID:   thisDepartment.EntityID,
		Name: thisDepartment.Entity.Name,
	}
	departmentInfo := define.DepartmentInfo{
		DepartmentBasicInfo: departmentBasicInfo,
		Entity:              entityBasicInfo,
	}
	departmentInfoRes := define.DepartmentInfoResponse{
		Department: departmentInfo,
	}
	ctx.Success(departmentInfoRes)
}

/*
Handle func for GET /entity/{entity_id}/department/{department_id}/department/list
*/
func (department *departmentApi) GetSubDepartments(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}

	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	departmentSuper := service.UserService.DepartmentSuper(ctx)
	if !entitySuper && !departmentSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	identity, err := service.DepartmentService.CheckDepartmentIdentity(ctx, entityID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if !identity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	departmentList, err := service.DepartmentService.GetSubDepartments(departmentID)
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
Handle func for GET /entity/{entity_id}/department/{department_id}/user/list
*/
func (department *departmentApi) GetAllUsersUnderDepartment(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}

	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	departmentSuper := service.UserService.DepartmentSuper(ctx)
	if !entitySuper && !departmentSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	identity, err := service.DepartmentService.CheckDepartmentIdentity(ctx, entityID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if !identity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	userList, err := service.DepartmentService.GetAllUsers(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	userListRes := []define.DepartmentUserInfo{}
	err = copier.Copy(&userListRes, userList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	userListResponse := define.DepartmentUserListResponse{
		UserList: userListRes,
	}
	ctx.Success(userListResponse)
}

/*
Handle func for POST /entity/{entity_id}/department/{department_id}/user
*/
func (department *departmentApi) CreateUserInDepartment(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	if !isInEntity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}

	var createUserReq define.CreateDepartmentUserReq
	err = ctx.MustBindWith(&createUserReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	existsUser, err := service.UserService.ExistsUser(createUserReq.UserName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if existsUser {
		ctx.BadRequest(myerror.USER_HAS_EXISTED, myerror.USER_HAS_EXISTED_INFO)
		return
	}

	err = service.DepartmentService.CreateDepartmentUser(createUserReq, entityID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /entity/{entity_id}/department/{department_id}/manager
*/
func (department *departmentApi) SetManager(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	if !isInEntity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}

	var setManagerReq define.SetDepartmentManagerReq
	err = ctx.MustBindWith(&setManagerReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	thisUser, err := service.UserService.GetUserByName(setManagerReq.UserName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}
	if thisUser.EntityID != entityID {
		ctx.BadRequest(myerror.USER_NOT_IN_ENTITY, myerror.USER_NOT_IN_ENTITY_INFO)
		return
	}

	err = service.DepartmentService.SetDepartmentManager(setManagerReq.UserName, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for DELETE /entity/{entity_id}/department/{department_id}/manager/{user_id}
*/
func (department *departmentApi) DeleteDepartmentManager(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	if !isInEntity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}

	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}
	if thisUser.EntityID != entityID {
		ctx.BadRequest(myerror.USER_NOT_IN_ENTITY, myerror.USER_NOT_IN_ENTITY_INFO)
		return
	}
	if thisUser.DepartmentID != departmentID {
		ctx.BadRequest(myerror.USER_NOT_IN_DEPARTMENT, myerror.USER_NOT_IN_DEPARTMENT_INFO)
		return
	}

	err = service.DepartmentService.DeleteDepartmentManager(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for GET /entity/{entity_id}/department/{department_id}/manager
*/
func (department *departmentApi) GetDepartmentManager(ctx *utils.Context) {
	entityID, departmentID, err := department.GetTwoIDs(ctx)
	if err != nil {
		return
	}

	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	if !isInEntity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	isValid := department.CheckEntityDepartmentValid(ctx, entityID, departmentID)
	if !isValid {
		return
	}

	managerList, err := service.DepartmentService.GetDepartmentManagerList(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	managerListRes := []define.DepartmentManager{}
	err = copier.Copy(&managerListRes, managerList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	managerListResponse := define.DepartmentManagerListResponse{
		ManagerList: managerListRes,
	}

	ctx.Success(managerListResponse)
}
