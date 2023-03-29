package middleware

import (
	"asset-management/app/define"
	"asset-management/myerror"
	"asset-management/utils"
)

func JWTMiddleware() utils.HandlerFunc {
	return func(ctx *utils.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.Unauthorized(myerror.TOKEN_EMPTY, "Cannot find token in request header.")
			ctx.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if utils.IsTokenExpiredError(err) {
			ctx.Unauthorized(myerror.TOKEN_EXPIRED, "Token has expired.")
			ctx.Abort()
			return
		}
		if utils.IsTokenInvalidError(err) {
			ctx.Unauthorized(myerror.TOKEN_INVALID, "Invaild token.")
			ctx.Abort()
			return
		}

		userInfo := define.UserBasicInfo{
			UserID:          claims.UserID,
			UserName:        claims.UserName,
			EntitySuper:     claims.EntitySuper,
			DepartmentSuper: claims.DepartmentSuper,
			SystemSuper:     claims.SystemSuper,
		}

		ctx.Set("user", userInfo)
		ctx.Set("token", token)
		ctx.Next()
	}
}
