package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/abelianl2/bridge-server/config"
	"github.com/abelianl2/bridge-server/server"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)
	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildLevel(xlog.Level(cfg.LogLevel)).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/bridge", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := server.NewService(cfg, xLog)

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	root.POST("/bridge/submit", srv.SaveTx)
	root.POST("/bridge/submitWithMemo", srv.SaveTxAndMemo)
	root.GET("/bridge/hash", srv.GetToAddress)
	root.POST("/bridge/notify/:id", srv.NotifyTx)

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
