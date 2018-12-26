package model

import (
	"github.com/go-sql-driver/mysql"
)

type ScaleLabels struct {
	ID              int    `json:"questionID" db:"question_id"`
	ScaleLabelRight string `json:"scale_label_right" db:"scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"  db:"scale_label_left"`
	ScaleMin        int    `json:"scale_min" db:"scale_min"`
	ScaleMax        int    `json:"scale_max" db:"scale_max"`
}

type ResponseBody struct {
	QuestionID     int      `json:"questionID"`
	QuestionType   string   `json:"question_type"`
	Response       string   `json:"response"`
	OptionResponse []string `json:"option_response"`
}

type Responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt mysql.NullTime `json:"submitted_at"`
	Body        []ResponseBody `json:"body"`
}
