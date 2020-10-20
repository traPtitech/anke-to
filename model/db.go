package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	db        *gorm.DB
	allTables = []interface{}{
		Questionnaires{},
		Question{},
		Respondents{},
		Response{},
		Administrators{},
		Options{},
		ScaleLabels{},
		Targets{},
		Validations{},
	}
)

// EstablishConnection DBと接続
func EstablishConnection() (*gorm.DB, error) {
	user := os.Getenv("MARIADB_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("MARIADB_PASSWORD")
	if pass == "" {
		pass = "password"
	}

	host := os.Getenv("MARIADB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("MARIADB_DATABASE")
	if dbname == "" {
		dbname = "anke-to"
	}

	_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4")
	db = _db
	db = db.BlockGlobalUpdate(true)
	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	return db, err
}

// Migrate DBのMigrationを行う
func Migrate() error {
	err := db.AutoMigrate(allTables...).Error
	if err != nil {
		return fmt.Errorf("failed in table's migration: %w", err)
	}

	err = db.
		Model(&Options{}).
		AddUniqueIndex("question_id", "question_id", "option_num").Error
	if err != nil {
		return fmt.Errorf("failed to add unique index(question_id): %w", err)
	}

	err = db.
		Model(&Administrators{}).
		AddForeignKey("questionnaire_id", "questionnaires(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(administrators.questionnaire_id): %w", err)
	}

	err = db.
		Model(&Options{}).
		AddForeignKey("question_id", "question(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(options.question_id): %w", err)
	}

	err = db.
		Model(&Question{}).
		AddForeignKey("questionnaire_id", "questionnaires(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(question.questionnaire_id): %w", err)
	}

	err = db.
		Model(&Respondents{}).
		AddForeignKey("questionnaire_id", "questionnaires(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(respondents.questionnaire_id): %w", err)
	}

	err = db.
		Model(&Response{}).
		AddForeignKey("response_id", "respondents(response_id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(response.response_id): %w", err)
	}

	err = db.
		Model(&Response{}).
		AddForeignKey("question_id", "question(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(response.question_id): %w", err)
	}

	err = db.
		Model(&ScaleLabels{}).
		AddForeignKey("question_id", "question(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(scale_labels.question_id): %w", err)
	}

	err = db.
		Model(&Targets{}).
		AddForeignKey("questionnaire_id", "questionnaires(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(targets.questionnaire_id): %w", err)
	}

	err = db.
		Model(&Validations{}).
		AddForeignKey("question_id", "question(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("failed to add foreingkey(validations.question_id): %w", err)
	}

	return nil
}
