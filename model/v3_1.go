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
	IsPublished              bool           `gorm:"type:boolean;not null;default:true"`
	IsAnonymous              bool           `gorm:"type:boolean;not null;default:false"`
	IsDuplicateAnswerAllowed bool           `gorm:"type:tinyint(4);size:4;not null;default:true"`
}

func (*v3_1Questionnaires) TableName() string {
	return "questionnaires"
}

func v3_1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3.1",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&v3_1Questionnaires{})
		},
	}
}
