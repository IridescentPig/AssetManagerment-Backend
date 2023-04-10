package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
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
		createDepartmentReq = define.CreateDepartmentReq{
			EntityID: entityID,
		}
	} else {
		departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
		if err != nil {
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
	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
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
