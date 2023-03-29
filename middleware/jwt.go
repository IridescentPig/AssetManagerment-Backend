package middleware

import (
	"asset-management/app/define"
	"asset-management/myerror"
	"asset-management/utils"
)

func JWTMiddleware() utils.HandlerFunc {
	return func(context *utils.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			context.Unauthorized(myerror.TOKEN_EMPTY, "Cannot find token in request header.")
			return
		}

		claims, err := utils.ParseToken(token)
		if utils.IsTokenExpiredError(err) {
			context.Unauthorized(myerror.TOKEN_EXPIRED, "Token has expired.")
			return
		}
		if utils.IsTokenInvalidError(err) {
			context.Unauthorized(myerror.TOKEN_INVALID, "Invaild token.")
			return
		}

		userInfo := define.UserBasicInfo{
			UserID:          claims.UserID,
			UserName:        claims.UserName,
			EntitySuper:     claims.EntitySuper,
			DepartmentSuper: claims.DepartmentSuper,
			SystemSuper:     claims.SystemSuper,
		}

		context.Set("user", userInfo)
		context.Set("token", token)
		context.Next()
	}
}
