package define

type UserRegisterReq struct {
	UserName string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}
