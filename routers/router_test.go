package routers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouter(t *testing.T) {
	r := gin.Default()
	_ = Router.Init(r)
}
