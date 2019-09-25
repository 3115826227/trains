package main

import (
	"context"
	"fmt"
	"git.tianrang-inc.com/data-brain/trains/config"
	"git.tianrang-inc.com/data-brain/trains/log"
	"git.tianrang-inc.com/data-brain/trains/middleware"
	"git.tianrang-inc.com/data-brain/trains/queue"
	"git.tianrang-inc.com/data-brain/trains/server"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func init() {

	var wg sync.WaitGroup
	router := gin.New()

	router.Use(middleware.Cors())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"ret":     404,
			"message": "找不到该路由",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"ret":     404,
			"message": "找不到该方法",
		})
	})

	server.Route(router)

	wg.Add(1)
	go Run(router)

	wg.Wait()
}

func Run(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + config.Getenv("SERVER_PORT", "8057"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Error(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//服务关闭前将内存中的数据进行redis写入
	queue.StorageData()

	log.Logger.Info("Shutdown Server....")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Error("Server Shutdown:" + err.Error())
	}
	log.Logger.Info("Server exiting")

	pid := fmt.Sprintf("%d", os.Getpid())
	_, openErr := os.OpenFile("pid", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr == nil {
		_ = ioutil.WriteFile("pid", []byte(pid), 0)
	}
}

func main() {
}
