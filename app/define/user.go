package define

import (
	"asset-management/app/model"

	"github.com/dgrijalva/jwt-go"
)

/*
.*Req struct are strictly defined according to the api
*/
type UserRegisterReq struct {
	UserName string `form:"userName" binding:"required" json:"userName" map:"userName,omitempty"`
	Password string `form:"password" binding:"required" json:"password" map:"password,omitempty"`
}

type UserLoginReq struct {
	UserName string `form:"userName" binding:"required" json:"userName" map:"userName,omitempty"`
	Password string `form:"password" binding:"required" json:"password" map:"password,omitempty"`
}

type UriInfo struct {
	UserName string `uri:"userName"`
}

type ResetReq struct {
	Method   int    `json:"method"`
	Identity int    `json:"identity"`
	Password string `json:"password"`
}

/*
Basic info of user, can be included in other info struct
*/
type UserBasicInfo struct {
	UserID          uint   `json:"user_id"`
	UserName        string `json:"username"`
	EntitySuper     bool   `json:"entity_super"`
	DepartmentSuper bool   `json:"department_super"`
	SystemSuper     bool   `json:"system_super"`
	EntityID        uint   `json:"entity_id"`
	DepartmentID    uint   `json:"department_id"`
}

/*
Used for jwt claims
*/
type UserClaims struct {
	UserBasicInfo
	jwt.StandardClaims
}

type UserInfo struct {
	UserID          uint              `json:"user_id" copier:"ID"`
	UserName        string            `json:"username" copier:"UserName"`
	Ban             bool              `json:"lock"`
	IsEmployee      bool              `json:"id0" default:"true"`
	DepartmentSuper bool              `json:"id1"`
	EntitySuper     bool              `json:"id2"`
	SystemSuper     bool              `json:"id3"`
	EntityID        uint              `json:"entity_id"`
	Entity          *model.Entity     `json:"entity"`
	DepartmentID    uint              `json:"department_id"`
	Department      *model.Department `json:"department"`
}

type UserInfoResponse struct {
	User UserInfo `json:"user"`
}

type UserListResponse struct {
	UserList []UserInfo `json:"user_list"`
}
