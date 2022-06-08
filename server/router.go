package server

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(eng *gin.Engine) {
	eng.GET("/song/list", restListSong)
}
