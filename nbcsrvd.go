package main

import (
	"github.com/gin-gonic/gin"
	"github.com/9chain/nbcsrvd01/api"
	"github.com/9chain/nbcsrvd01/config"
	"github.com/9chain/nbcsrvd01/web"
)

func main() {
	config.LoadConfig()

	r := gin.Default()
	r.Use(gin.Recovery())

	r.Static("public", "./public")
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "public")
	})
	api.InitApi(r.Group("api"))
	web.InitWeb(r.Group("web"))

	r.Run() // listen and serve on 0.0.0.0:8080
}
