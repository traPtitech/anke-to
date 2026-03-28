package model

import (
	"fmt"
	"sort"
	"strings"
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
			if err := migrateQuestionTable(tx); err != nil {
				return err
			}
			return nil
		},
	}
}

func migrateQuestionTable(tx *gorm.DB) error {
	hasQuestion := tx.Migrator().HasTable(v3LegacyQuestionsTableName)
	hasQuestions := tx.Migrator().HasTable(v3QuestionsTableName)

	switch {
	case hasQuestion && hasQuestions:
		return fmt.Errorf("both %s and %s tables exist", v3LegacyQuestionsTableName, v3QuestionsTableName)
	case hasQuestion:
		definitions, err := migrateQuestionForeignKeys(tx)
		if err != nil {
			return err
		}
		if err := tx.Migrator().RenameTable(v3LegacyQuestionsTableName, v3QuestionsTableName); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&v3Questions{}); err != nil {
			return err
		}
		if err := recreateQuestionForeignKeys(tx, definitions); err != nil {
			return err
		}
	case hasQuestions:
		if err := tx.AutoMigrate(&v3Questions{}); err != nil {
			return err
		}
	default:
		if err := tx.AutoMigrate(&v3Questions{}); err != nil {
			return err
		}
	}

	return nil
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
	IsPublished              bool             `json:"is_published" gorm:"type:boolean;default:false"`
	IsAnonymous              bool             `json:"is_anonymous" gorm:"type:boolean;not null;default:false"`
	IsDuplicateAnswerAllowed bool             `json:"is_duplicate_answer_allowed" gorm:"type:tinyint(4);size:4;not null;default:0"`
}

func (*v3Questionnaires) TableName() string {
	return "questionnaires"
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
	DeletedAt       gorm.DeletedAt `json:"-"                   gorm:"type:TIMESTAMP NULL;default:NULL"`
	CreatedAt       time.Time      `json:"created_at"          gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (*v3Questions) TableName() string {
	return "questions"
}

type v3TargetUsers struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

func (*v3TargetUsers) TableName() string {
	return "targets_users"
}

type v3TargetGroups struct {
	QuestionnaireID int       `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	GroupID         uuid.UUID `gorm:"type:char(36);size:36;not null;primaryKey"`
}

func (*v3TargetGroups) TableName() string {
	return "targets_groups"
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

type v3QuestionForeignKeyDefinition struct {
	columns    []string
	refColumns []string
	updateRule string
	deleteRule string
	constraint string
	table      string
}

type v3QuestionForeignKeyColumn struct {
	ConstraintName       string `gorm:"column:CONSTRAINT_NAME"`
	TableName            string `gorm:"column:TABLE_NAME"`
	ColumnName           string `gorm:"column:COLUMN_NAME"`
	ReferencedColumnName string `gorm:"column:REFERENCED_COLUMN_NAME"`
	OrdinalPosition      int    `gorm:"column:ORDINAL_POSITION"`
	UpdateRule           string `gorm:"column:UPDATE_RULE"`
	DeleteRule           string `gorm:"column:DELETE_RULE"`
}

const (
	v3LegacyQuestionsTableName = "question"
	v3QuestionsTableName       = "questions"
)

var v3ReferentialActions = map[string]struct{}{
	"CASCADE":  {},
	"RESTRICT": {},
	"SET NULL": {},
	"NO ACTION": {},
}

func migrateQuestionForeignKeys(tx *gorm.DB) ([]v3QuestionForeignKeyDefinition, error) {
	definitions, err := loadLegacyQuestionForeignKeys(tx)
	if err != nil {
		return nil, err
	}
	for _, definition := range definitions {
		dropSQL := fmt.Sprintf(
			"ALTER TABLE %s DROP FOREIGN KEY %s",
			quoteIdentifier(definition.table),
			quoteIdentifier(definition.constraint),
		)
		if err := tx.Exec(dropSQL).Error; err != nil {
			return nil, err
		}
	}

	return definitions, nil
}

func recreateQuestionForeignKeys(tx *gorm.DB, definitions []v3QuestionForeignKeyDefinition) error {
	if len(definitions) == 0 {
		return nil
	}

	sort.Slice(definitions, func(i, j int) bool {
		if definitions[i].table == definitions[j].table {
			return definitions[i].constraint < definitions[j].constraint
		}
		return definitions[i].table < definitions[j].table
	})

	for _, definition := range definitions {
		addSQL := fmt.Sprintf(
			"ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON UPDATE %s ON DELETE %s",
			quoteIdentifier(definition.table),
			quoteIdentifier(definition.constraint),
			joinIdentifiers(definition.columns),
			quoteIdentifier(v3QuestionsTableName),
			joinIdentifiers(definition.refColumns),
			definition.updateRule,
			definition.deleteRule,
		)
		if err := tx.Exec(addSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func loadLegacyQuestionForeignKeys(tx *gorm.DB) ([]v3QuestionForeignKeyDefinition, error) {
	const fkQuery = `
SELECT
	kcu.CONSTRAINT_NAME,
	kcu.TABLE_NAME,
	kcu.COLUMN_NAME,
	kcu.REFERENCED_COLUMN_NAME,
	kcu.ORDINAL_POSITION,
	rc.UPDATE_RULE,
	rc.DELETE_RULE
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
	ON kcu.CONSTRAINT_SCHEMA = rc.CONSTRAINT_SCHEMA
	AND kcu.CONSTRAINT_NAME = rc.CONSTRAINT_NAME
WHERE kcu.CONSTRAINT_SCHEMA = DATABASE()
	AND kcu.REFERENCED_TABLE_NAME = ?
ORDER BY kcu.CONSTRAINT_NAME, kcu.TABLE_NAME, kcu.ORDINAL_POSITION
`

	var columns []v3QuestionForeignKeyColumn
	if err := tx.Raw(fkQuery, v3LegacyQuestionsTableName).Scan(&columns).Error; err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, nil
	}

	type fkKey struct {
		constraintName string
		tableName      string
	}

	definitions := map[fkKey]*v3QuestionForeignKeyDefinition{}
	for _, column := range columns {
		key := fkKey{
			constraintName: column.ConstraintName,
			tableName:      column.TableName,
		}
		definition, exists := definitions[key]
		if !exists {
			updateRule, err := validateReferentialRule(column.UpdateRule)
			if err != nil {
				return nil, err
			}
			deleteRule, err := validateReferentialRule(column.DeleteRule)
			if err != nil {
				return nil, err
			}
			definition = &v3QuestionForeignKeyDefinition{
				columns:    []string{},
				refColumns: []string{},
				updateRule: updateRule,
				deleteRule: deleteRule,
				constraint: column.ConstraintName,
				table:      column.TableName,
			}
			definitions[key] = definition
		}
		definition.columns = append(definition.columns, column.ColumnName)
		definition.refColumns = append(definition.refColumns, column.ReferencedColumnName)
	}

	result := make([]v3QuestionForeignKeyDefinition, 0, len(definitions))
	for _, definition := range definitions {
		result = append(result, *definition)
	}

	return result, nil
}

func quoteIdentifier(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}

func joinIdentifiers(identifiers []string) string {
	quoted := make([]string, len(identifiers))
	for i, identifier := range identifiers {
		quoted[i] = quoteIdentifier(identifier)
	}
	return strings.Join(quoted, ", ")
}

func validateReferentialRule(rule string) (string, error) {
	if _, ok := v3ReferentialActions[rule]; ok {
		return rule, nil
	}
	return "", fmt.Errorf("unexpected referential action: %s", rule)
}
