package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"context"
	"errors"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
)

type feishuService struct{}

var FeishuService *feishuService

var appId = "cli_a4d26de640e4500c"

var appSecret = "qTGl1gT9HReTRxZAxAwAjewGlxeyZTfr"

var Client = lark.NewClient(appId, appSecret)

func newFeishuService() *feishuService {
	return &feishuService{}
}

func init() {
	FeishuService = newFeishuService()
}

func (feishu *feishuService) GetAccessToken(code string) (res *larkext.AuthenAccessTokenResp, err error) {
	res, err = Client.Ext.Authen.AuthenAccessToken(context.Background(),
		larkext.NewAuthenAccessTokenReqBuilder().
			Body(larkext.NewAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeAuthorizationCode).
				Code(code).
				Build()).
			Build())
	if err != nil {
		return
	}
	if !res.Success() {
		err = errors.New(res.Msg)
	}
	return
}

func (feishu *feishuService) GetUserInfo(token string) (res *larkext.AuthenUserInfoResp, err error) {
	res, err = Client.Ext.Authen.AuthenUserInfo(context.Background(),
		larkcore.WithUserAccessToken(token))
	if err != nil {
		return
	}
	if !res.Success() {
		err = errors.New(res.Msg)
	}
	return
}

func (feishu *feishuService) RefreshToken(token string) (res *larkext.RefreshAuthenAccessTokenResp, err error) {
	res, err = Client.Ext.Authen.RefreshAuthenAccessToken(context.Background(),
		larkext.NewRefreshAuthenAccessTokenReqBuilder().
			Body(larkext.NewRefreshAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeRefreshCode).
				RefreshToken(token).Build()).Build())
	if err != nil {
		return
	}
	if !res.Success() {
		err = errors.New(res.Msg)
	}
	return
}

func (feishu *feishuService) FindUserByFeishuID(FeishuID string) (user *model.User, err error) {
	user, err = dao.UserDao.GetUserByFeishuID(FeishuID)
	return
}

func (feishu *feishuService) FeishuLoginAndGetInfo(code string) (FeishuID string, AccessToken string, RefreshToken string, err error) {
	res, err := feishu.GetAccessToken(code)
	if err != nil {
		return
	}
	AccessToken = res.Data.AccessToken
	RefreshToken = res.Data.RefreshToken
	info_res, err := feishu.GetUserInfo(AccessToken)
	if err != nil {
		return
	}
	FeishuID = info_res.Data.UserID
	return
}

func (feishu *feishuService) StoreToken(UserID uint, AccessToken string, RefreshToken string) error {
	return dao.UserDao.UpdateFeishuToken(UserID, AccessToken, RefreshToken)
}

func (feishu *feishuService) BindFeishu(UserID uint, FeishuID string) error {
	return dao.UserDao.BindFeishu(UserID, FeishuID)
}
