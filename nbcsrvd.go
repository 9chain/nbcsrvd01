package main

import (
	"github.com/gin-gonic/gin"
	"github.com/9chain/nbcsrvd01/api"
	"github.com/9chain/nbcsrvd01/web"
)

func main() {
	r := gin.Default()

	api.InitApi(r.Group("api"))
	web.InitWeb(r.Group("web"))

	r.Run() // listen and serve on 0.0.0.0:8080
}
