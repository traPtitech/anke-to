package model

import (
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

// Questionnaire QuestionnaireRepositoryの実装
type Questionnaire struct{}

// NewQuestionnaire Questionnaireのコンストラクター
func NewQuestionnaire() *Questionnaire {
	return new(Questionnaire)
}

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
func (questionnaire *Questionnaires) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set ModifiedAt: %w", err)
	}
	err = scope.SetColumn("CreatedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set CreatedAt: %w", err)
	}

	return nil
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

//QuestionnaireDetail Questionnaireの詳細
type QuestionnaireDetail struct {
	Targets        []string
	Respondents    []string
	Administrators []string
	Questionnaires
}

//TargettedQuestionnaire targetになっているアンケートの情報
type TargettedQuestionnaire struct {
	Questionnaires
	RespondedAt null.Time `json:"responded_at"`
	HasResponse bool      `json:"has_response"`
}

type ResponseReadPrivilegeInfo struct {
	ResSharedTo     string
	IsAdministrator bool
	IsRespondent    bool
}

//InsertQuestionnaire アンケートの追加
func (*Questionnaire) InsertQuestionnaire(title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error) {
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
func (*Questionnaire) UpdateQuestionnaire(title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error {
	if !resTimeLimit.Valid {
		questionnaire := map[string]interface{}{
			"title":          title,
			"description":    description,
			"res_time_limit": gorm.Expr("NULL"),
			"res_shared_to":  resSharedTo,
		}

		result := db.
			Model(&Questionnaires{}).
			Where("id = ?", questionnaireID).
			Update(questionnaire)
		err := result.Error
		if err != nil {
			return fmt.Errorf("failed to update a questionnaire record: %w", err)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("failed to update a questionnaire record: %w", ErrNoRecordUpdated)
		}

		return nil
	}

	questionnaire := Questionnaires{
		Title:        title,
		Description:  description,
		ResTimeLimit: resTimeLimit,
		ResSharedTo:  resSharedTo,
	}

	result := db.
		Model(&Questionnaires{}).
		Where("id = ?", questionnaireID).
		Update(&questionnaire)
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to update a questionnaire record: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update a questionnaire record: %w", ErrNoRecordUpdated)
	}

	return nil
}

//DeleteQuestionnaire アンケートの削除
func (*Questionnaire) DeleteQuestionnaire(questionnaireID int) error {
	result := db.Delete(&Questionnaires{ID: questionnaireID})
	err := result.Error
	if err != nil {
		return fmt.Errorf("failed to delete questionnaire: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete questionnaire: %w", ErrNoRecordDeleted)
	}

	return nil
}

/*GetQuestionnaires アンケートの一覧
2つ目の戻り値はページ数の最大値*/
func (*Questionnaire) GetQuestionnaires(userID string, sort string, search string, pageNum int, nontargeted bool) ([]QuestionnaireInfo, int, error) {
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
	if len(search) != 0 {
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
		return []QuestionnaireInfo{}, 0, nil
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
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, 0, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaires, pageMax, nil
}

// GetAdminQuestionnaires 自分が管理者のアンケートの取得
func (*Questionnaire) GetAdminQuestionnaires(userID string) ([]Questionnaires, error) {
	questionnaires := []Questionnaires{}
	err := db.
		Table("questionnaires").
		Joins("INNER JOIN administrators ON questionnaires.id = administrators.questionnaire_id").
		Where("administrators.user_traqid = ?", userID).
		Order("questionnaires.modified_at DESC").
		Find(&questionnaires).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get a questionnaire: %w", err)
	}

	return questionnaires, nil
}

//GetQuestionnaireInfo アンケートの詳細な情報取得
func (*Questionnaire) GetQuestionnaireInfo(questionnaireID int) (*Questionnaires, []string, []string, []string, error) {
	questionnaire := Questionnaires{}
	targets := []string{}
	administrators := []string{}
	respondents := []string{}

	err := db.
		Model(&Questionnaires{}).
		Where("questionnaires.id = ?", questionnaireID).
		First(&questionnaire).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get a questionnaire: %w", err)
	}

	err = db.
		Table("targets").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &targets).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get targets: %w", err)
	}

	err = db.
		Table("administrators").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &administrators).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get administrators: %w", err)
	}

	err = db.
		Table("respondents").
		Where("questionnaire_id = ? AND deleted_at IS NULL AND submitted_at IS NOT NULL", questionnaire.ID).
		Pluck("user_traqid", &respondents).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get respondents: %w", err)
	}

	return &questionnaire, targets, administrators, respondents, nil
}

//GetTargettedQuestionnaires targetになっているアンケートの取得
func (*Questionnaire) GetTargettedQuestionnaires(userID string, answered string, sort string) ([]TargettedQuestionnaire, error) {
	query := db.
		Table("questionnaires").
		Where("questionnaires.res_time_limit > ? OR questionnaires.res_time_limit IS NULL", time.Now()).
		Joins("INNER JOIN targets ON questionnaires.id = targets.questionnaire_id").
		Where("targets.user_traqid = ? OR targets.user_traqid = 'traP'", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id AND respondents.user_traqid = ? AND respondents.deleted_at IS NULL", userID).
		Group("questionnaires.id,respondents.user_traqid").
		Select("questionnaires.*, MAX(respondents.submitted_at) AS responded_at, COUNT(respondents.response_id) != 0 AS has_response")

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
		return nil, fmt.Errorf("invalid answered parameter value(%s): %w", answered, ErrInvalidAnsweredParam)
	}

	questionnaires := []TargettedQuestionnaire{}
	err = query.Find(&questionnaires).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaires, nil
}

//GetQuestionnaireLimit アンケートの回答期限の取得
func (*Questionnaire) GetQuestionnaireLimit(questionnaireID int) (null.Time, error) {
	res := Questionnaires{}

	err := db.
		Model(Questionnaires{}).
		Where("id = ?", questionnaireID).
		Select("res_time_limit").
		Scan(&res).Error
	if err != nil {
		return null.NewTime(time.Time{}, false), fmt.Errorf("failed to get the questionnaires: %w", err)
	}

	return res.ResTimeLimit, nil
}

// GetQuestionnaireLimitByResponseID 回答のIDからアンケートの回答期限を取得
func (*Questionnaire) GetQuestionnaireLimitByResponseID(responseID int) (null.Time, error) {
	res := Questionnaires{}

	err := db.
		Table("respondents").
		Joins("INNER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Where("respondents.response_id = ? AND respondents.deleted_at IS NULL", responseID).
		Select("questionnaires.res_time_limit").
		Scan(&res).Error
	if err != nil {
		return null.NewTime(time.Time{}, false), fmt.Errorf("failed to get the questionnaires: %w", err)
	}

	return res.ResTimeLimit, nil
}

func (*Questionnaire) GetResponseReadPrivilegeInfoByResponseID(userID string, responseID int) (*ResponseReadPrivilegeInfo, error) {
	responseReadPrivilegeInfo := ResponseReadPrivilegeInfo{}
	err := db.
		Table("respondents").
		Where("respondents.response_id = ? AND respondents.submitted_at IS NOT NULL", responseID).
		Joins("INNER JOIN questionnaires ON questionnaires.id = respondents.questionnaire_id").
		Joins("LEFT OUTER JOIN administrators ON questionnaires.id = administrators.questionnaire_id AND administrators.user_traqid = ?", userID).
		Joins("LEFT OUTER JOIN respondents AS respondents2 ON questionnaires.id = respondents2.questionnaire_id AND respondents2.user_traqid = ? AND respondents2.submitted_at IS NOT NULL", userID).
		Select("questionnaires.res_shared_to, administrators.questionnaire_id IS NOT NULL AS is_administrator, respondents2.response_id IS NOT NULL AS is_respondent").
		Scan(&responseReadPrivilegeInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get response read privilege info: %w", err)
	}

	return &responseReadPrivilegeInfo, nil
}

func (*Questionnaire) GetResponseReadPrivilegeInfoByQuestionnaireID(userID string, questionnaireID int) (*ResponseReadPrivilegeInfo, error) {
	responseReadPrivilegeInfo := ResponseReadPrivilegeInfo{}
	err := db.
		Table("questionnaires").
		Where("questionnaires.id = ?", questionnaireID).
		Joins("LEFT OUTER JOIN administrators ON questionnaires.id = administrators.questionnaire_id AND administrators.user_traqid = ?", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id AND respondents.user_traqid = ? AND respondents.submitted_at IS NOT NULL", userID).
		Select("questionnaires.res_shared_to, administrators.questionnaire_id IS NOT NULL AS is_administrator, respondents.response_id IS NOT NULL AS is_respondent").
		Scan(&responseReadPrivilegeInfo).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get response read privilege info: %w", err)
	}

	return &responseReadPrivilegeInfo, nil
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
		return nil, ErrInvalidSortParam
	}
	query = query.Order("questionnaires.id desc")

	return query, nil
}
