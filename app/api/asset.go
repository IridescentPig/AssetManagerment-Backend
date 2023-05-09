package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
	"sort"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
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

func (asset *assetApi) CheckAssetExistsAndValid(ctx *utils.Context, departmentID uint) (uint, *model.Asset, bool) {
	assetID, err := service.EntityService.GetParamID(ctx, "asset_id")
	if err != nil {
		return 0, nil, false
	}

	thisAsset, err := service.AssetService.GetAssetByID(assetID)
	if err != nil {
		ctx.InternalError(err.Error())
		return 0, nil, false
	} else if thisAsset == nil {
		ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
		return 0, nil, false
	} else if thisAsset.DepartmentID != departmentID {
		ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_CLASS_NOT_FOUND_INFO)
		return 0, nil, false
	}

	return assetID, thisAsset, true
}

func (asset *assetApi) CheckAssetsValid(ctx *utils.Context, departmentID uint, assetList []define.ExpireAssetReq) ([]uint, bool) {
	assetIDs := []uint{}
	for _, assetID := range assetList {
		exists, err := service.AssetService.ExistAsset(assetID.AssetID)
		if err != nil {
			ctx.InternalError(err.Error())
			return nil, false
		} else if !exists {
			ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
			return nil, false
		}
		isInDepartment, err := service.AssetService.CheckAssetInDepartment(assetID.AssetID, departmentID)
		if err != nil {
			ctx.InternalError(err.Error())
			return nil, false
		}
		if !isInDepartment {
			ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_NOT_IN_DEPARTMENT_INFO)
			return nil, false
		}
		assetIDs = append(assetIDs, assetID.AssetID)
	}

	return assetIDs, true
}

/*
Handle func for GET /department/{department_id}/asset/list
*/
func (asset *assetApi) GetAssetList(ctx *utils.Context) {
	// hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	// if err != nil {
	// 	return
	// } else if !hasIdentity {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	departmentID, err := service.EntityService.GetParamID(ctx, "department_id")
	if err != nil {
		return
	}

	existsDepartment, err := service.DepartmentService.ExistsDepartmentByID(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !existsDepartment {
		ctx.NotFound(myerror.DEPARTMENT_NOT_FOUND, myerror.DEPARTMENT_NOT_FOUND_INFO)
		return
	}

	thisUser := UserApi.GetOperatorInfo(ctx)
	if thisUser.DepartmentID != departmentID {
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

	// assetID, err := service.EntityService.GetParamID(ctx, "asset_id")
	// if err != nil {
	// 	return
	// }

	// thisAsset, err := service.AssetService.GetAssetByID(assetID)
	// if err != nil {
	// 	ctx.InternalError(err.Error())
	// 	return
	// } else if thisAsset == nil {
	// 	ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
	// 	return
	// } else if thisAsset.DepartmentID != departmentID {
	// 	ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_CLASS_NOT_FOUND_INFO)
	// 	return
	// }

	assetID, thisAsset, isOK := asset.CheckAssetExistsAndValid(ctx, departmentID)
	if !isOK {
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
	minimalPrice := decimal.NewFromFloat(0)
	maxiumPrice, _ := decimal.NewFromString("99999999.99")
	if minimalPrice.Cmp(modifyAssetReq.Price) == 1 || maxiumPrice.Cmp(modifyAssetReq.Price) == -1 {
		ctx.BadRequest(myerror.PRICE_OUT_OF_RANGE, myerror.PRICE_OUT_OF_RANGE_INFO)
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

	if modifyAssetReq.Price != decimal.Zero || modifyAssetReq.Expire != 0 {
		_ = service.AssetService.UpdateNetWorth(assetID)
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

	minimalPrice := decimal.NewFromFloat(0)
	maxiumPrice, _ := decimal.NewFromString("99999999.99")

	for _, asset := range assetsCreateReq.AssetList {
		exists, err := service.AssetClassService.ExistsAssetClass(asset.ClassID)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		} else if !exists {
			ctx.BadRequest(myerror.ASSET_CLASS_NOT_FOUND, myerror.ASSET_CLASS_NOT_FOUND_INFO)
			return
		}

		if minimalPrice.Cmp(asset.Price) == 1 || maxiumPrice.Cmp(asset.Price) == -1 {
			ctx.BadRequest(myerror.PRICE_OUT_OF_RANGE, myerror.PRICE_OUT_OF_RANGE_INFO)
			return
		}

		if asset.ParentID != 0 {
			exists, err := service.AssetService.ExistAsset(asset.ParentID)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			} else if !exists {
				ctx.BadRequest(myerror.PARENT_ASSET_NOT_FOUND, myerror.PARENT_ASSET_NOT_FOUND_INFO)
				return
			}
		}

		err = service.AssetService.CreateAsset(&asset, departmentID, asset.ParentID, userID)
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

	// assetIDs := []uint{}
	// for _, assetID := range expireReq.ExpireList {
	// 	exists, err := service.AssetService.ExistAsset(assetID.AssetID)
	// 	if err != nil {
	// 		ctx.InternalError(err.Error())
	// 		return
	// 	} else if !exists {
	// 		ctx.BadRequest(myerror.ASSET_NOT_FOUND, myerror.ASSET_NOT_FOUND_INFO)
	// 		return
	// 	}
	// 	isInDepartment, err := service.AssetService.CheckAssetInDepartment(assetID.AssetID, departmentID)
	// 	if err != nil {
	// 		ctx.InternalError(err.Error())
	// 		return
	// 	}
	// 	if !isInDepartment {
	// 		ctx.BadRequest(myerror.ASSET_NOT_IN_DEPARTMENT, myerror.ASSET_NOT_IN_DEPARTMENT_INFO)
	// 		return
	// 	}
	// 	assetIDs = append(assetIDs, assetID.AssetID)
	// }

	assetIDs, isOK := asset.CheckAssetsValid(ctx, departmentID, expireReq.ExpireList)
	if !isOK {
		return
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

	assetIDs, isOK := asset.CheckAssetsValid(ctx, departmentID, transferReq.Assets)
	if !isOK {
		return
	}

	err = service.AssetService.TransferAssets(assetIDs, targetUser.ID, targetUser.DepartmentID, departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

func (asset *assetApi) userAssetPrevillige(ctx *utils.Context) (*model.User, bool) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return nil, false
	}

	operatorUser := UserApi.GetOperatorInfo(ctx)
	if operatorUser.UserID != userID {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return nil, false
	}

	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return nil, false
	} else if thisUser == nil {
		ctx.NotFound(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return nil, false
	}

	if thisUser.EntityID == 0 {
		ctx.BadRequest(myerror.USER_NOT_IN_ENTITY, myerror.USER_NOT_IN_ENTITY_INFO)
		return nil, false
	} else if thisUser.DepartmentID == 0 {
		ctx.BadRequest(myerror.USER_NOT_IN_DEPARTMENT, myerror.USER_NOT_IN_DEPARTMENT_INFO)
		return nil, false
	}

	return thisUser, true
}

/*
Handle func for GET /users/:user_id/assets/maintain
*/
func (asset *assetApi) GetUserMaintainAssets(ctx *utils.Context) {
	// userID, err := service.EntityService.GetParamID(ctx, "user_id")
	// if err != nil {
	// 	return
	// }

	// operatorUser := UserApi.GetOperatorInfo(ctx)
	// if operatorUser.UserID != userID {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }

	// thisUser, err := service.UserService.GetUserByID(userID)
	// if err != nil {
	// 	ctx.InternalError(err.Error())
	// 	return
	// } else if thisUser == nil {
	// 	ctx.NotFound(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
	// 	return
	// }

	// if thisUser.EntityID == 0 {
	// 	ctx.BadRequest(myerror.USER_NOT_IN_ENTITY, myerror.USER_NOT_IN_ENTITY_INFO)
	// 	return
	// } else if thisUser.DepartmentID == 0 {
	// 	ctx.BadRequest(myerror.USER_NOT_IN_DEPARTMENT, myerror.USER_NOT_IN_DEPARTMENT_INFO)
	// 	return
	// }

	thisUser, isOK := asset.userAssetPrevillige(ctx)
	if !isOK {
		return
	}

	var assetListRes []*define.AssetInfo

	assetList, err := service.AssetService.GetUserMaintainAssets(thisUser.ID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	err = copier.Copy(&assetListRes, assetList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	for _, thisAsset := range assetListRes {
		thisAsset.Children = nil
	}

	assetListResponse := define.AssetListResponse{
		AssetList: assetListRes,
	}

	ctx.Success(assetListResponse)
}

/*
Handle func for POST /users/:user_id/assets/:asset_id/maintain
*/
func (asset *assetApi) FinishMaintenance(ctx *utils.Context) {
	thisUser, isOK := asset.userAssetPrevillige(ctx)
	if !isOK {
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
	} else if thisAsset.State != 2 {
		ctx.BadRequest(myerror.ASSET_NOT_IN_MAINTENCE, myerror.ASSET_NOT_IN_MAINTENCE_INFO)
		return
	} else if thisAsset.MaintainerID != thisUser.ID {
		ctx.BadRequest(myerror.NOT_YOUR_MAINTENCE_ASSET, myerror.NOT_YOUR_MAINTENCE_ASSET_INFO)
		return
	}

	err = service.AssetService.ModifyAssetMaintainerAndState([]uint{assetID}, 0)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /department/{department_id}/asset/{asset_id}/property
*/
func (asset *assetApi) CreateAssetProperty(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetID, _, isOK := asset.CheckAssetExistsAndValid(ctx, departmentID)
	if !isOK {
		return
	}

	var createPropertyReq define.AssetPropertyReq
	err = ctx.MustBindWith(&createPropertyReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	isExist, err := service.AssetService.ExistsProperty(assetID, createPropertyReq.Key)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if isExist {
		ctx.BadRequest(myerror.PROPERTY_HAS_EXIST, myerror.PROPERTY_HAS_EXIST_INFO)
		return
	}

	err = service.AssetService.SetProperty(assetID, createPropertyReq.Key, createPropertyReq.Value)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for PATCH /department/{department_id}/asset/{asset_id}/property
*/
func (asset *assetApi) ModifyAssetProperty(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetID, _, isOK := asset.CheckAssetExistsAndValid(ctx, departmentID)
	if !isOK {
		return
	}

	var modifyPropertyReq define.AssetPropertyReq
	err = ctx.MustBindWith(&modifyPropertyReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	isExist, err := service.AssetService.ExistsProperty(assetID, modifyPropertyReq.Key)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !isExist {
		ctx.BadRequest(myerror.PROPERTY_NOT_EXIST, myerror.PROPERTY_NOT_EXIST_INFO)
		return
	}

	err = service.AssetService.SetProperty(assetID, modifyPropertyReq.Key, modifyPropertyReq.Value)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for DELETE /department/{department_id}/asset/{asset_id}/property
*/
func (asset *assetApi) DeleteAssetProperty(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetID, _, isOK := asset.CheckAssetExistsAndValid(ctx, departmentID)
	if !isOK {
		return
	}

	var deletePropertyReq define.DeleteAssetPropertyReq
	err = ctx.MustBindWith(&deletePropertyReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	isExist, err := service.AssetService.ExistsProperty(assetID, deletePropertyReq.Key)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !isExist {
		ctx.BadRequest(myerror.PROPERTY_NOT_EXIST, myerror.PROPERTY_NOT_EXIST_INFO)
		return
	}

	err = service.AssetService.DeleteProperty(assetID, deletePropertyReq.Key)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for GET /department/:department_id/asset/:asset_id/history
*/
func (asset *assetApi) GetAssetHistory(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetViewIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	assetID, _, isOK := asset.CheckAssetExistsAndValid(ctx, departmentID)
	if !isOK {
		return
	}

	approvedTaskList, err := service.AssetService.GetAssetHistory(assetID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	assetHistory := funk.Map(approvedTaskList, func(task *model.Task) *define.AssetHistory {
		history := define.AssetHistory{
			Type:         task.TaskType,
			ReviewTime:   task.ReviewAt,
			UserID:       task.UserID,
			Username:     task.User.UserName,
			DepartmentID: task.DepartmentID,
		}
		if task.TaskType >= 2 {
			history.TargetID = task.TargetID
			history.TargetName = task.Target.UserName
			history.TargetDepartmentID = task.Target.DepartmentID
		}

		return &history
	}).([]*define.AssetHistory)

	sort.Slice(assetHistory, func(i, j int) bool {
		return time.Time(*assetHistory[i].ReviewTime).After(time.Time(*assetHistory[j].ReviewTime))
	})

	assetHistoryRes := define.AssetHistoryResponse{
		History: assetHistory,
	}

	ctx.Success(assetHistoryRes)
}
