package define

import "github.com/dgrijalva/jwt-go"

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
	UserID          uint
	UserName        string
	EntitySuper     bool
	DepartmentSuper bool
	SystemSuper     bool
}

/*
Used for jwt claims
*/
type UserClaims struct {
	UserBasicInfo
	jwt.StandardClaims
}
