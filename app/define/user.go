package define

type UserRegisterReq struct {
	UserName string `form:"userName" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
}

type UserLoginReq struct {
	UserName string `form:"userName" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
}
