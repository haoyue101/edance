package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func success(c *gin.Context, data interface{}) {
	if data != nil {
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func failed(c *gin.Context, status int, err error) {
	res := map[string]interface{}{}
	if err != nil {
		res["error"] = err.Error()
		c.AbortWithStatusJSON(status, res)
	}
	c.AbortWithStatus(status)
}
