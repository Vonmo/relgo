package models

import (
	"time"

	"github.com/elzor/relgo/log"
	"github.com/elzor/relgo/metrics"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// -----------------------------------------------------------------------------

var db *sqlx.DB

func InitDB(dataSourceName string, maxCons int, maxIdleConns int, metrics *metrics.Metrics) {
	var err error

	if dataSourceName == "" {
		log.Error("db is not configured")
		return
	}

	db, err = sqlx.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	log.Debug("db connected")

	db.SetConnMaxLifetime(time.Second * 30)
	db.SetMaxOpenConns(maxCons)
	db.SetMaxIdleConns(maxIdleConns)
	log.Debug("db pool: ok [m ", maxCons, ", i ", maxIdleConns, "]")
	go func() {
		for {
			if err = db.Ping(); err != nil {
				log.Error(err)
				time.Sleep(5 * time.Second)
				InitDB(dataSourceName, maxCons, maxIdleConns, metrics)
				return
			}
			currentOpenConnections := db.Stats().OpenConnections
			metrics.Set("db.open_connections", int64(currentOpenConnections))
			time.Sleep(time.Second)
		}
	}()
}
