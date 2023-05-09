package service

import (
	"context"
	"errors"

	lark "github.com/larksuite/oapi-sdk-go/v3"
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

func (feishu *feishuService) GetAccessToken(code string) (err error) {
	res, err := Client.Ext.Authen.AuthenAccessToken(context.Background(),
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

func (feishu *feishuService) Login() {

}
