package main

import (
	"context"
	"fmt"
	"git.tianrang-inc.com/data-brain/trains/config"
	"git.tianrang-inc.com/data-brain/trains/log"
	"git.tianrang-inc.com/data-brain/trains/middleware"
	"git.tianrang-inc.com/data-brain/trains/server"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

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

	var wg sync.WaitGroup

	c := gin.Default()

	//中间件初始化，接口允许AJAX跨域访问
	c.Use(middleware.Cors())

	server.Route(c)

	wg.Add(1)
	go Run(c)

	wg.Wait()

}
