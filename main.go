package main

import (
	"asset-management/app/api"
	"asset-management/routers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	api.Initial()
	r := gin.Default()

	routers.Router.Init(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if len(os.Getenv("DEBUG")) == 0 {
		r.Run("0.0.0.0:8080")
	} else {
		r.Run("0.0.0.0:80")
	}
	// 监听并在 0.0.0.0:8080 上启动服务
}
