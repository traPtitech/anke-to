package model

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/traPtitech/anke-to/model/migration"
)

// Migrations is all db migrations
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migration.V1(), // Questionnariesにis_anonymousカラムを追加
	}
}

func AllTables() []interface{} {
	return []interface{}{
		&Questionnaires{},
		&Questions{},
		&Respondents{},
		&Responses{},
		&Administrators{},
		&Options{},
		&ScaleLabels{},
		&Targets{},
		&Validations{},
	}
}
