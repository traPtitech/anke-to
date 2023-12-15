package model

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

// Migrations is all db migrations
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{},
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
