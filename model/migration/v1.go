package migration

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func V1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&v1Questionnaires{},
			); err != nil {
				return err
			}
			return nil
		},
	}
}

type v1Questionnaires struct {
	ID             int                `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title          string             `json:"title"           gorm:"type:char(50);size:50;not null"`
	Description    string             `json:"description"     gorm:"type:text;not null"`
	ResTimeLimit   null.Time          `json:"res_time_limit,omitempty"  gorm:"type:TIMESTAMP NULL;default:NULL;"`
	DeletedAt      gorm.DeletedAt     `json:"-"      gorm:"type:TIMESTAMP NULL;default:NULL;"`
	ResSharedTo    string             `json:"res_shared_to"   gorm:"type:char(30);size:30;not null;default:administrators"`
	IsAnonymous    bool               `json:"is_anonymous" gorm:"type:boolean;not null;default:false"`
	CreatedAt      time.Time          `json:"created_at"      gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt     time.Time          `json:"modified_at"     gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Administrators []v1Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets        []v1Targets        `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Questions      []v1Questions      `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Respondents    []v1Respondents    `json:"-"  gorm:"foreignKey:QuestionnaireID"`
}

type v1Administrators struct {
	QuestionnaireID int    `gorm:"type:int(11);not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

type v1Targets struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

type v1Questions struct {
	ID              int             `json:"id"                  gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int             `json:"questionnaireID"     gorm:"type:int(11);not null"`
	PageNum         int             `json:"page_num"            gorm:"type:int(11);not null"`
	QuestionNum     int             `json:"question_num"        gorm:"type:int(11);not null"`
	Type            string          `json:"type"                gorm:"type:char(20);size:20;not null"`
	Body            string          `json:"body"                gorm:"type:text;default:NULL"`
	IsRequired      bool            `json:"is_required"         gorm:"type:tinyint(4);size:4;not null;default:0"`
	DeletedAt       gorm.DeletedAt  `json:"-"          gorm:"type:TIMESTAMP NULL;default:NULL"`
	CreatedAt       time.Time       `json:"created_at"          gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Options         []v1Options     `json:"-"  gorm:"foreignKey:QuestionID"`
	Responses       []v1Responses   `json:"-"  gorm:"foreignKey:QuestionID"`
	ScaleLabels     []v1ScaleLabels `json:"-"  gorm:"foreignKey:QuestionID"`
	Validations     []v1Validations `json:"-"  gorm:"foreignKey:QuestionID"`
}

type v1Respondents struct {
	ResponseID      int            `json:"responseID" gorm:"column:response_id;type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID" gorm:"type:int(11);not null"`
	UserTraqid      string         `json:"user_traq_id,omitempty" gorm:"type:varchar(32);size:32;default:NULL"`
	ModifiedAt      time.Time      `json:"modified_at,omitempty" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty" gorm:"type:TIMESTAMP NULL;default:NULL"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:TIMESTAMP NULL;default:NULL"`
	Responses       []v1Responses  `json:"-"  gorm:"foreignKey:ResponseID;references:ResponseID"`
}

type v1Options struct {
	ID         int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionID int    `gorm:"type:int(11);not null"`
	OptionNum  int    `gorm:"type:int(11);not null"`
	Body       string `gorm:"type:text;default:NULL;"`
}

type v1Responses struct {
	ResponseID int            `json:"-" gorm:"type:int(11);not null"`
	QuestionID int            `json:"-" gorm:"type:int(11);not null"`
	Body       null.String    `json:"response" gorm:"type:text;default:NULL"`
	ModifiedAt time.Time      `json:"-" gorm:"type:timestamp;not null;dafault:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"type:TIMESTAMP NULL;default:NULL"`
}

type v1ScaleLabels struct {
	QuestionID      int    `json:"questionID"        gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	ScaleLabelRight string `json:"scale_label_right" gorm:"type:text;default:NULL;"`
	ScaleLabelLeft  string `json:"scale_label_left"  gorm:"type:text;default:NULL;"`
	ScaleMin        int    `json:"scale_min"         gorm:"type:int(11);default:NULL;"`
	ScaleMax        int    `json:"scale_max"         gorm:"type:int(11);default:NULL;"`
}

type v1Validations struct {
	QuestionID   int    `json:"questionID"    gorm:"type:int(11);not null;primaryKey"`
	RegexPattern string `json:"regex_pattern" gorm:"type:text;default:NULL"`
	MinBound     string `json:"min_bound"     gorm:"type:text;default:NULL"`
	MaxBound     string `json:"max_bound"     gorm:"type:text;default:NULL"`
}
