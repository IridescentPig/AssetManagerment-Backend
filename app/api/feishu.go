package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
	"log"
	"strconv"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"github.com/thoas/go-funk"
)

type feishuApi struct {
}

var FeishuApi *feishuApi

func newFeishuApi() *feishuApi {
	return &feishuApi{}
}

func init() {
	FeishuApi = newFeishuApi()
}

/*
Handle func for POST /user/feishu/login
*/
func (feishu *feishuApi) FeishuLogin(ctx *utils.Context) {
	var req define.FeishuBindOrLoginRequest

	err := ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	token_res, err := service.FeishuService.GetAccessToken(req.Code)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_FEISHU_CODE, myerror.INVALID_FEISHU_CODE_INFO)
		return
	}

	access_token := token_res.Data.AccessToken
	refresh_token := token_res.Data.RefreshToken

	info_res, err := service.FeishuService.GetUserInfo(access_token)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	feishu_id := info_res.Data.UserID

	user, err := service.FeishuService.FindUserByFeishuID(feishu_id)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if user == nil {
		ctx.BadRequest(myerror.FEISHU_NOT_BIND, myerror.FEISHU_NOT_BIND_INFO)
		return
	}

	err = service.FeishuService.StoreToken(user.ID, access_token, refresh_token)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	userBasicInfo := define.UserBasicInfo{
		UserID:          user.ID,
		UserName:        user.UserName,
		EntitySuper:     user.EntitySuper,
		DepartmentSuper: user.DepartmentSuper,
		SystemSuper:     user.SystemSuper,
		EntityID:        user.EntityID,
		DepartmentID:    user.DepartmentID,
	}
	token, err := utils.CreateToken(userBasicInfo)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	var userInfo define.UserInfo
	err = copier.Copy(&userInfo, user)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	data := define.UserLoginResponse{
		Token: token,
		User:  userInfo,
	}

	err = service.FeishuService.FeishuSync(userInfo.EntityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(data)
}

/*
Handle func for POST /user/feishu/bind
*/
func (feishu *feishuApi) FeishuBind(ctx *utils.Context) {
	var req define.FeishuBindOrLoginRequest

	err := ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	token_res, err := service.FeishuService.GetAccessToken(req.Code)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_FEISHU_CODE, myerror.INVALID_FEISHU_CODE_INFO)
		return
	}

	access_token := token_res.Data.AccessToken
	refresh_token := token_res.Data.RefreshToken

	info_res, err := service.FeishuService.GetUserInfo(access_token)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	current_user_id := UserApi.GetOperatorID(ctx)
	exist_user, err := service.FeishuService.FindUserByFeishuID(info_res.Data.UserID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if exist_user != nil && exist_user.ID != current_user_id {
		ctx.BadRequest(myerror.FEISHU_DUPLICATE_BIND, myerror.FEISHU_DUPLICATE_BIND_INFO)
		return
	}
	err = service.FeishuService.BindFeishu(current_user_id, info_res.Data.UserID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	err = service.FeishuService.StoreToken(current_user_id, access_token, refresh_token)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /user/feishu/bind
*/
func (feishu *feishuApi) FeishuCallBack(ctx *utils.Context) {
	// action_type, exists := ctx.Get("action_type")
	// if !exists {
	// 	ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
	// 	return
	// }

	var req define.FeishuCallBackReq
	err := ctx.MustBindWith(&req, binding.JSON)
	if err != nil {
		log.Println(err.Error())
		ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
		return
	}
	// log.Println(req)

	action_map := map[string]uint{
		"APPROVE": 1,
		"REJECT":  2,
	}

	matched_token := "sdjkljkx9lsadf110"
	// token, exists := ctx.Get("token")
	// if !exists || token != matched_token {
	// 	ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
	// 	return
	// }
	if req.Token != matched_token {
		log.Println("1")
		ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
		return
	}

	// task_id, exists := ctx.Get("instance_id")
	// if !exists {
	// 	ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
	// 	return
	// }
	instanceID, err := strconv.ParseUint(req.InstanceID, 10, 0)
	if err != nil {
		log.Println(err.Error())
		ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
		return
	}
	thisTask, err := service.TaskService.GetTaskInfoByID(uint(instanceID))
	if err != nil {
		ctx.InternalError("callback_error")
		return
	} else if thisTask == nil {
		ctx.InternalError("callback_error")
		return
	}
	// log.Println(instanceID)
	// if thisTask.State != 0 {
	// 	ctx.InternalError(myerror.TASK_NOT_PENDING_INFO)
	// 	return
	// }
	err = service.TaskService.ModifyTaskState(uint(instanceID), action_map[req.ActionType])
	thisTask.State = action_map[req.ActionType]
	if err != nil {
		ctx.InternalError("callback_error")
		return
	}

	thisUser, err := service.FeishuService.FindUserByFeishuID(req.UserID)
	if err != nil || thisUser == nil {
		ctx.InternalError("callback_error")
		return
	}

	assetIDs := funk.Map(thisTask.AssetList, func(thisAsset *model.Asset) uint {
		return thisAsset.ID
	}).([]uint)

	if req.ActionType == "APPROVE" {
		if thisTask.TaskType == 0 {
			err = service.AssetService.AcquireAssets(assetIDs, thisTask.UserID)
		} else if thisTask.TaskType == 1 {
			err = service.AssetService.CancelAssets(assetIDs, thisUser.ID)
		} else if thisTask.TaskType == 2 {
			err = service.AssetService.ModifyAssetMaintainerAndState(assetIDs, thisTask.TargetID)
		} else {
			err = service.AssetService.TransferAssets(assetIDs, thisTask.TargetID, thisTask.Target.DepartmentID, thisTask.DepartmentID)
		}

		if err != nil {
			ctx.InternalError("callback_error")
			return
		}
	}

	approvalCode, err := service.FeishuService.CreateApprovalDefination()
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	thisTask.State = action_map[req.ActionType]
	err = service.FeishuService.PutApproval(*thisTask, thisTask.User.FeishuID, approvalCode)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}
