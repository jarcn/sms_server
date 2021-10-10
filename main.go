package main

import (
	"errors"
	"fmt"
	"os"
	"sms_server/conf"
	"sms_server/ctrl"
	"sms_server/r2m"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	c := conf.Config{}
	conf.Set(c)
	switch {
	case c.Env == "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func init() {
	if appkey := os.Getenv("yun.pian.appkey"); appkey != "" {
		conf.Cfg.YPAppKey = appkey
	}
	if conf.Cfg.YPAppKey == "" {
		conf.Cfg.YPAppKey = "1234567890"
		fmt.Println("云片appkey不存在")
	}
	if timeOut := os.Getenv("check.code.time.out"); timeOut != "" {
		if val, err := strconv.Atoi(strings.TrimSpace(timeOut)); err != nil {
			panic(errors.New("[check.code.time.out] 格式错误"))
		} else {
			conf.Cfg.KeyTimeOut = int64(val)
		}
	}
}

func initRouters(g *gin.Engine) {
	g.POST("/code/send", ctrl.SendCode)
	g.POST("/code/check", ctrl.CheckCode)
	g.POST("/async", ctrl.AsynchronousPost)
}

func main() {
	fmt.Println("start sms_server ...")
	g := gin.Default()
	initRouters(g)
	r2m.Job{}.Start()
	g.Run()
}
