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
	r.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "http://localhost:8080, AssetManagement-Frontend-dev-BinaryAbstract.app.secoder.net, ")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Next()
	})

	routers.Router.Init(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if gin.Mode() == gin.DebugMode {
		if os.Getenv("DEBUG") == "" {
			r.Run("0.0.0.0:8080")
		} else {
			r.Run("0.0.0.0:80")
		}
	} else {
		r.Run("0.0.0.0:8080")
	}
}
