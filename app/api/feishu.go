package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
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
	action_type, exists := ctx.Get("action_type")
	if !exists {
		ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
		return
	}
	action_map := map[string]uint{
		"APPROVE": 1,
		"REJECT":  2,
	}

	matched_token := "sdjkljkx9lsadf110"
	token, exists := ctx.Get("token")
	if !exists || token != matched_token {
		ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
		return
	}

	task_id, exists := ctx.Get("instance_id")
	if !exists {
		ctx.BadRequest(myerror.FEISHU_CALLBACK_ERROR, myerror.FEISHU_CALLBACK_ERROR_INFO)
		return
	}

	err := service.TaskService.ModifyTaskState(task_id.(uint), action_map[action_type.(string)])
	if err != nil {
		ctx.InternalError("callback_error")
	}

	ctx.Success(nil)
}
