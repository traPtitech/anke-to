package model

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type v3_2Questionnaires struct {
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
	RandomOrderSalt          string         `gorm:"type:char(64);size:64;not null;default:''"`
}

func (*v3_2Questionnaires) TableName() string {
	return "questionnaires"
}

func v3_2() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3.2",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AddColumn(&v3_2Questionnaires{}, "RandomOrderSalt"); err != nil {
				return err
			}

			var questionnaireIDs []int
			if err := tx.Table("questionnaires").Pluck("id", &questionnaireIDs).Error; err != nil {
				return err
			}
			for _, questionnaireID := range questionnaireIDs {
				salt, err := generateRandomOrderSalt()
				if err != nil {
					return err
				}
				if err := tx.Table("questionnaires").
					Where("id = ?", questionnaireID).
					Update("random_order_salt", salt).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
