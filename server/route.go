package server

import (
	"git.tianrang-inc.com/data-brain/trains/data/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Test(c *gin.Context) {
	var resp response.Response
	c.JSON(http.StatusOK, resp)
}

func Route(router *gin.Engine) {
	router.GET("/test", Test)
}
