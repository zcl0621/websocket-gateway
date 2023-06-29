package main

import (
	"os"
	"ws-gateway/config"
	"ws-gateway/logger"
	"ws-gateway/redisUtils"
	"ws-gateway/server"
)

func main() {
	config.InitConf()
	redisUtils.InitRedis()
	if err := server.StartGinServer(); err != nil {
		logger.Logger("main 启动http服务失败", "error", err, "")
		os.Exit(2)
		return
	}
}
