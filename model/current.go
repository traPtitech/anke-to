package model

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

// Migrations is all db migrations
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v3(),
	}
}

func AllTables() []interface{} {
	return []interface{}{
		&Questionnaires{},
		&Questions{},
		&Respondents{},
		&Responses{},
		&Administrators{},
		&AdministratorUsers{},
		&AdministratorGroups{},
		&Options{},
		&ScaleLabels{},
		&Targets{},
		&TargetUsers{},
		&TargetGroups{},
		&Validations{},
	}
}
