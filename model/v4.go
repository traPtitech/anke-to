package model

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v4() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "4",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&ReminderTargets{})
		},
	}
}
