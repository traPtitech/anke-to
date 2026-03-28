package model

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
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
			if err := tx.AutoMigrate(&v3TargetUsers{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&v3TargetGroups{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&v3AdministratorUsers{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&v3AdministratorGroups{}); err != nil {
				return err
			}
			if err := tx.Exec("INSERT INTO target_users (questionnaire_id, user_traqid) SELECT questionnaire_id, user_traqid FROM targets").Error; err != nil {
				return err
			}
			if err := tx.Exec("INSERT INTO administrator_users (questionnaire_id, user_traqid) SELECT questionnaire_id, user_traqid FROM administrators").Error; err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("question", "questions"); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&v3Questions{}); err != nil {
				return err
			}
			if err := tx.Migrator().RenameTable("response", "responses"); err != nil {
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
	ID                       int              `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title                    string           `json:"title"           gorm:"type:char(50);size:50;not null"`
	Description              string           `json:"description"     gorm:"type:text;not null"`
	ResTimeLimit             null.Time        `json:"res_time_limit,omitempty"  gorm:"type:TIMESTAMP NULL;default:NULL;"`
	DeletedAt                gorm.DeletedAt   `json:"-"      gorm:"type:TIMESTAMP NULL;default:NULL;"`
	ResSharedTo              string           `json:"res_shared_to"   gorm:"type:char(30);size:30;not null;default:administrators"`
	CreatedAt                time.Time        `json:"created_at"      gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt               time.Time        `json:"modified_at"     gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Administrators           []Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets                  []Targets        `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	TargetGroups             []TargetGroups   `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Questions                []Questions      `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Respondents              []Respondents    `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	IsPublished              bool             `json:"is_published" gorm:"type:boolean;not null;default:true"`
	IsAnonymous              bool             `json:"is_anonymous" gorm:"type:boolean;not null;default:false"`
	IsDuplicateAnswerAllowed bool             `json:"is_duplicate_answer_allowed" gorm:"type:tinyint(4);size:4;not null;default:true"`
}

func (*v3Questionnaires) TableName() string {
	return "questionnaires"
}

type v3TargetUsers struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

func (*v3TargetUsers) TableName() string {
	return "target_users"
}

type v3TargetGroups struct {
	QuestionnaireID int       `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	GroupID         uuid.UUID `gorm:"type:char(36);size:36;not null;primaryKey"`
}

func (*v3TargetGroups) TableName() string {
	return "target_groups"
}

type v3AdministratorUsers struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

func (*v3AdministratorUsers) TableName() string {
	return "administrator_users"
}

type v3AdministratorGroups struct {
	QuestionnaireID int       `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	GroupID         uuid.UUID `gorm:"type:char(36);size:36;not null;primaryKey"`
}

func (*v3AdministratorGroups) TableName() string {
	return "administrator_groups"
}

type v3Questions struct {
	ID              int            `json:"id"                  gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID"     gorm:"type:int(11);not null"`
	PageNum         int            `json:"page_num"            gorm:"type:int(11);not null"`
	QuestionNum     int            `json:"question_num"        gorm:"type:int(11);not null"`
	Type            string         `json:"type"                gorm:"type:char(20);size:20;not null"`
	Body            string         `json:"body"                gorm:"type:text;default:NULL"`
	Description     string         `json:"description"         gorm:"type:text;default:NULL"`
	IsRequired      bool           `json:"is_required"         gorm:"type:tinyint(4);size:4;not null;default:0"`
	DeletedAt       gorm.DeletedAt `json:"-"          gorm:"type:TIMESTAMP NULL;default:NULL"`
	CreatedAt       time.Time      `json:"created_at"          gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Options         []Options      `json:"-"  gorm:"foreignKey:QuestionID"`
	Responses       []Responses    `json:"-"  gorm:"foreignKey:QuestionID"`
	ScaleLabels     []ScaleLabels  `json:"-"  gorm:"foreignKey:QuestionID"`
	Validations     []Validations  `json:"-"  gorm:"foreignKey:QuestionID"`
}

func (*v3Questions) TableName() string {
	return "question"
}
