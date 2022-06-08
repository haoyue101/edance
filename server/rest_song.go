package server

import (
	"edance/dao/song"
	"github.com/gin-gonic/gin"
	"net/http"
)

func restListSong(context *gin.Context) {
	listSong, err := song.ListSong()
	if err != nil {
		failed(context, http.StatusNotFound, err)
	}
	success(context, listSong)
}
