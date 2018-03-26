package main

import (
	"fmt"
	"github.com/9chain/nbcsrvd01/api"
	"github.com/9chain/nbcsrvd01/config"
	"github.com/9chain/nbcsrvd01/state"
	"github.com/9chain/nbcsrvd01/web"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"os"
)

func initSeelog() {
	cfgPath := config.Cfg.App.SeeLogXml
	if _, err := os.Stat(cfgPath); err == nil {
		logger, err := log.LoggerFromConfigAsFile(cfgPath)
		if err != nil {
			panic(err)
		}

		log.ReplaceLogger(logger)
		return
	}

	fmt.Println("use default seelog config")

	defaultConfig := `
<seelog>
    <outputs formatid="main">
        <console />
    </outputs>
    <formats>
        <format id="main" format="%l %Date %Time %File:%Line %Msg%n"/>
    </formats>
</seelog>`
	logger, err := log.LoggerFromConfigAsString(defaultConfig)
	if err != nil {
		panic(err)
	}

	log.ReplaceLogger(logger)
}

func main() {
	config.LoadConfig()
	initSeelog()

	state.Init()

	r := gin.Default()
	r.Use(gin.Recovery())

	r.Static("public", "./public")
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "public")
	})

	web.InitWeb(r.Group("web"))
	api.InitApi(r.Group("api"))

	r.Run() // listen and serve on 0.0.0.0:8080
}
