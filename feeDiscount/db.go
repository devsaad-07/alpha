package postgres

import (
	"fmt"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpgx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sync"
	"time"
)

var (
	singleton  sync.Once
	PostgresDB = &gorm.DB{}
)

func InitDB() *gorm.DB {
	var err error
	singleton.Do(func() {
		host := "host=localhost "
		user := "user=anubhav.s "
		password := "password=postgres "
		dbname := "dbname=cs_india "
		port := "port=5433 "
		sslmode := "sslmode=disable"

		dsn := host + user + password + dbname + port + sslmode
		PostgresDB, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "nrpgx",
			DSN:        dsn,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic(fmt.Sprintf("not able to connect to the database. Error:- %s", err.Error()))
		}

		db, err := PostgresDB.DB()
		if err != nil {
			panic(fmt.Sprintf("error occurred while getting db instance object. Error:- %s", err.Error()))
		}

		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(10)
		db.SetConnMaxLifetime(1 * time.Minute)
	})
	return PostgresDB
}
