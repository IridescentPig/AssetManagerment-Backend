package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"context"
	"errors"
	"fmt"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type feishuService struct{}

var FeishuService *feishuService

var appId = "cli_a4d26de640e4500c"

var appSecret = "qTGl1gT9HReTRxZAxAwAjewGlxeyZTfr"

var Client = lark.NewClient(appId, appSecret, lark.WithEnableTokenCache(true))

var approval_code string

func newFeishuService() *feishuService {
	return &feishuService{}
}

func init() {
	FeishuService = newFeishuService()
	err := FeishuService.CreateApprovalDefination()
	if err != nil {
		fmt.Printf("code Defination error: %s", err)
	}
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

// 审批相关
func (feishu *feishuService) CreateApprovalDefination() error {
	req := larkapproval.NewCreateExternalApprovalReqBuilder().
		DepartmentIdType(`open_department_id`).
		UserIdType("user_id").
		ExternalApproval(larkapproval.NewExternalApprovalBuilder().
			ApprovalName(`@i18n@1`).
			ApprovalCode(`AssetApproval`).
			GroupCode(`ApprovalRequest`).
			GroupName(`@i18n@2`).
			External(larkapproval.NewApprovalCreateExternalBuilder().
				CreateLinkMobile(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
				CreateLinkPc(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
				SupportPc(true).
				SupportMobile(true).
				SupportBatchRead(false).
				EnableMarkReaded(false).
				EnableQuickOperate(true).
				ActionCallbackUrl(`http://feishu.cn/approval/openapi/operate`). //记得改
				ActionCallbackToken(`sdjkljkx9lsadf110`).
				Build()).
			Viewers([]*larkapproval.ApprovalCreateViewers{
				larkapproval.NewApprovalCreateViewersBuilder().
					ViewerType(`TENANT`).
					Build(),
			}).
			I18nResources([]*larkapproval.I18nResource{
				larkapproval.NewI18nResourceBuilder().
					Locale(`zh-CN`).
					Texts([]*larkapproval.I18nResourceText{
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@1`).
							Value(`资产审批`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@2`).
							Value(`审批请求`).
							Build(),
					}).
					IsDefault(true).
					Build(),
			}).
			Managers([]string{`e2ba357b`}).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := Client.Approval.ExternalApproval.Create(context.Background(), req)

	approval_code = *resp.Data.ApprovalCode

	// 处理错误
	if err != nil {
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		return errors.New(resp.Msg)
	}
	return err
}

func (feishu *feishuService) PutApproval(task model.Task, FeishuID string) error {
	StateMap := map[uint]string{
		0: `PENDING`,
		1: `APPROVED`,
		2: `REJECTED`,
		3: `CANCELED`,
	}

	managers, err := dao.DepartmentDao.GetDepartmentManager(task.DepartmentID)
	if err != nil {
		return err
	}
	var TaskList []*larkapproval.ExternalInstanceTaskNode
	for index, manager := range managers {
		if len(manager.FeishuID) != 0 {
			TaskList = append(TaskList, larkapproval.NewExternalInstanceTaskNodeBuilder().
				TaskId((string)(index)).
				UserId(manager.FeishuID).
				Title(`同意`).
				Links(larkapproval.NewExternalInstanceLinkBuilder().
					PcLink(`http://`).
					MobileLink(`http://`).
					Build()).
				Status(`PENDING`).
				Extra(``).
				CreateTime(`1638468921000`).
				EndTime(`0`).
				UpdateTime(`1638468921000`).
				ActionContext(`123456`).
				ActionConfigs([]*larkapproval.ActionConfig{
					larkapproval.NewActionConfigBuilder().
						ActionType(`APPROVE`).
						ActionName(`@i18n@1`).
						IsNeedReason(true).
						IsReasonRequired(true).
						IsNeedAttachment(true).
						Build(),
				}).
				Build())
		}
	}

	req := larkapproval.NewCreateExternalInstanceReqBuilder().
		ExternalInstance(larkapproval.NewExternalInstanceBuilder().
			ApprovalCode(approval_code).
			Status(StateMap[task.State]).
			Extra(``).
			InstanceId((string)(task.ID)).
			Links(larkapproval.NewExternalInstanceLinkBuilder().
				PcLink(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
				MobileLink(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
				Build()).
			Form([]*larkapproval.ExternalInstanceForm{
				larkapproval.NewExternalInstanceFormBuilder().
					Name(`@i18n@1`).
					Value(`@i18n@2`).
					Build(),
			}).
			UserId(FeishuID).
			StartTime((string)(time.Now().UnixNano() / 1e6)).
			EndTime(`0`).
			UpdateTime((string)(time.Now().UnixNano() / 1e6)).
			UpdateMode(`REPLACE`).
			TaskList([]*larkapproval.ExternalInstanceTaskNode{
				larkapproval.NewExternalInstanceTaskNodeBuilder().
					TaskId((string)(task.ID)).
					UserId(`16fb9ff3`).
					Title(`同意`).
					Links(larkapproval.NewExternalInstanceLinkBuilder().
						PcLink(`http://`).
						MobileLink(`http://`).
						Build()).
					Status(`PENDING`).
					Extra(``).
					CreateTime(`1638468921000`).
					EndTime(`0`).
					UpdateTime(`1638468921000`).
					ActionContext(`123456`).
					ActionConfigs([]*larkapproval.ActionConfig{
						larkapproval.NewActionConfigBuilder().
							ActionType(`APPROVE`).
							ActionName(`@i18n@1`).
							IsNeedReason(true).
							IsReasonRequired(true).
							IsNeedAttachment(true).
							Build(),
					}).
					Build(),
			}).
			CcList([]*larkapproval.CcNode{
				larkapproval.NewCcNodeBuilder().
					CcId(`1231243`).
					UserId(`16fb9ff3`).
					OpenId(``).
					Links(larkapproval.NewExternalInstanceLinkBuilder().
						PcLink(`http://`).
						MobileLink(`http://`).
						Build()).
					ReadStatus(`READ`).
					Extra(``).
					Title(`XXX`).
					CreateTime(`1657093395000`).
					UpdateTime(`1657093395000`).
					Build(),
			}).
			I18nResources([]*larkapproval.I18nResource{
				larkapproval.NewI18nResourceBuilder().
					Locale(`zh-CN`).
					Texts([]*larkapproval.I18nResourceText{
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@1`).
							Value(`测试`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@2`).
							Value(`天`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@3`).
							Value(`2022-07-06`).
							Build(),
					}).
					IsDefault(true).
					Build(),
			}).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := Client.Approval.ExternalInstance.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		return errors.New(resp.Msg)
	}
	return err
}
