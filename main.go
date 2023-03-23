package main

import (
	"asset-management/app/api"
	"asset-management/routers"

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
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
