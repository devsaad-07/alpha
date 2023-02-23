package db

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var DatabaseSchemaMigrations = []*gormigrate.Migration{
	{
		ID: "1",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Rules{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&Rules{})
		},
	},
	{
		ID: "2",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserMetrics{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&UserMetrics{})
		},
	},
}
