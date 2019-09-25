package config

import (
	"os"
	"time"
)

const (
	MaxOpenConn        = 100
	MaxIdleConn        = 100
	MaxLifeTime        = 30 * time.Second
	DBPoolSize         = 20
	PostgresDriverName = "postgres"
	MySQLDriverName    = "mysql"
)

var (
	PgDataSource    string
	MySQLDataSource string
)

func init() {
	PgDataSource = Getenv("PG_DATA_SOURCE", "")
	MySQLDataSource = Getenv("MYSQL_DATA_SOURCE", "")
}

func Getenv(env, def string) string {
	e := os.Getenv(env)
	if e != "" {
		return e
	}

	return def
}
