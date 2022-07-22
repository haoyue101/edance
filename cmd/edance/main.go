package main

import (
	"edance/router"
	"edance/util"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

const PortEnv = "EDANCE_PORT"
const DefaultPort = "10088"
const LogFile = "D:/project-go/1/edance/log/edance.log"

func main() {
	eng := gin.Default()
	util.Init(LogFile)
	eng.Use(util.XLog())
	router.InitRouters(eng)
	_ = eng.Run(":" + resolvePort())
}

func resolvePort() string {
	if port := os.Getenv(PortEnv); port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil || 0 < portInt || portInt > 65535 {
			util.Warn("Environment variable " + PortEnv + " is invalid, using default port " + DefaultPort)
			return DefaultPort
		}
		return port
	} else {
		util.Info("Environment variable " + PortEnv + " is undefined, using default port " + DefaultPort)
		return DefaultPort
	}
}
