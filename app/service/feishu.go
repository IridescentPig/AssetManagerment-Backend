package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"context"
	"errors"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type feishuService struct{}

var FeishuService *feishuService

var appId = "cli_a4d26de640e4500c"

var appSecret = "qTGl1gT9HReTRxZAxAwAjewGlxeyZTfr"

var Client = lark.NewClient(appId, appSecret, lark.WithEnableTokenCache(true))

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

/*func (feishu *feishuService) FeishuLoginAndGetInfo(code string) (FeishuID string, AccessToken string, RefreshToken string, err error) {
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
}*/

func (feishu *feishuService) StoreToken(UserID uint, AccessToken string, RefreshToken string) error {
	return dao.UserDao.UpdateFeishuToken(UserID, AccessToken, RefreshToken)
}

func (feishu *feishuService) BindFeishu(UserID uint, FeishuID string) error {
	return dao.UserDao.BindFeishu(UserID, FeishuID)
}

func (feishu *feishuService) GetTokenByUserID(UserID uint) (access_token string, refresh_token string, err error) {
	user, err := dao.UserDao.GetUserByID(UserID)
	if err != nil {
		return
	}
	access_token = user.FeishuToken
	refresh_token = user.RefreshToken
	return
}

func (feishu *feishuService) TextTokenAndRefresh(UserID uint) (bool, error) {
	access_token, refresh_token, err := feishu.GetTokenByUserID(UserID)
	if err != nil {
		return false, err
	}
	_, err = feishu.GetUserInfo(access_token)
	if err == nil {
		return true, nil
	}
	resp, err := feishu.RefreshToken(refresh_token)
	if err != nil {
		return false, err
	}
	err = feishu.StoreToken(UserID, resp.Data.AccessToken, resp.Data.RefreshToken)
	if err != nil {
		return false, err
	}
	return true, err
}

func (feishu *feishuService) SendMessage(UserId uint, text string) error {
	user, err := dao.UserDao.GetUserByID(UserId)
	if err != nil {
		return err
	}

	text_content := fmt.Sprintf(`{"text":"%s"}`, text)

	req := larkim.NewCreateMessageReqBuilder().
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(user.FeishuID).
			MsgType(`text`).
			Content(text_content).
			Build()).
		Build()

	resp, err := Client.Im.Message.Create(context.Background(), req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return errors.New(resp.Msg)
	}

	return nil
}
