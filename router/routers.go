package router

import (
	"github.com/gin-gonic/gin"
)

var (
	song := SongApi{
		prefix: "/song"
	}
)

func InitRouters(eng *gin.Engine) {
	eng.GET(song.Prefix("/list"), song.ListSongs)
}
