package utils

import (
	"asset-management/myerror"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func GetNewContext() Context {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	return Context{Context: ctx}
}

func TestResponse(t *testing.T) {
	{
		context := GetNewContext()
		context.BadRequest(1, "Invalid request body")
		assert.Equal(t, http.StatusBadRequest, context.Writer.Status(), "response error")
	}
	{
		context := GetNewContext()
		context.Unauthorized(myerror.TOKEN_INVALID, "Token invalid")
		assert.Equal(t, http.StatusUnauthorized, context.Writer.Status(), "response error")
	}
	{
		context := GetNewContext()
		context.Forbidden(1, "Permission Denied")
		assert.Equal(t, http.StatusForbidden, context.Writer.Status(), "response error")
	}
	{
		context := GetNewContext()
		context.NotFound(1, "Route not found")
		assert.Equal(t, http.StatusNotFound, context.Writer.Status(), "response error")
	}
	{
		context := GetNewContext()
		context.Success(nil)
		assert.Equal(t, http.StatusOK, context.Writer.Status(), "response error")
	}
	{
		context := GetNewContext()
		context.InternalError("Database operation failed")
		assert.Equal(t, http.StatusInternalServerError, context.Writer.Status(), "response error")
	}

	helloFunc := func(ctx *Context) {
		log.Println(ctx.Writer.Status())
	}

	handleFunc := Handler(helloFunc)
	context := GetNewContext()
	context.Success(nil)
	handleFunc(context.Context)
}
