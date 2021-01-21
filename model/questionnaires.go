package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

var (
	// ErrTooLargePageNum too large page number
	ErrTooLargePageNum = errors.New("too large page number")
	// ErrInvalidRegex invalid regexp
	ErrInvalidRegex = errors.New("invalid regexp")
)

//Questionnaires questionnairesテーブルの構造体
type Questionnaires struct {
	ID           int       `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT NOT NULL PRIMARY KEY;"`
	Title        string    `json:"title"           gorm:"type:char(50) NOT NULL;"`
	Description  string    `json:"description"     gorm:"type:text NOT NULL;"`
	ResTimeLimit null.Time `json:"res_time_limit,omitempty"  gorm:"type:timestamp NULL;default:NULL;"`
	DeletedAt    null.Time `json:"deleted_at,omitempty"      gorm:"type:timestamp NULL;default:NULL;"`
	ResSharedTo  string    `json:"res_shared_to"   gorm:"type:char(30) NOT NULL;default:\"administrators\";"`
	CreatedAt    time.Time `json:"created_at"      gorm:"type:timestamp NOT NULL;default:CURRENT_TIMESTAMP;"`
	ModifiedAt   time.Time `json:"modified_at"     gorm:"type:timestamp NOT NULL;default:CURRENT_TIMESTAMP;"`
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (questionnaire *Questionnaires) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set ModifiedAt: %w", err)
	}

	return nil
}

//QuestionnaireInfo Questionnaireにtargetかの情報追加
type QuestionnaireInfo struct {
	Questionnaires
	IsTargeted bool `json:"is_targeted" gorm:"type:boolean"`
}

type QuestionnaireDetail struct {
	Targets        []string
	Respondents    []string
	Administrators []string
	Questionnaires
}

//TargettedQuestionnaire targetになっているアンケートの情報
type TargettedQuestionnaire struct {
	Questionnaires
	RespondedAt string `json:"responded_at"`
}

//InsertQuestionnaire アンケートの追加
func InsertQuestionnaire(title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error) {
	var questionnaire Questionnaires
	if !resTimeLimit.Valid {
		questionnaire = Questionnaires{
			Title:       title,
			Description: description,
			ResSharedTo: resSharedTo,
		}
	} else {
		questionnaire = Questionnaires{
			Title:        title,
			Description:  description,
			ResTimeLimit: resTimeLimit,
			ResSharedTo:  resSharedTo,
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&questionnaire).Error
		if err != nil {
			return fmt.Errorf("failed to insert a questionnaire: %w", err)
		}

		err = tx.
			Select("id").
			Last(&questionnaire).Error
		if err != nil {
			return fmt.Errorf("failed to get the last id: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed in the transaction: %w", err)
	}

	return questionnaire.ID, nil
}

//UpdateQuestionnaire アンケートの更新
func UpdateQuestionnaire(title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error {
	if !resTimeLimit.Valid {
		questionnaire := map[string]interface{}{
			"title":          title,
			"description":    description,
			"res_time_limit": gorm.Expr("NULL"),
			"res_shared_to":  resSharedTo,
		}

		err := db.
			Model(&Questionnaires{}).
			Where("id = ?", questionnaireID).
			Update(questionnaire).Error
		if err != nil {
			return fmt.Errorf("failed to update a questionnaire record: %w", err)
		}

		return nil
	}

	questionnaire := Questionnaires{
		Title:        title,
		Description:  description,
		ResTimeLimit: resTimeLimit,
		ResSharedTo:  resSharedTo,
	}

	err := db.
		Model(&Questionnaires{}).
		Where("id = ?", questionnaireID).
		Update(&questionnaire).Error
	if err != nil {
		return fmt.Errorf("failed to update a questionnaire record: %w", err)
	}

	return nil
}

//DeleteQuestionnaire アンケートの削除
func DeleteQuestionnaire(questionnaireID int) error {
	err := db.Delete(&Questionnaires{ID: questionnaireID}).Error
	if err != nil {
		return fmt.Errorf("failed to delete questionnaire: %w", err)
	}

	return nil
}

/*GetQuestionnaires アンケートの一覧
2つ目の戻り値はページ数の最大値*/
func GetQuestionnaires(userID string, sort string, search string, pageNum int, nontargeted bool) ([]QuestionnaireInfo, int, error) {
	questionnaires := make([]QuestionnaireInfo, 0, 20)

	query := db.
		Table("questionnaires").
		Joins("LEFT OUTER JOIN targets ON questionnaires.id = targets.questionnaire_id")

	query, err := setQuestionnairesOrder(query, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to set the order of the questionnaire table: %w", err)
	}

	if nontargeted {
		query = query.Where("targets.questionnaire_id IS NULL OR (targets.user_traqid != ? AND targets.user_traqid != 'traP')", userID)
	}
	if search != "" {
		// MySQLでのregexpの構文は少なくともGoのregexpの構文でvalidである必要がある
		_, err := regexp.Compile(search)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid search param: %w", ErrInvalidRegex)
		}

		// BINARYをつけていないので大文字小文字区別しない
		query = query.Where("questionnaires.title REGEXP ?", search)
	}

	count := 0
	err = query.
		Group("questionnaires.id").
		Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve the number of questionnaires: %w", err)
	}
	if count == 0 {
		return questionnaires, 0, nil
	}
	pageMax := (count + 19) / 20

	if pageNum > pageMax {
		return nil, 0, fmt.Errorf("failed to set page offset: %w", ErrTooLargePageNum)
	}

	offset := (pageNum - 1) * 20
	query = query.Limit(20).Offset(offset)

	err = query.
		Group("questionnaires.id").
		Select("questionnaires.*, (targets.user_traqid = ? OR targets.user_traqid = 'traP') AS is_targeted", userID).
		Find(&questionnaires).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, 0, fmt.Errorf("failed to get the targeted questionnaires: %w", gorm.ErrRecordNotFound)
	} else if err != nil {
		return nil, 0, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaires, pageMax, nil
}

// GetAdminQuestionnaires 自分が管理者のアンケートの取得
func GetAdminQuestionnaires(userID string) ([]Questionnaires, error) {
	questionnaires := []Questionnaires{}
	err := db.
		Table("questionnaires").
		Joins("INNER JOIN administrators ON questionnaires.id = administrators.questionnaire_id").
		Where("administrators.user_traqid = ?", userID).
		Order("questionnaires.modified_at DESC").
		Find(&questionnaires).Error
	if gorm.IsRecordNotFoundError(err) {
		return []Questionnaires{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get a questionnaire: %w", err)
	}

	return questionnaires, nil
}

//GetQuestionnaireInfo アンケートの詳細な情報取得
func GetQuestionnaireInfo(questionnaireID int) (*Questionnaires, []string, []string, []string, error) {
	questionnaire := Questionnaires{}
	targets := []string{}
	administrators := []string{}
	respondents := []string{}

	err := db.
		Model(&Questionnaires{}).
		Where("questionnaires.id = ?", questionnaireID).
		First(&questionnaire).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, fmt.Errorf("failed to get a questionnaire: %w", gorm.ErrRecordNotFound)
		}
		return nil, nil, nil, nil, fmt.Errorf("failed to get a questionnaire: %w", err)
	}

	err = db.
		Table("targets").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &targets).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, fmt.Errorf("failed to get targets: %w", gorm.ErrRecordNotFound)
		}
		return nil, nil, nil, nil, fmt.Errorf("failed to get targets: %w", err)
	}

	err = db.
		Table("administrators").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &administrators).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, fmt.Errorf("failed to get administrators: %w", gorm.ErrRecordNotFound)
		}
		return nil, nil, nil, nil, fmt.Errorf("failed to get administrators: %w", err)
	}

	err = db.
		Table("respondents").
		Where("questionnaire_id = ? AND deleted_at IS NULL AND submitted_at IS NOT NULL", questionnaire.ID).
		Pluck("user_traqid", &respondents).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, fmt.Errorf("failed to get respondents: %w", gorm.ErrRecordNotFound)
		}
		return nil, nil, nil, nil, fmt.Errorf("failed to get respondents: %w", err)
	}

	return &questionnaire, targets, administrators, respondents, nil
}

//GetTargettedQuestionnaires targetになっているアンケートの取得
func GetTargettedQuestionnaires(userID string, answered string, sort string) ([]TargettedQuestionnaire, error) {
	query := db.
		Table("questionnaires").
		Where("questionnaires.res_time_limit > ? OR questionnaires.res_time_limit IS NULL", time.Now()).
		Joins("INNER JOIN targets ON questionnaires.id = targets.questionnaire_id").
		Where("targets.user_traqid = ? OR targets.user_traqid = 'traP'", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id AND respondents.user_traqid = ?", userID).
		Where("respondents.user_traqid IS NULL").
		Group("questionnaires.id,respondents.user_traqid").
		Select("questionnaires.*, MAX(respondents.submitted_at) AS responsed_at")

	query, err := setQuestionnairesOrder(query, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to set the order of the questionnaire table: %w", err)
	}

	query = query.
		Order("questionnaires.res_time_limit").
		Order("questionnaires.modified_at desc")

	switch answered {
	case "answered":
		query = query.Where("respondents.questionnaire_id IS NOT NULL")
	case "unanswered":
		query = query.Where("respondents.questionnaire_id IS NULL")
	case "":
	default:
		return nil, fmt.Errorf("invalid answered parameter value(%s)", answered)
	}

	questionnaires := []TargettedQuestionnaire{}
	err = query.Find(&questionnaires).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("failed to get the targeted questionnaires: %w", gorm.ErrRecordNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaires, nil
}

//GetQuestionnaireLimit アンケートの回答期限の取得
func GetQuestionnaireLimit(questionnaireID int) (string, error) {
	res := Questionnaires{}

	err := db.
		Model(Questionnaires{}).
		Where("id = ?", questionnaireID).
		Select("res_time_limit").
		Scan(&res).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get the questionnaires: %w", err)
	}

	return NullTimeToString(res.ResTimeLimit), nil
}

//GetResShared アンケートの回答の公開範囲の取得
func GetResShared(questionnaireID int) (string, error) {
	res := Questionnaires{}

	err := db.
		Model(Questionnaires{}).
		Where("id = ?", questionnaireID).
		Select("res_shared_to").
		Scan(&res).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", fmt.Errorf("failed to get resShared: %w", gorm.ErrRecordNotFound)
		}
		return "", fmt.Errorf("failed to get resShared: %w", err)
	}

	return res.ResSharedTo, nil
}

func setQuestionnairesOrder(query *gorm.DB, sort string) (*gorm.DB, error) {
	switch sort {
	case "created_at":
		query = query.Order("questionnaires.created_at")
	case "-created_at":
		query = query.Order("questionnaires.created_at desc")
	case "title":
		query = query.Order("questionnaires.title")
	case "-title":
		query = query.Order("questionnaires.title desc")
	case "modified_at":
		query = query.Order("questionnaires.modified_at")
	case "-modified_at":
		query = query.Order("questionnaires.modified_at desc")
	case "":
	default:
		return nil, errors.New("invalid sort type")
	}

	return query, nil
}
