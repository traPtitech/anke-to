package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

//Respondents respondentsテーブルの構造体
type Respondents struct {
	ResponseID int `gorm:"type:int(11);NOT NULL;PRIMARY_KEY;AUTO_INCREMENT"`
	QuestionnaireID int `gorm:"type:int(11);NOT NULL;"`
	UserTraqid null.String `gorm:"type:char(30);"`
	ModifiedAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT CURRENT_TIMESTAMP;"`
	SubmitedAt null.Time `gorm:"type:timestamp;"`
	DeletedAT null.Time `gorm:"type:timestamp;"`
}