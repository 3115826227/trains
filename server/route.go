package server

import (
	"git.tianrang-inc.com/data-brain/trains/data/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"git.tianrang-inc.com/data-brain/trains/server/handle"
)

func Test(c *gin.Context) {
	var resp response.Response
	c.JSON(http.StatusOK, resp)
}

func Route(router *gin.Engine) {

	cfg := router.Group("/config", )
	{
		cfg.POST("/postgres", handle.PostgresConnect)
	}

	router.GET("/test", Test)
}
