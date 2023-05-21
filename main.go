package main

import (
	"asset-management/app/api"
	"asset-management/app/timing"
	"asset-management/middleware"
	"asset-management/routers"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = timezone
	c := timing.Init()
	c.Start()

	api.Initial()
	if os.Getenv("RELEASE") != "" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.Use(middleware.Cors())

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
		r.Run("0.0.0.0:80")
	}
}
