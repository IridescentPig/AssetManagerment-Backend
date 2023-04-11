package middleware

import (
	"asset-management/app/define"
	"asset-management/myerror"
	"asset-management/utils"
)

func CheckSystemSuper() utils.HandlerFunc {
	return func(ctx *utils.Context) {
		userInfo, exists := ctx.Get("user")
		if exists {
			if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
				if userInfo.SystemSuper {
					ctx.Next()
				} else {
					ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
					ctx.Abort()
					return
				}
			}
		} else {
			ctx.Unauthorized(myerror.TOKEN_INVALID, myerror.TOKEN_INVALID_INFO)
			ctx.Abort()
			return
		}
	}
}

func CheckEntitySuper() utils.HandlerFunc {
	return func(ctx *utils.Context) {
		userInfo, exists := ctx.Get("user")
		if exists {
			if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
				if userInfo.EntitySuper {
					ctx.Next()
				} else {
					ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
					ctx.Abort()
					return
				}
			}
		} else {
			ctx.Unauthorized(myerror.TOKEN_INVALID, myerror.TOKEN_INVALID_INFO)
			ctx.Abort()
			return
		}
	}
}

func CheckDepartmentSuper() utils.HandlerFunc {
	return func(ctx *utils.Context) {
		userInfo, exists := ctx.Get("user")
		if exists {
			if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
				if userInfo.DepartmentSuper {
					ctx.Next()
				} else {
					ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
					ctx.Abort()
					return
				}
			}
		} else {
			ctx.Unauthorized(myerror.TOKEN_INVALID, myerror.TOKEN_INVALID_INFO)
			ctx.Abort()
			return
		}
	}
}
