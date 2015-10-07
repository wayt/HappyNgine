package sql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"time"
)

var DB *gorm.DB

func init() {

	username := env.Get("HAPPY_SQL_USERNAME")
	password := env.Get("HAPPY_SQL_PASSWORD")
	host := env.Get("HAPPY_SQL_HOST")
	port := env.Get("HAPPY_SQL_PORT")
	dbName := env.Get("HAPPY_SQL_DB_NAME")
	pingInterval := env.GetInt("HAPPY_SQL_PING_INTERVAL")
	maxIdleConns := env.GetInt("HAPPY_SQL_MAX_IDLE_CONNS")
	maxOpenConns := env.GetInt("HAPPY_SQL_MAX_OPEN_CONNS")

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", username, password, host, port, dbName))
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)

	db.LogMode(env.Get("DEBUG") == "1")

	DB = &db

	go func() {
		for true {
			time.Sleep(time.Duration(pingInterval) * time.Second)

			log.Debugln("Ping sql to keep-alive...")
			db.DB().Ping()
		}
	}()
}
