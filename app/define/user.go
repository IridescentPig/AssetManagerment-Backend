package define

import "github.com/dgrijalva/jwt-go"

/*
.*Req struct are strictly defined according to the api
*/
type UserRegisterReq struct {
	UserName string `form:"userName" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
}

type UserLoginReq struct {
	UserName string `form:"userName" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
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
}

/*
Used for jwt claims
*/
type UserClaims struct {
	UserBasicInfo
	jwt.StandardClaims
}

type UserInfoResponse struct {
	UserBasicInfo
	EntityID     uint
	DepartmentID uint
	SystemID     uint
}
