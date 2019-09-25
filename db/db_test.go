package db

import (
	"testing"
	"time"
	"fmt"
	"git.tianrang-inc.com/data-brain/trains/data"
	"git.tianrang-inc.com/data-brain/trains/log"
)

func TestNewConnDB(t *testing.T) {
	conn := NewConnDB(
		"postgres",
		"host=localhost port=5432 user=postgres password=ps4 dbname=ticket_data sslmode=disable",
		100,
		100,
		30*time.Second,
		20,
	)

	err := conn.Add()
	fmt.Println(err)

	engine, err := conn.Connect()
	err = engine.Ping()

	engine.Sync(
		new(data.User),
		new(data.Admin),
	)
	fmt.Println(err)

	defer conn.Close(engine)
}

func TestConnDB_NewTable(t *testing.T) {
	//fmt.Println(PgConn.PoolSize)
	//count, err := PgConn.BatchAdd(PgConn.PoolSize)
	//fmt.Println(count, err)
	//engines := make([]*xorm.Engine, count)
	for i := 0; i < 50; i++ {
		engine, err := PgConn.Connect()
		if err != nil {
			log.Logger.Warn(err.Error())
			return
		}
		fmt.Println(engine.Ping())
	}
	//defer PgConn.BatchClose(engines)
	//engine := PgConn.Connect()

	//defer PgConn.Close(engine)
}
