package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
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

	feishu_id := info_res.Data.OpenID

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

	ctx.Success(nil)
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
	err = service.FeishuService.BindFeishu(current_user_id, info_res.Data.OpenID)
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
