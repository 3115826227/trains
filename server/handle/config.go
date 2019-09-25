package handle

import (
	"github.com/gin-gonic/gin"
	"git.tianrang-inc.com/data-brain/trains/config"
	"git.tianrang-inc.com/data-brain/trains/db"
	"fmt"
	"net/http"
	"git.tianrang-inc.com/data-brain/trains/data/response"
)

func PostgresConnect(c *gin.Context) {
	host := c.PostForm("host")
	port := c.PostForm("port")
	user := c.PostForm("user")
	password := c.PostForm("password")
	dbName := c.PostForm("db")
	config.PgDataSource = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		host, port, user, password, dbName)
	conn := db.NewConnDB(config.PostgresDriverName, config.PgDataSource, config.MaxOpenConn,
		config.MaxIdleConn, config.MaxLifeTime, config.DBPoolSize)
	err := conn.Add()
	if err != nil {
		c.JSON(http.StatusOK, response.Response{
			Message: err.Error(),
		})
		return
	}

	engine, err := conn.Connect()
	if err != nil {
		c.JSON(http.StatusOK, response.Response{
			Message: err.Error(),
		})
		return
	}

	if engine.Ping() != nil {
		c.JSON(http.StatusOK, response.Response{
			Message: err.Error(),
		})
		return
	}

	db.PgConn = conn
	db.PgConn.Init()

	c.JSON(http.StatusOK, response.Response{
		Data: string("连接成功"),
	})
}
