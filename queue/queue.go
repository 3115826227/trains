package queue

import (
	"git.tianrang-inc.com/data-brain/trains/log"
	"time"
)

func StorageData() {
	log.Logger.Info("Data storage...start to push redis")
	time.Sleep(10 * time.Second)
}
