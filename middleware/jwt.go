package middleware

import (
	"asset-management/app/define"
	"asset-management/myerror"
	"asset-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			utils.NewResponseJson(context).Error(http.StatusUnauthorized, myerror.TOKEN_EMPTY, "Cannot find token in request header.", nil)
			return
		}

		claims, err := utils.ParseToken(token)
		if utils.IsTokenExpiredError(err) {
			utils.NewResponseJson(context).Error(http.StatusUnauthorized, myerror.TOKEN_EXPIRED, "Token has expired.", nil)
			return
		}
		if utils.IsTokenInvalidError(err) {
			utils.NewResponseJson(context).Error(http.StatusUnauthorized, myerror.TOKEN_INVALID, "Invaild token.", nil)
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
