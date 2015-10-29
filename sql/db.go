package sql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"strings"
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

// Run severale queries within a transactional context
// The transaction is automaticly retried if needed
func Transaction(queries ...func(tx *gorm.DB) error) error {

	// Should never happens
	if len(queries) == 0 {
		return errors.New("No Transaction to be run...")
	}

	// Last error, used to return what happened if needed
	var lastErr error = nil

	// Retry transaction up to 3 times
	for i := 0; i < 3; i++ {

		if i != 0 { // Sleep for 50ms if it's not the first try
			time.Sleep(50 * time.Millisecond)
		}

		// BEGIN
		tx := DB.Begin()
		rollback := true // In case of rollback
		defer func() {
			if rollback {
				tx.Rollback()
			}
		}()

		// Run each queries with the tx context
		for _, q := range queries {

			if lastErr = q(tx); lastErr != nil {
				if strings.Contains(lastErr.Error(), "try restarting transaction") { // Generic error handling, can occure for multiple errors
					continue
				} else {
					break // Another error, break and return via lastErr
				}
			}
		}

		// Try to commit the transaction, no retry here
		if err := tx.Commit().Error; err != nil {
			return err
		}

		// Tell to not rollback this one
		rollback = false
		break // we're done
	}

	return lastErr
}
