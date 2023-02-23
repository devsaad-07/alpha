package db

import (
	"fmt"
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

func init() {
	var err error
	singleton.Do(func() {
		dsn := fmt.Sprintf("host=localhost user=postgres password=postgres dbname=cs_india port=5432 sslmode=disable")
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
		db.SetConnMaxLifetime(time.Duration(1 * time.Minute))
	})
	return
}

func GetDb() *gorm.DB {
	return PostgresDB
}

func RunMigration() {
	db := GetDb()
	Migrate(db, DatabaseSchemaMigrations, false)
}
