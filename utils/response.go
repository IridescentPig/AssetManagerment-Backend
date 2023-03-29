package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseData struct {
	Error ErrorData   `json:"error"`
	Data  interface{} `json:"data"`
}

var NoError = ErrorData{0, ""}

func (ctx *Context) BadRequest(code int, err string) {
	ctx.JSON(http.StatusBadRequest, ResponseData{
		ErrorData{code, err},
		nil,
	})
}

func (ctx *Context) Unauthorized(code int, err string) {
	ctx.JSON(http.StatusUnauthorized, ResponseData{
		ErrorData{code, err},
		nil,
	})
}

func (ctx *Context) Forbidden(code int, err string) {
	ctx.JSON(http.StatusForbidden, ResponseData{
		ErrorData{code, err},
		nil,
	})
}

func (ctx *Context) NotFound(code int, err string) {
	ctx.JSON(http.StatusNotFound, ResponseData{
		ErrorData{code, err},
		nil,
	})
}

func (ctx *Context) Success(data interface{}) {
	ctx.JSON(http.StatusOK, ResponseData{
		NoError,
		data,
	})
}

func (ctx *Context) InternalError(err string) {
	ctx.JSON(http.StatusInternalServerError, ResponseData{
		ErrorData{-1, err},
		nil,
	})
}

type HandlerFunc func(ctx *Context)

func Handler(f HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f(&Context{ctx})
	}
}
