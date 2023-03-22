package define

type UserRegisterReq struct {
	UserName string `form:"userName" json:"username"`
	Password string `form:"passWord" json:"password"`
}

type UserLoginReq struct {
	UserName string `form:"userName" json:"username"`
	Password string `form:"passWord" json:"password"`
}
