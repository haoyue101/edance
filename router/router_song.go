package router

import (
	"edance/dao/song"
	"edance/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SongApi struct {
	prefix string
}

func (api *SongApi) Url(url string) string {
	return api.prefix + url
}

func (api *SongApi) ListSongs(context *gin.Context) {
	listSong, err := song.ListSong()
	if err != nil {
		server.failed(context, http.StatusNotFound, err)
	}
	server.success(context, listSong)
}
