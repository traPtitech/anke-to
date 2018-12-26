package main

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

type questionnaires struct {
	ID           int            `json:"questionnaireID" db:"id"`
	Title        string         `json:"title"           db:"title"`
	Description  string         `json:"description"     db:"description"`
	ResTimeLimit mysql.NullTime `json:"res_time_limit"  db:"res_time_limit"`
	DeletedAt    mysql.NullTime `json:"deleted_at"      db:"deleted_at"`
	ResSharedTo  string         `json:"res_shared_to"   db:"res_shared_to"`
	CreatedAt    time.Time      `json:"created_at"      db:"created_at"`
	ModifiedAt   time.Time      `json:"modified_at"     db:"modified_at"`
}

type questions struct {
	ID              int            `json:"id"                  db:"id"`
	QuestionnaireId int            `json:"questionnaireID"     db:"questionnaire_id"`
	PageNum         int            `json:"page_num"            db:"page_num"`
	QuestionNum     int            `json:"question_num"        db:"question_num"`
	Type            string         `json:"type"                db:"type"`
	Body            string         `json:"body"                db:"body"`
	IsRequrired     bool           `json:"is_required"         db:"is_required"`
	DeletedAt       mysql.NullTime `json:"deleted_at"          db:"deleted_at"`
	CreatedAt       time.Time      `json:"created_at"          db:"created_at"`
}

type scaleLabels struct {
	ID              int    `json:"questionID" db:"question_id"`
	ScaleLabelRight string `json:"scale_label_right" db:"scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"  db:"scale_label_left"`
	ScaleMin        int    `json:"scale_min" db:"scale_min"`
	ScaleMax        int    `json:"scale_max" db:"scale_max"`
}

type responseBody struct {
	QuestionID     int      `json:"questionID"`
	QuestionType   string   `json:"question_type"`
	Response       string   `json:"response"`
	OptionResponse []string `json:"option_response"`
}

type responses struct {
	ID          int            `json:"questionnaireID"`
	SubmittedAt mysql.NullTime `json:"submitted_at"`
	Body        []responseBody `json:"body"`
}

type questionIDType struct {
	ID   int    `db:"id"`
	Type string `db:"type"`
}
