package model

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v3_1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3.1",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE questionnaires MODIFY COLUMN title varchar(1024) NOT NULL").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
