package api

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
)

type userApi struct {
}

var UserApi *userApi

func newUserApi() *userApi {
	return &userApi{}
}

func init() {
	UserApi = newUserApi()
}

func (user *userApi) CheckIdentity(ctx *utils.Context, entityID uint) bool {
	systemSuper := service.UserService.SystemSuper(ctx)
	if systemSuper {
		return true
	}
	entitySuper := service.UserService.EntitySuper(ctx)
	if !entitySuper {
		return false
	}
	isInEntity := service.EntityService.CheckIsInEntity(ctx, entityID)
	return isInEntity
}

func (user *userApi) GetOperatorID(ctx *utils.Context) uint {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			return userInfo.UserID
		}
	}
	return 0
}

func (user *userApi) GetOperatorInfo(ctx *utils.Context) *define.UserBasicInfo {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			return &userInfo
		}
	}
	return nil
}

func (user *userApi) UserRegister(ctx *utils.Context) {
	var req define.UserRegisterReq

	if err := ctx.MustBindWith(&req, binding.JSON); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	exists, err := service.UserService.ExistsUser(req.UserName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if exists {
		ctx.BadRequest(1, "Duplicated Name")
		return
	}
	err = service.UserService.CreateUser(req.UserName, req.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

func (user *userApi) UserLogin(ctx *utils.Context) {
	// 是否需要加密传输？
	// 是否需要通过中间件处理？
	var req define.UserLoginReq
	var thisUser *model.User

	if err := ctx.MustBindWith(&req, binding.JSON); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	var err error
	var token string
	token, thisUser, err = service.UserService.VerifyPasswordAndGetUser(req.UserName, req.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(1, "Wrong Username Or Password")
		return
	} else if thisUser.Ban {
		ctx.BadRequest(3, "User Banned")
		return
	}

	var userInfo define.UserInfo
	err = copier.Copy(&userInfo, thisUser)
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

func (user *userApi) UserLogout(ctx *utils.Context) {
	// 使用中间件验证 token 是否正确
	ctx.Success(nil)
}

func (user *userApi) UserCreate(ctx *utils.Context) {
	// TODO: 暂时使用 register 的 req
	var req define.UserRegisterReq
	if err := ctx.MustBindWith(&req, binding.JSON); err != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	// 暂时按照只有超级用户可以创建来处理
	if !service.UserService.SystemSuper(ctx) {
		ctx.Forbidden(2, "Permission Denied.")
		return
	}
	exists, err := service.UserService.ExistsUser(req.UserName)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if exists {
		ctx.BadRequest(1, "Duplicated Name")
		return
	}
	err = service.UserService.CreateUser(req.UserName, req.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	ctx.Success(nil)
}

func (user *userApi) ResetContent(ctx *utils.Context) {
	var req define.ResetReq
	username := ctx.Param("username")
	errReq := ctx.MustBindWith(&req, binding.JSON)
	if errReq != nil {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	//查找用户是否存在
	exists, err := service.UserService.ExistsUser(username)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if !exists {
		ctx.BadRequest(1, "User Not Found")
		return
	}
	if req.Method == 0 {
		// 修改身份
		// 暂时按照只有超级用户可以创建来处理
		if !service.UserService.SystemSuper(ctx) {
			ctx.Forbidden(2, "Permission Denied.")
			return
		} else {
			err = service.UserService.ModifyUserIdentity(username, req.Identity)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			}
		}
	} else if req.Method == 1 {
		// 修改密码
		// 超级用户和自己都应该可以修改密码
		// 自己修改密码需要验证是否为本人
		if !service.UserService.SystemSuper(ctx) {
			get_username, err := service.UserService.UserName(ctx)
			if err != nil {
				ctx.InternalError(err.Error())
				return
			}
			if username != get_username {
				ctx.Forbidden(2, "Permission Denied.")
			}
		}
		err = service.UserService.ModifyUserPassword(username, req.Password)
		if err != nil {
			ctx.InternalError(err.Error())
			return
		}
	} else {
		ctx.BadRequest(-1, "Invalid request body.")
		return
	}
	ctx.Success(nil)
}

func (user *userApi) LockUser(ctx *utils.Context) {
	username := ctx.Param("username")

	thisUser, err := service.UserService.GetUserByName(username)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}
	hasIdentity := user.CheckIdentity(ctx, thisUser.EntityID)
	if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	err = service.UserService.ModifyUserBanstate(username, true)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

func (user *userApi) UnlockUser(ctx *utils.Context) {
	username := ctx.Param("username")

	thisUser, err := service.UserService.GetUserByName(username)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}
	hasIdentity := user.CheckIdentity(ctx, thisUser.EntityID)
	if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	err = service.UserService.ModifyUserBanstate(username, false)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for DELETE /user/{user_id}
*/
func (user *userApi) DeleteUser(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}

	hasIdentity := user.CheckIdentity(ctx, thisUser.EntityID)
	if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}
	if user.GetOperatorID(ctx) == userID {
		ctx.BadRequest(myerror.DELETE_USER_SELF, myerror.DELETE_USER_SELF_INFO)
		return
	}

	err = service.UserService.DeleteUser(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
修改密码权限：超级管理员或实体系统管理员或自己
*/
func (user *userApi) CheckChangePasswdIdentity(ctx *utils.Context) (*model.User, bool) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return nil, false
	}

	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return nil, false
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return nil, false
	}

	isSelf := user.GetOperatorID(ctx) == userID
	hasIdentity := user.CheckIdentity(ctx, thisUser.EntityID)
	if !hasIdentity && !isSelf {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return nil, false
	}

	return thisUser, true
}

/*
Handle func for GET /user/info/{user_id}
*/
func (user *userApi) GetUserInfoByID(ctx *utils.Context) {
	// userID, err := service.EntityService.GetParamID(ctx, "user_id")
	// if err != nil {
	// 	return
	// }

	// thisUser, err := service.UserService.GetUserByID(userID)
	// if err != nil {
	// 	ctx.InternalError(err.Error())
	// 	return
	// } else if thisUser == nil {
	// 	ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
	// 	return
	// }

	// isSelf := user.GetOperatorID(ctx) == userID
	// hasIdentity := user.CheckIdentity(ctx, thisUser.EntityID)
	// if !hasIdentity && !isSelf {
	// 	ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
	// 	return
	// }
	thisUser, isOK := user.CheckChangePasswdIdentity(ctx)
	if !isOK {
		return
	}

	var userInfoRes define.UserInfo
	err := copier.Copy(&userInfoRes, thisUser)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	userInfoResponse := define.UserInfoResponse{
		User: userInfoRes,
	}

	ctx.Success(userInfoResponse)
}

/*
Handle func for GET /user/list
*/
func (user *userApi) GetAllUsers(ctx *utils.Context) {
	systemSuper := service.UserService.SystemSuper(ctx)
	if !systemSuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	userList, err := service.UserService.GetAllUsers()
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	var userListRes []define.UserInfo
	err = copier.Copy(&userListRes, userList)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	userListResponse := define.UserListResponse{
		UserList: userListRes,
	}

	ctx.Success(userListResponse)
}

/*
Handle func for POST /user/info/:user_id/password
*/
func (user *userApi) ChangePassword(ctx *utils.Context) {
	thisUser, isOK := user.CheckChangePasswdIdentity(ctx)
	if !isOK {
		return
	}

	var changePasswordReq define.ChangePasswordReq
	err := ctx.MustBindWith(&changePasswordReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	err = service.UserService.ModifyUserPassword(thisUser.UserName, changePasswordReq.Password)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /user/info/:user_id/entity
*/
func (user *userApi) ChangeUserEntity(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}

	var changeUserEntityReq define.ChangeUserEntityReq
	err = ctx.MustBindWith(&changeUserEntityReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
	}

	exists, err := service.EntityService.ExistsEntityByID(changeUserEntityReq.EntityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}
	if !exists {
		ctx.NotFound(myerror.ENTITY_NOT_FOUND, myerror.ENTITY_NOT_FOUND_INFO)
		return
	}

	err = service.UserService.ModifyUserEntity(userID, changeUserEntityReq.EntityID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for POST /user/info/{user_id}/department
*/
func (user *userApi) ChangeUserDepartment(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}

	entityID, err := service.EntityService.GetParamID(ctx, "entity_id")
	if err != nil {
		return
	}
	hasIdentity := DepartmentApi.CheckDepartmentModifyIdentity(ctx, entityID)
	if !hasIdentity {
		return
	}

	var changeUserDepartmentReq define.ChangeUserDepartmentReq
	err = ctx.MustBindWith(&changeUserDepartmentReq, binding.JSON)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_BODY, myerror.INVALID_BODY_INFO)
		return
	}

	isValid := DepartmentApi.CheckEntityDepartmentValid(ctx, thisUser.EntityID, changeUserDepartmentReq.DepartmentID)
	if !isValid {
		return
	}

	err = service.UserService.ModifyUserDepartment(userID, changeUserDepartmentReq.DepartmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(nil)
}

/*
Handle func for GET /users/{userId}/assets
*/
func (user *userApi) GetAssetsByUser(ctx *utils.Context) {
	userID, err := service.EntityService.GetParamID(ctx, "user_id")
	if err != nil {
		return
	}

	thisUser, err := service.UserService.GetUserByID(userID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	} else if thisUser == nil {
		ctx.BadRequest(myerror.USER_NOT_FOUND, myerror.USER_NOT_FOUND_INFO)
		return
	}

	assets, err := service.AssetService.GetAssetByUser(thisUser.ID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	ctx.Success(assets)

}
