package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
	"errors"

	"github.com/gin-gonic/gin/binding"
)

type assetClassApi struct {
}

var AssetClassApi *assetClassApi

func newAssetClassApi() *assetClassApi {
	return &assetClassApi{}
}

func init() {
	AssetClassApi = newAssetClassApi()
}

/*
 */
func (assetClass *assetClassApi) CheckAssetIdentity(ctx *utils.Context) (bool, uint, error) {
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return false, departmentID, err
	}
	existsDepartment, err := service.DepartmentService.ExistsDepartmentByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false, departmentID, err
	} else if !existsDepartment {
		ctx.NotFound(myerror.DEPARTMENT_NOT_FOUND, myerror.DEPARTMENT_NOT_FOUND_INFO)
		return false, departmentID, errors.New("")
	}
	isDepartmentSuper := service.UserService.DepartmentSuper(ctx)
	if !isDepartmentSuper {
		return false, departmentID, nil
	}
	isInDepartment := service.DepartmentService.CheckIsInDepartment(ctx, departmentID)
	return isInDepartment, departmentID, nil
}

/*
 */
func (assetClass *assetClassApi) CheckAssetIdentityReturnDepartment(ctx *utils.Context) (bool, *model.Department, error) {
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return false, nil, err
	}
	thisDepartment, err := service.DepartmentService.GetDepartmentInfoByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false, nil, err
	} else if thisDepartment == nil {
		ctx.NotFound(myerror.DEPARTMENT_NOT_FOUND, myerror.DEPARTMENT_NOT_FOUND_INFO)
		return false, nil, errors.New("")
	}
	isDepartmentSuper := service.UserService.DepartmentSuper(ctx)
	if !isDepartmentSuper {
		return false, nil, nil
	}
	isInDepartment := service.DepartmentService.CheckIsInDepartment(ctx, departmentID)
	return isInDepartment, thisDepartment, nil
}

func (assetClass *assetClassApi) CheckAssetViewIdentity(ctx *utils.Context) (bool, uint, error) {
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return false, departmentID, err
	}
	existsDepartment, err := service.DepartmentService.ExistsDepartmentByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return false, departmentID, err
	} else if !existsDepartment {
		ctx.NotFound(myerror.DEPARTMENT_NOT_FOUND, myerror.DEPARTMENT_NOT_FOUND_INFO)
		return false, departmentID, errors.New("")
	}
	isInDepartment := service.DepartmentService.CheckIsInDepartment(ctx, departmentID)
	return isInDepartment, departmentID, nil
}

/*
Handle func for POST /department/{department_id}/asset_class
*/
func (assetClass *assetClassApi) CreateAssetClass(ctx *utils.Context) {
	hasIdentity, departmentID, err := assetClass.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var createAssetClassReq define.CreateAssetClassReq
	err = ctx.MustBindWith(&createAssetClassReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}
	if createAssetClassReq.Type > 2 || createAssetClassReq.Type < 0 {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	if createAssetClassReq.ParentID != 0 {
		existsParentClass, err := service.AssetClassService.ExistsAssetClass(createAssetClassReq.ParentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if !existsParentClass {
			ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
			return
		}
	}

	err = service.AssetClassService.CreateAssetClass(createAssetClassReq, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handler func for GET /department/{department_id}/asset_class
*/
func (assetClass *assetClassApi) GetAssetClassTree(ctx *utils.Context) {
	hasIdentity, departmentID, err := assetClass.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetClassTree, err := service.AssetClassService.GetSubAssetClass(0, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	assetClassTreeResponse := define.AssetClassTreeResponse{
		AssetClassTree: assetClassTree,
	}

	ctx.Success(assetClassTreeResponse)
}

/*
Handle func for PATCH /department/{department_id}/asset_class/{class_id}
*/
func (assetClass *assetClassApi) ModifyAssetClassInfo(ctx *utils.Context) {
	hasIdentity, _, err := assetClass.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var modifyAssetClassReq define.ModifyAssetClassReq
	err = ctx.MustBindWith(&modifyAssetClassReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}
	if modifyAssetClassReq.Type > 2 || modifyAssetClassReq.Type < 0 {
		ctx.BadRequest(myerror.INVALID_TYPE_OF_CLASS, myerror.INVALID_TYPE_OF_CLASS_INFO)
		return
	}

	classID, err := service.EntityService.GetParamID(ctx, "class_id")
	if err != nil {
		return
	}

	existAssetClass, err := service.AssetClassService.ExistsAssetClass(classID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	if !existAssetClass {
		ctx.BadRequest(myerror.ASSET_CLASS_NOT_FOUND, myerror.ASSET_CLASS_NOT_FOUND_INFO)
		return
	}
	if modifyAssetClassReq.ParentID != nil && *modifyAssetClassReq.ParentID != 0 {
		existAssetClass, err = service.AssetClassService.ExistsAssetClass(*modifyAssetClassReq.ParentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}

		if !existAssetClass {
			ctx.BadRequest(myerror.PARENT_ASSET_CLASS_NOT_FOUND, myerror.PARENT_ASSET_CLASS_NOT_FOUND_INFO)
			return
		}

		isAncestor, err := service.AssetClassService.CheckIsAncestor(classID, *modifyAssetClassReq.ParentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if isAncestor {
			ctx.BadRequest(myerror.PARENT_CANNOOT_BE_SUCCESSOR, myerror.PARENT_CANNOOT_BE_SUCCESSOR_INFO)
			return
		}
	}

	err = service.AssetClassService.ModifyAssetClassInfo(modifyAssetClassReq, classID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for DELETE /department/:department_id/asset_class/:class_id
*/
func (assetClass *assetClassApi) DeleteAssetClass(ctx *utils.Context) {
	hasIdentity, departmentID, err := assetClass.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	classID, err := service.EntityService.GetParamID(ctx, "class_id")
	if err != nil {
		return
	}

	existAssetClass, err := service.AssetClassService.ExistsAssetClass(classID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if !existAssetClass {
		ctx.BadRequest(myerror.ASSET_CLASS_NOT_FOUND, myerror.ASSET_CLASS_NOT_FOUND_INFO)
		return
	}

	hasAsset, err := service.AssetClassService.ClassHasAsset(classID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if hasAsset {
		ctx.BadRequest(myerror.CLASS_HAS_ASSET, myerror.CLASS_HAS_ASSET_INFO)
		return
	}

	hasSubClass, err := service.AssetClassService.ClassHasSubClass(classID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
	} else if hasSubClass {
		ctx.BadRequest(myerror.ClASS_HAS_SUB_CLASS, myerror.CLASS_HAS_SUB_CLASS_INFO)
		return
	}

	err = service.AssetClassService.DeleteAssetClass(classID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for /department/:department_id/asset_class/:class_id
*/
func (assetClass *assetClassApi) GetSubAssetClass(ctx *utils.Context) {
	hasIdentity, departmentID, err := assetClass.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	classID, err := service.EntityService.GetParamID(ctx, "class_id")
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	assetClassTree, err := service.AssetClassService.GetSubAssetClass(classID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	assetClassTreeResponse := define.AssetClassTreeResponse{
		AssetClassTree: assetClassTree,
	}

	ctx.Success(assetClassTreeResponse)
}
