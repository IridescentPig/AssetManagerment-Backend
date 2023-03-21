package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseJson struct {
	context *gin.Context
}

type ResponseData struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}

func NewResponseJson(ctx *gin.Context) *responseJson {
	return &responseJson{
		context: ctx,
	}
}

func (r *responseJson) Success(info string, data interface{}) {
	r.context.JSON(http.StatusOK, ResponseData{
		Code: 0,
		Info: info,
		Data: data,
	})
}
