package db

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB, migrations []*gormigrate.Migration, rollback bool) {

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations)

	if rollback == true {
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("Could not rollback: %v", err)
		}
		log.Printf("Rollback did run successfully")
	} else {
		if err := m.Migrate(); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}
		log.Printf("Migration did run successfully")
	}
}
