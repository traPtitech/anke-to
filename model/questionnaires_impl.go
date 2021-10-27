package model

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// Questionnaire QuestionnaireRepositoryの実装
type Questionnaire struct{}

// NewQuestionnaire Questionnaireのコンストラクター
func NewQuestionnaire() *Questionnaire {
	return new(Questionnaire)
}

//Questionnaires questionnairesテーブルの構造体
type Questionnaires struct {
	ID             int              `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title          string           `json:"title"           gorm:"type:char(50);size:50;not null"`
	Description    string           `json:"description"     gorm:"type:text;not null"`
	ResTimeLimit   null.Time        `json:"res_time_limit,omitempty"  gorm:"type:TIMESTAMP NULL;default:NULL;"`
	DeletedAt      gorm.DeletedAt   `json:"-"      gorm:"type:TIMESTAMP NULL;default:NULL;"`
	ResSharedTo    string           `json:"res_shared_to"   gorm:"type:char(30);size:30;not null;default:administrators"`
	CreatedAt      time.Time        `json:"created_at"      gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt     time.Time        `json:"modified_at"     gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Administrators []Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets        []Targets        `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Questions      []Questions      `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Respondents    []Respondents    `json:"-"  gorm:"foreignKey:QuestionnaireID"`
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (questionnaire *Questionnaires) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	questionnaire.ModifiedAt = now
	questionnaire.CreatedAt = now

	return nil
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (questionnaire *Questionnaires) BeforeUpdate(tx *gorm.DB) error {
	questionnaire.ModifiedAt = time.Now()

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
func (*Questionnaire) InsertQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error) {
	db, err := getTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get tx: %w", err)
	}

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

	err = db.Create(&questionnaire).Error
	if err != nil {
		return 0, fmt.Errorf("failed to insert a questionnaire: %w", err)
	}

	return questionnaire.ID, nil
}

//UpdateQuestionnaire アンケートの更新
func (*Questionnaire) UpdateQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}

	var questionnaire interface{}
	if resTimeLimit.Valid {
		questionnaire = Questionnaires{
			Title:        title,
			Description:  description,
			ResTimeLimit: resTimeLimit,
			ResSharedTo:  resSharedTo,
		}
	} else {
		questionnaire = map[string]interface{}{
			"title":          title,
			"description":    description,
			"res_time_limit": gorm.Expr("NULL"),
			"res_shared_to":  resSharedTo,
		}
	}

	result := db.
		Model(&Questionnaires{}).
		Where("id = ?", questionnaireID).
		Updates(questionnaire)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update a questionnaire record: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update a questionnaire record: %w", ErrNoRecordUpdated)
	}

	return nil
}

//DeleteQuestionnaire アンケートの削除
func (*Questionnaire) DeleteQuestionnaire(ctx context.Context, questionnaireID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}

	result := db.Delete(&Questionnaires{ID: questionnaireID})
	err = result.Error
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
func (*Questionnaire) GetQuestionnaires(ctx context.Context, userID string, sort string, search string, pageNum int, nontargeted bool) ([]QuestionnaireInfo, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	db, err := getTx(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tx: %w", err)
	}

	questionnaires := make([]QuestionnaireInfo, 0, 20)

	query := db.
		Table("questionnaires").
		Joins("LEFT OUTER JOIN targets ON questionnaires.id = targets.questionnaire_id")

	query, err = setQuestionnairesOrder(query, sort)
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

	var count int64
	err = query.
		Session(&gorm.Session{}).
		Group("questionnaires.id").
		Count(&count).Error
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, 0, ErrDeadlineExceeded
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve the number of questionnaires: %w", err)
	}

	if count == 0 {
		return []QuestionnaireInfo{}, 0, nil
	}
	pageMax := (int(count) + 19) / 20

	if pageNum > pageMax {
		return nil, 0, fmt.Errorf("failed to set page offset: %w", ErrTooLargePageNum)
	}

	offset := (pageNum - 1) * 20

	err = query.
		Limit(20).
		Offset(offset).
		Group("questionnaires.id").
		Select("questionnaires.*, (targets.user_traqid = ? OR targets.user_traqid = 'traP') AS is_targeted", userID).
		Find(&questionnaires).Error
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, 0, ErrDeadlineExceeded
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaires, pageMax, nil
}

// GetAdminQuestionnaires 自分が管理者のアンケートの取得
func (*Questionnaire) GetAdminQuestionnaires(ctx context.Context, userID string) ([]Questionnaires, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	questionnaires := []Questionnaires{}
	err = db.
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
func (*Questionnaire) GetQuestionnaireInfo(ctx context.Context, questionnaireID int) (*Questionnaires, []string, []string, []string, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get tx: %w", err)
	}

	questionnaire := Questionnaires{}
	targets := []string{}
	administrators := []string{}
	respondents := []string{}

	err = db.
		Where("questionnaires.id = ?", questionnaireID).
		First(&questionnaire).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get a questionnaire: %w", err)
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Table("targets").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &targets).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get targets: %w", err)
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Table("administrators").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &administrators).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get administrators: %w", err)
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Table("respondents").
		Where("questionnaire_id = ? AND deleted_at IS NULL AND submitted_at IS NOT NULL", questionnaire.ID).
		Pluck("user_traqid", &respondents).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get respondents: %w", err)
	}

	return &questionnaire, targets, administrators, respondents, nil
}

//GetTargettedQuestionnaires targetになっているアンケートの取得
func (*Questionnaire) GetTargettedQuestionnaires(ctx context.Context, userID string, answered string, sort string) ([]TargettedQuestionnaire, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	query := db.
		Table("questionnaires").
		Where("questionnaires.res_time_limit > ? OR questionnaires.res_time_limit IS NULL", time.Now()).
		Joins("INNER JOIN targets ON questionnaires.id = targets.questionnaire_id").
		Where("targets.user_traqid = ? OR targets.user_traqid = 'traP'", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id AND respondents.user_traqid = ? AND respondents.deleted_at IS NULL", userID).
		Group("questionnaires.id,respondents.user_traqid").
		Select("questionnaires.*, MAX(respondents.submitted_at) AS responded_at, COUNT(respondents.response_id) != 0 AS has_response")

	query, err = setQuestionnairesOrder(query, sort)
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
func (*Questionnaire) GetQuestionnaireLimit(ctx context.Context, questionnaireID int) (null.Time, error) {
	db, err := getTx(ctx)
	if err != nil {
		return null.NewTime(time.Time{}, false), fmt.Errorf("failed to get tx: %w", err)
	}

	var res Questionnaires

	err = db.
		Where("id = ?", questionnaireID).
		Select("res_time_limit").
		First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return null.NewTime(time.Time{}, false), ErrRecordNotFound
	}
	if err != nil {
		return null.NewTime(time.Time{}, false), fmt.Errorf("failed to get the questionnaires: %w", err)
	}

	return res.ResTimeLimit, nil
}

// GetQuestionnaireLimitByResponseID 回答のIDからアンケートの回答期限を取得
func (*Questionnaire) GetQuestionnaireLimitByResponseID(ctx context.Context, responseID int) (null.Time, error) {
	db, err := getTx(ctx)
	if err != nil {
		return null.NewTime(time.Time{}, false), fmt.Errorf("failed to get tx: %w", err)
	}

	var res Questionnaires

	err = db.
		Joins("INNER JOIN respondents ON questionnaires.id = respondents.questionnaire_id").
		Where("respondents.response_id = ? AND respondents.deleted_at IS NULL", responseID).
		Select("questionnaires.res_time_limit").
		First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return null.NewTime(time.Time{}, false), ErrRecordNotFound
	}
	if err != nil {
		return null.NewTime(time.Time{}, false), fmt.Errorf("failed to get the questionnaires: %w", err)
	}

	return res.ResTimeLimit, nil
}

func (*Questionnaire) GetResponseReadPrivilegeInfoByResponseID(ctx context.Context, userID string, responseID int) (*ResponseReadPrivilegeInfo, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	var responseReadPrivilegeInfo ResponseReadPrivilegeInfo
	err = db.
		Table("respondents").
		Where("respondents.response_id = ? AND respondents.submitted_at IS NOT NULL", responseID).
		Joins("INNER JOIN questionnaires ON questionnaires.id = respondents.questionnaire_id").
		Joins("LEFT OUTER JOIN administrators ON questionnaires.id = administrators.questionnaire_id AND administrators.user_traqid = ?", userID).
		Joins("LEFT OUTER JOIN respondents AS respondents2 ON questionnaires.id = respondents2.questionnaire_id AND respondents2.user_traqid = ? AND respondents2.submitted_at IS NOT NULL", userID).
		Select("questionnaires.res_shared_to, administrators.questionnaire_id IS NOT NULL AS is_administrator, respondents2.response_id IS NOT NULL AS is_respondent").
		Take(&responseReadPrivilegeInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get response read privilege info: %w", err)
	}

	return &responseReadPrivilegeInfo, nil
}

func (*Questionnaire) GetResponseReadPrivilegeInfoByQuestionnaireID(ctx context.Context, userID string, questionnaireID int) (*ResponseReadPrivilegeInfo, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	var responseReadPrivilegeInfo ResponseReadPrivilegeInfo
	err = db.
		Table("questionnaires").
		Where("questionnaires.id = ?", questionnaireID).
		Joins("LEFT OUTER JOIN administrators ON questionnaires.id = administrators.questionnaire_id AND administrators.user_traqid = ?", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id AND respondents.user_traqid = ? AND respondents.submitted_at IS NOT NULL", userID).
		Select("questionnaires.res_shared_to, administrators.questionnaire_id IS NOT NULL AS is_administrator, respondents.response_id IS NOT NULL AS is_respondent").
		Take(&responseReadPrivilegeInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
