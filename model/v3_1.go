package model

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type v3_1Questionnaires struct {
	ID                       int            `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title                    string         `gorm:"type:varchar(1024);size:1024;not null"`
	Description              string         `gorm:"type:text;not null"`
	ResTimeLimit             null.Time      `gorm:"type:TIMESTAMP NULL;default:NULL;"`
	DeletedAt                gorm.DeletedAt `gorm:"type:TIMESTAMP NULL;default:NULL;"`
	ResSharedTo              string         `gorm:"type:char(30);size:30;not null;default:administrators"`
	CreatedAt                time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt               time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	IsPublished              bool           `gorm:"type:boolean;not null;default:false"`
	IsAnonymous              bool           `gorm:"type:boolean;not null;default:false"`
	IsDuplicateAnswerAllowed bool           `gorm:"type:boolean;not null;default:false"`
}

func (*v3_1Questionnaires) TableName() string {
	return "questionnaires"
}

func v3_1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3.1",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AlterColumn(&v3_1Questionnaires{}, "Title"); err != nil {
				return err
			}
			if err := tx.Migrator().AlterColumn(&v3_1Questionnaires{}, "IsPublished"); err != nil {
				return err
			}
			return tx.Migrator().AlterColumn(&v3_1Questionnaires{}, "IsDuplicateAnswerAllowed")
		},
	}
}
