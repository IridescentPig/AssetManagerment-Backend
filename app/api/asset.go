package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
)

type assetApi struct {
}

var AssetApi *assetApi

func newAssetApi() *assetApi {
	return &assetApi{}
}

func init() {
	AssetApi = newAssetApi()
}

/*
Handle func for GET /department/{department_id}/asset/list
*/
func (asset *assetApi) GetAssetList(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetTree, err := service.AssetService.GetSubAsset(0, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	assetListRespone := define.AssetListResponse{
		AssetList: assetTree,
	}

	ctx.Success(assetListRespone)
}

/*
Handle func for PATCH /department/{department_id}/asset/{asset_id}
*/
func (asset *assetApi) ModifyAssetInfo(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetID, err := service.EntityService.GetParamID(ctx, "asset_id")
	if err != nil {
		return
	}

	thisAsset, err := service.AssetService.GetAssetByID(assetID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisAsset == nil {
		ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
		return
	} else if thisAsset.DepartmentID != departmentID {
		ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_CLASS_NOT_FOUND_INFO)
		return
	}

	var modifyAssetReq define.ModifyAssetInfoReq
	err = ctx.MustBindWith(&modifyAssetReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	} else if modifyAssetReq.Type < 0 || modifyAssetReq.Type > 2 {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	if modifyAssetReq.ParentID != nil && *modifyAssetReq.ParentID != 0 {
		exists, err := service.AssetService.ExistAsset(*modifyAssetReq.ParentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if !exists {
			ctx.BadRequest(myerror.PARENT_ASSET_NOT_FOUND, myerror.PARENT_ASSET_NOT_FOUND_INFO)
			return
		}
		isAncestor, err := service.AssetService.CheckIsAncestor(thisAsset.ID, *modifyAssetReq.ParentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if isAncestor {
			ctx.BadRequest(myerror.PARENT_CANNOOT_BE_SUCCESSOR, myerror.PARENT_CANNOOT_BE_SUCCESSOR_INFO)
			return
		}
	}

	err = service.AssetService.ModifyAssetInfo(assetID, modifyAssetReq)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /department/{department_id}/asset
*/
func (asset *assetApi) CreateAssets(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	userID := UserApi.GetOperatorID(ctx)

	var assetsCreateReq define.CreateAssetListReq
	err = ctx.MustBindWith(&assetsCreateReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	for _, asset := range assetsCreateReq.AssetList {
		exists, err := service.AssetClassService.ExistsAssetClass(asset.ClassID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		} else if !exists {
			ctx.BadRequest(myerror.ASSET_CLASS_NOT_FOUND, myerror.ASSET_CLASS_NOT_FOUND_INFO)
			return
		}
		err = service.AssetService.CreateAsset(&asset, departmentID, 0, userID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	}

	ctx.Success(nil)
}

/*
Handle func for PATCH /department/{department_id}/asset/expire
*/
func (asset *assetApi) ExpireAsset(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var expireReq define.ExpireAssetListReq
	err = ctx.MustBindWith(&expireReq, binding.JSON)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	assetIDs := []uint{}
	for _, assetID := range expireReq.ExpireList {
		exists, err := service.AssetService.ExistAsset(assetID.AssetID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		} else if !exists {
			ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
			return
		}
		isInDepartment, err := service.AssetService.CheckAssetInDepartment(assetID.AssetID, departmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if !isInDepartment {
			ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_NOT_IN_DEPARTMENT_INFO)
			return
		}
		assetIDs = append(assetIDs, assetID.AssetID)
	}

	err = service.AssetService.ExpireAssets(assetIDs)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /department/{department_id}/asset/transfer
*/
func (asset *assetApi) TransferAssets(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	var transferReq define.AssetTransferReq
	err = ctx.MustBindWith(&transferReq, binding.JSON)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	targetUser, err := service.UserService.GetUserByID(transferReq.UserID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if targetUser == nil {
		ctx.BadRequest(myerror.TARGET_USER_NOT_FOUND, myerror.TARGET_USER_NOT_FOUND_INFO)
		return
	} else if targetUser.EntityID != thisUser.EntityID {
		ctx.BadRequest(myerror.NOT_IN_SAME_ENTITY, myerror.NOT_IN_SAME_ENTITY_INFO)
		return
	} else if !targetUser.DepartmentSuper {
		ctx.BadRequest(myerror.TARGET_NOT_DEPARTMENT_SUPER, myerror.TARGET_NOT_DEPARTMENT_SUPER_INFO)
		return
	}

	assetIDs := []uint{}

	for _, assetID := range transferReq.Assets {
		exists, err := service.AssetService.ExistAsset(assetID.AssetID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		} else if !exists {
			ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
			return
		}
		isInDepartment, err := service.AssetService.CheckAssetInDepartment(assetID.AssetID, departmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
		if !isInDepartment {
			ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_NOT_IN_DEPARTMENT_INFO)
			return
		}
		assetIDs = append(assetIDs, assetID.AssetID)
	}

	err = service.AssetService.TransferAssets(assetIDs, targetUser.ID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}
