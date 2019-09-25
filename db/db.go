package db

import (
	"git.tianrang-inc.com/data-brain/trains/log"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"time"
	"go.uber.org/zap"
	"errors"
	"fmt"
)

//DB连接客户端
type ClientDB interface {
	//初始化
	Init()

	//新建engine，写入pool池中
	Add() error

	//批量新建engines，写入pool池中，直至池满
	BatchAdd(int) (int, error)

	//从pool池中取出一个engine，建立连接
	Connect() (*xorm.Engine, error)

	//将engine回收到pool池中
	Close(*xorm.Engine)

	//批量回收engines
	BatchClose([]*xorm.Engine)

	//批量自动建表
	NewTable(...interface{}) error
}

type ConnDB struct {
	DriverName     string            `json:"driver_name"`
	DataSourceName string            `json:"data_source_name"`
	MaxOpenConn    int               `json:"max_open_conn"`
	MaxIdleConn    int               `json:"max_idle_conn"`
	MaxLifeTime    time.Duration     `json:"max_life_time"`
	PoolSize       int               `json:"pool_size"`
	Pool           chan *xorm.Engine `json:"pool"`
}

var PgConn *ConnDB

func init() {
	//PgConn = NewConnDB(config.PostgresDriverName, config.PgDataSource, config.MaxOpenConn,
	//	config.MaxIdleConn, config.MaxLifeTime, config.DBPoolSize)
	//PgConn.BatchAdd(PgConn.PoolSize)
}

func NewConnDB(driverName, dataSourceName string, maxOpenConn, maxIdleConn int, maxLifeTime time.Duration, poolSize int) *ConnDB {
	return &ConnDB{
		DriverName:     driverName,
		DataSourceName: dataSourceName,
		MaxOpenConn:    maxOpenConn,
		MaxIdleConn:    maxIdleConn,
		MaxLifeTime:    maxLifeTime,
		PoolSize:       poolSize,
		Pool:           make(chan *xorm.Engine, poolSize),
	}
}

func (conn *ConnDB) Init() {
	conn.BatchAdd(conn.PoolSize)
}

func (conn *ConnDB) Add() (err error) {

	var engine *xorm.Engine
	engine, err = xorm.NewEngine(conn.DriverName, conn.DataSourceName)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	engine.SetMaxOpenConns(conn.MaxOpenConn)
	engine.SetMaxIdleConns(conn.MaxIdleConn)
	engine.SetConnMaxLifetime(conn.MaxLifeTime)

	err = engine.Ping()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	select {
	case conn.Pool <- engine:
	default:
		return errors.New("pool is full")
	}

	return
}

func (conn *ConnDB) BatchAdd(num int) (count int, err error) {
	if num > conn.PoolSize {
		num = conn.PoolSize
	}

	for i := 0; i < num; i++ {
		err = conn.Add()
		if err != nil {
			return
		}
		count += 1
	}

	return
}

func (conn *ConnDB) Connect() (engine *xorm.Engine, err error) {
	select {
	case engine = <-conn.Pool:
		return
	default:
		err = errors.New("pool is empty")
		fmt.Println(engine, err)
		return
	}
}

func (conn *ConnDB) Close(engine *xorm.Engine) {
	select {
	case conn.Pool <- engine:
	default:
		return
	}
}

func (conn *ConnDB) BatchClose(engines []*xorm.Engine) {
	for _, engine := range engines {
		conn.Close(engine)
	}
}

func (conn *ConnDB) NewTable(beans ...interface{}) (err error) {
	engine, err := conn.Connect()
	if err != nil {
		return
	}
	if engine.Ping() != nil {
		log.Logger.Warn(err.Error())
		return
	}
	engine.ShowSQL(true)
	err = engine.Sync2(beans)
	if err != nil {
		log.Logger.Warn("table create failed", zap.String("err", err.Error()))
	}
	defer conn.Close(engine)
	return
}
