package model

import (
	"fmt"
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
		Migrate: func(tx *gorm.DB) (err error) {
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
			if err := tx.Exec("SET SESSION FOREIGN_KEY_CHECKS = 0").Error; err != nil {
				return err
			}
			defer func() {
				if restoreErr := tx.Exec("SET SESSION FOREIGN_KEY_CHECKS = 1").Error; restoreErr != nil && err == nil {
					err = restoreErr
				}
			}()
			if err := tx.Migrator().RenameTable("question", "questions"); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&v3Questions{}); err != nil {
				return err
			}
			if err := migrateQuestionForeignKeys(tx); err != nil {
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

const (
	v3QuestionsTableName    = "questions"
	v3OldQuestionsTableName = "question"
)

var v3ValidReferentialActions = []string{"CASCADE", "RESTRICT", "SET NULL", "NO ACTION", "SET DEFAULT"}

func (*v3Questions) TableName() string {
	return v3QuestionsTableName
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

func migrateQuestionForeignKeys(tx *gorm.DB) error {
	const fkQuery = `
SELECT kcu.CONSTRAINT_NAME, kcu.TABLE_NAME, kcu.COLUMN_NAME, kcu.REFERENCED_COLUMN_NAME,
       kcu.ORDINAL_POSITION, rc.UPDATE_RULE, rc.DELETE_RULE
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
  ON kcu.CONSTRAINT_SCHEMA = rc.CONSTRAINT_SCHEMA
 AND kcu.CONSTRAINT_NAME = rc.CONSTRAINT_NAME
WHERE kcu.CONSTRAINT_SCHEMA = DATABASE()
  AND kcu.REFERENCED_TABLE_NAME = ?
ORDER BY kcu.CONSTRAINT_NAME, kcu.TABLE_NAME, kcu.ORDINAL_POSITION
`

	var columns []v3QuestionForeignKeyColumn
	if err := tx.Raw(fkQuery, v3OldQuestionsTableName).Scan(&columns).Error; err != nil {
		return err
	}
	if len(columns) == 0 {
		return nil
	}

	type foreignKeyIdentifier struct {
		constraintName string
		tableName      string
	}
	type foreignKeyDefinition struct {
		columns       []string
		refColumns    []string
		updateRule    string
		deleteRule    string
	}

	definitions := map[foreignKeyIdentifier]*foreignKeyDefinition{}
	for _, column := range columns {
		key := foreignKeyIdentifier{
			constraintName: column.ConstraintName,
			tableName:      column.TableName,
		}
		definition, exists := definitions[key]
		if !exists {
			updateRule, err := validateReferentialRule(column.UpdateRule)
			if err != nil {
				return err
			}
			deleteRule, err := validateReferentialRule(column.DeleteRule)
			if err != nil {
				return err
			}
			definition = &foreignKeyDefinition{
				columns:    []string{},
				refColumns: []string{},
				updateRule: updateRule,
				deleteRule: deleteRule,
			}
			definitions[key] = definition
		}
		definition.columns = append(definition.columns, column.ColumnName)
		definition.refColumns = append(definition.refColumns, column.ReferencedColumnName)
	}

	for key, definition := range definitions {
		tableName, err := quoteIdentifier(key.tableName)
		if err != nil {
			return err
		}
		constraintName, err := quoteIdentifier(key.constraintName)
		if err != nil {
			return err
		}
		columnList, err := joinIdentifiers(definition.columns)
		if err != nil {
			return err
		}
		refColumnList, err := joinIdentifiers(definition.refColumns)
		if err != nil {
			return err
		}
		referencedTable, err := quoteIdentifier(v3QuestionsTableName)
		if err != nil {
			return err
		}

		dropSQL := fmt.Sprintf(
			"ALTER TABLE %s DROP FOREIGN KEY %s",
			tableName,
			constraintName,
		)
		if err := tx.Exec(dropSQL).Error; err != nil {
			return err
		}

		addSQL := fmt.Sprintf(
			"ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON UPDATE %s ON DELETE %s",
			tableName,
			constraintName,
			columnList,
			referencedTable,
			refColumnList,
			definition.updateRule,
			definition.deleteRule,
		)
		if err := tx.Exec(addSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func quoteIdentifier(identifier string) (string, error) {
	if strings.ContainsRune(identifier, '\x00') {
		return "", fmt.Errorf("invalid identifier %q", identifier)
	}
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`", nil
}

func joinIdentifiers(identifiers []string) (string, error) {
	quoted := make([]string, len(identifiers))
	for i, identifier := range identifiers {
		quotedIdentifier, err := quoteIdentifier(identifier)
		if err != nil {
			return "", err
		}
		quoted[i] = quotedIdentifier
	}
	return strings.Join(quoted, ", "), nil
}

func validateReferentialRule(rule string) (string, error) {
	for _, action := range v3ValidReferentialActions {
		if rule == action {
			return rule, nil
		}
	}
	return "", fmt.Errorf("invalid referential action %q: must be one of %v", rule, v3ValidReferentialActions)
}
