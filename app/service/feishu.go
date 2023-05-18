package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
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

	req := larkim.NewCreateMessageReqBuilder().ReceiveIdType("user_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(user.FeishuID).
			MsgType(`text`).
			Content(text_content).
			Build()).
		Build()

	resp, err := Client.Im.Message.Create(context.Background(), req)
	if err != nil {
		// log.Println(err.Error())
		return err
	}

	if !resp.Success() {
		// log.Println(resp.Code)
		// log.Println(resp.Err.Details)
		return errors.New(resp.Msg)
	}

	return nil
}

//var CallBackUrl = "http://0.0.0.0:8070/user/feishu/callback"

// 审批相关
func (feishu *feishuService) CreateApprovalDefination() (approval_code string, err error) {
	var CallBackUrl string
	if gin.Mode() == gin.DebugMode {
		// CallBackUrl = "http://AssetManagement-Backend-dev-BinaryAbstract.app.secoder.net/user/feishu/callback"
		CallBackUrl = "http://49.233.51.221:8080/user/feishu/callback"
	} else {
		CallBackUrl = "http://AssetManagement-Backend-BinaryAbstract.app.secoder.net/user/feishu/callback"
	}
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
				ActionCallbackUrl(CallBackUrl). //记得改
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

	// 处理错误
	if err != nil {
		return
	}

	approval_code = *resp.Data.ApprovalCode

	// 服务端错误处理
	if !resp.Success() {
		err = errors.New(resp.Msg)
	}
	return
}

func (feishu *feishuService) PutApproval(task model.Task, FeishuID string, approval_code string) error {
	//log.Print("test approval:", strconv.FormatInt(int64(task.ID), 10))
	StateMap := map[uint]string{
		0: `PENDING`,
		1: `APPROVED`,
		2: `REJECTED`,
		3: `REJECTED`,
	}

	managers, err := dao.DepartmentDao.GetDepartmentManager(task.DepartmentID)
	if err != nil {
		return err
	}
	var TaskList []*larkapproval.ExternalInstanceTaskNode
	for index, manager := range managers {
		if len(manager.FeishuID) != 0 {
			TaskList = append(TaskList, larkapproval.NewExternalInstanceTaskNodeBuilder().
				TaskId(strconv.FormatInt(int64(task.ID), 10)+"&"+strconv.FormatInt(int64(index), 10)).
				UserId(manager.FeishuID).
				Title(task.TaskDescription).
				Links(larkapproval.NewExternalInstanceLinkBuilder().
					PcLink(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
					MobileLink(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
					Build()).
				Status(StateMap[task.State]).
				Extra(``).
				CreateTime(strconv.FormatInt(time.Now().UnixMilli(), 10)). //改时间戳
				EndTime(`0`).
				UpdateTime(strconv.FormatInt(time.Now().UnixMilli(), 10)).
				ActionConfigs([]*larkapproval.ActionConfig{
					larkapproval.NewActionConfigBuilder().
						ActionType(`APPROVE`).
						ActionName(`@i18n@7`).
						Build(),
					larkapproval.NewActionConfigBuilder().
						ActionType(`REJECT`).
						ActionName(`@i18n@8`).
						Build(),
				}).
				Build())
		}
	}

	TaskTypeMap := map[uint]string{
		0: "领用",
		1: "退库",
		2: "维保",
		3: "转移",
	}

	req := larkapproval.NewCreateExternalInstanceReqBuilder().
		ExternalInstance(larkapproval.NewExternalInstanceBuilder().
			ApprovalCode(approval_code).
			Status(StateMap[task.State]).
			Extra(``).
			InstanceId(strconv.FormatInt(int64(task.ID), 10)).
			Links(larkapproval.NewExternalInstanceLinkBuilder().
				PcLink(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
				MobileLink(`http://assetmanagement-frontend-binaryabstract.app.secoder.net/#/asset/list`).
				Build()).
			Form([]*larkapproval.ExternalInstanceForm{
				larkapproval.NewExternalInstanceFormBuilder().
					Name(`@i18n@1`).
					Value(`@i18n@2`).
					Build(),
				larkapproval.NewExternalInstanceFormBuilder().
					Name(`@i18n@3`).
					Value(`@i18n@4`).
					Build(),
				larkapproval.NewExternalInstanceFormBuilder().
					Name(`@i18n@5`).
					Value(`@i18n@6`).
					Build(),
			}).
			UserId(FeishuID).
			StartTime(strconv.FormatInt(time.Now().UnixMilli(), 10)). //改时间戳
			EndTime(`0`).
			UpdateTime(strconv.FormatInt(time.Now().UnixMilli(), 10)).
			UpdateMode(`REPLACE`).
			TaskList(TaskList).
			CcList([]*larkapproval.CcNode{}).
			I18nResources([]*larkapproval.I18nResource{
				larkapproval.NewI18nResourceBuilder().
					Locale(`zh-CN`).
					Texts([]*larkapproval.I18nResourceText{
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@1`).
							Value(`任务类型`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@2`).
							Value(TaskTypeMap[task.TaskType]).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@3`).
							Value(`发起者`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@4`).
							Value(task.User.UserName).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@5`).
							Value(`任务描述`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@6`).
							Value(task.TaskDescription).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@7`).
							Value(`同意`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@8`).
							Value(`拒绝`).
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

func (feishu *feishuService) GetSubDEpartments(PageToken string) (departments []*string, err error) {
	req := larkcontact.NewChildrenDepartmentReqBuilder().
		UserIdType(`user_id`).
		DepartmentIdType(`open_department_id`).
		DepartmentId(strconv.FormatInt(0, 10)).
		PageSize(50).
		PageToken(PageToken).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := Client.Contact.Department.Children(context.Background(), req)

	// 处理错误
	if err != nil {
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		err = errors.New(resp.Msg)
		return
	}

	for _, department := range resp.Data.Items {
		departments = append(departments, department.DepartmentId)
	}

	if *resp.Data.HasMore {
		next_deparments, in_err := feishu.GetSubDEpartments(*resp.Data.PageToken)
		if in_err != nil {
			err = in_err
			return
		}
		departments = append(departments, next_deparments...)
	}

	return
}

func (feishu *feishuService) GetDepartmentUser(DepartmentID string, PageToken string) (feishuIDs []*larkcontact.User, err error) {
	req := larkcontact.NewFindByDepartmentUserReqBuilder().
		UserIdType(`user_id`).
		DepartmentIdType(`open_department_id`).
		DepartmentId(DepartmentID).
		PageSize(50).
		PageToken(PageToken).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := Client.Contact.User.FindByDepartment(context.Background(), req, larkcore.WithTenantAccessToken("t-g1045beaMUXCE5AKFM5T2GNNVOM5BDZESUIMVTZT"))

	// 处理错误
	if err != nil {
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		err = errors.New(resp.Msg)
		return
	}

	feishuIDs = append(feishuIDs, resp.Data.Items...)

	if *resp.Data.HasMore {
		next_users, in_err := feishu.GetDepartmentUser(DepartmentID, *resp.Data.PageToken)
		if in_err != nil {
			err = in_err
			return
		}
		feishuIDs = append(feishuIDs, next_users...)
	}

	return
}

func (feishu *feishuService) GetAllUsers() (feishuIDs []*larkcontact.User, err error) {
	root_department := "0"
	departments := []*string{&root_department}
	sub_departments, err := feishu.GetSubDEpartments("")
	departments = append(departments, sub_departments...)
	//log.Print("FeishuDepartments: ", departments)
	if err != nil {
		return
	}
	for _, deparment := range departments {
		users, in_err := feishu.GetDepartmentUser(*deparment, "")
		if in_err != nil {
			err = in_err
			return
		}
		feishuIDs = append(feishuIDs, users...)
	}
	return
}

func (feishu *feishuService) CheckUserAndBind(FeishuUser *larkcontact.User, EntityID uint) (bool, error) {
	user, err := dao.UserDao.GetUserByFeishuID(*FeishuUser.UserId)
	if err != nil {
		return false, err
	}
	if user != nil {
		return false, err
	}

	new_user := model.User{
		UserName: *FeishuUser.Name,
		Password: *FeishuUser.Name,
		EntityID: EntityID,
		FeishuID: *FeishuUser.UserId,
	}
	err = dao.UserDao.Create(new_user)

	return true, err
}

func (feishu *feishuService) FeishuSync(EntityID uint) error {

	FeishuIDs, err := feishu.GetAllUsers()
	if err != nil {
		return err
	}
	//log.Print("FeishuUsers: ", FeishuIDs)
	for _, feishuID := range FeishuIDs {
		_, err = feishu.CheckUserAndBind(feishuID, EntityID)
		if err != nil {
			return err
		}
	}
	return err
}
