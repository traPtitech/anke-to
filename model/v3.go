package model

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func v3() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&v3Targets{}); err != nil {
				return err
			}
      if err := tx.AutoMigrate(&v3Questionnaires{}); err != nil {
				return err
      }
			return nil
		},
	}
}

type v3Targets struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
	IsCanceled      bool   `gorm:"type:tinyint(1);not null;default:0"`
}

func (*v3Targets) TableName() string {
	return "targets"
}

type v3Questionnaires struct {
	ID             int              `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title          string           `json:"title"           gorm:"type:char(50);size:50;not null"`
	Description    string           `json:"description"     gorm:"type:text;not null"`
	ResTimeLimit   null.Time        `json:"res_time_limit,omitempty"  gorm:"type:TIMESTAMP NULL;default:NULL;"`
	DeletedAt      gorm.DeletedAt   `json:"-"      gorm:"type:TIMESTAMP NULL;default:NULL;"`
	ResSharedTo    string           `json:"res_shared_to"   gorm:"type:char(30);size:30;not null;default:administrators"`
	CreatedAt      time.Time        `json:"created_at"      gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt     time.Time        `json:"modified_at"     gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Administrators []Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets        []Targets        `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	TargetGroups   []TargetGroups   `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Questions      []Questions      `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Respondents    []Respondents    `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	IsPublished    bool             `json:"is_published" gorm:"type:boolean;default:false"`
}

func (*v3Questionnaires) TableName() string {
	return "questionnaires"
}