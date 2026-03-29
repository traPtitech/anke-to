package model

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// Respondent RespondentRepositoryの実装
type Respondent struct{}

// NewRespondent Respondentのコンストラクター
func NewRespondent() *Respondent {
	return new(Respondent)
}

// Respondents respondentsテーブルの構造体
type Respondents struct {
	ResponseID      int            `json:"responseID" gorm:"column:response_id;type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID" gorm:"type:int(11);not null"`
	UserTraqid      string         `json:"user_traq_id,omitempty" gorm:"type:varchar(32);size:32;default:NULL"`
	ModifiedAt      time.Time      `json:"modified_at,omitempty" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty" gorm:"type:TIMESTAMP NULL;default:NULL"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:TIMESTAMP NULL;default:NULL"`
	Responses       []Responses    `json:"-"  gorm:"foreignKey:ResponseID;references:ResponseID"`
}

// BeforeCreate insert時に自動でmodifiedAt更新
func (r *Respondents) BeforeCreate(_ *gorm.DB) error {
	r.ModifiedAt = time.Now()

	return nil
}

// BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (r *Respondents) BeforeUpdate(_ *gorm.DB) error {
	r.ModifiedAt = time.Now()

	return nil
}

// RespondentInfo 回答とその周辺情報の構造体
type RespondentInfo struct {
	Title        string    `json:"questionnaire_title"`
	Description  string    `json:"description"`
	ResTimeLimit null.Time `json:"res_time_limit"`
	Respondents
}

// RespondentDetail 回答の詳細情報の構造体
type RespondentDetail struct {
	ResponseID      int            `json:"responseID,omitempty"`
	TraqID          string         `json:"traqID,omitempty"`
	QuestionnaireID int            `json:"questionnaireID,omitempty"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty"`
	ModifiedAt      time.Time      `json:"modified_at,omitempty"`
	Responses       []ResponseBody `json:"body"`
}

// InsertRespondent 回答の追加
func (*Respondent) InsertRespondent(ctx context.Context, userID string, questionnaireID int, submittedAt null.Time) (int, error) {
	db, err := getTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get tx: %w", err)
	}

	var questionnaire Questionnaires
	var respondent Respondents

	err = db.
		Where("id = ?", questionnaireID).
		First(&questionnaire).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrRecordNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get questionnaire: %w", err)
	}

	if submittedAt.Valid {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
			SubmittedAt:     submittedAt,
		}
	} else {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
		}
	}

	if !questionnaire.IsDuplicateAnswerAllowed {
		err = db.
			Where("questionnaire_id = ? AND user_traqid = ?", questionnaireID, userID).
			First(&Respondents{}).Error
		if err == nil {
			return 0, ErrDuplicatedAnswered
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("failed to check duplicate answer: %w", err)
		}

	}

	err = db.Create(&respondent).Error
	if err != nil {
		return 0, fmt.Errorf("failed to insert a respondent record: %w", err)
	}

	return respondent.ResponseID, nil

}

// UpdateSubmittedAt 投稿日時更新
func (*Respondent) UpdateSubmittedAt(ctx context.Context, responseID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}

	err = db.
		Model(&Respondents{}).
		Where("response_id = ?", responseID).
		Update("submitted_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("failed to update response's submitted_at: %w", err)
	}

	return nil
}

// UpdateModifiedAt 編集日時更新
func (*Respondent) UpdateModifiedAt(ctx context.Context, responseID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}

	err = db.
		Model(&Respondents{}).
		Where("response_id = ?", responseID).
		Update("modified_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("failed to update response's modified_at: %w", err)
	}

	return nil
}

// DeleteRespondent 回答の削除
func (*Respondent) DeleteRespondent(ctx context.Context, responseID int) error {
	db, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}

	result := db.
		Where("response_id = ?", responseID).
		Delete(&Respondents{})
	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete respondent: %w", err)
	}

	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

// CheckRespondentByResponseID 回答者かどうかの確認
func (*Respondent) GetRespondent(ctx context.Context, responseID int) (*Respondents, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	var respondent Respondents

	err = db.
		Where("response_id = ?", responseID).
		First(&respondent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	return &respondent, nil
}

// GetRespondentInfos ユーザーの回答とその周辺情報一覧の取得
func (*Respondent) GetRespondentInfos(ctx context.Context, userID string, questionnaireIDs ...int) ([]RespondentInfo, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	respondentInfos := []RespondentInfo{}

	query := db.
		Table("respondents").
		Joins("LEFT OUTER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Order("respondents.submitted_at DESC").
		Where("user_traqid = ? AND respondents.deleted_at IS NULL AND questionnaires.deleted_at IS NULL", userID)

	if len(questionnaireIDs) > 1 {
		// 空配列か1要素の取得にしか用いない
		return nil, errors.New("illegal function usage")
	}
	if len(questionnaireIDs) != 0 {
		questionnaireID := questionnaireIDs[0]
		query = query.Where("questionnaire_id = ?", questionnaireID)
	}

	err = query.
		Select("respondents.questionnaire_id, respondents.response_id, respondents.modified_at, respondents.submitted_at, questionnaires.title, questionnaires.description, questionnaires.res_time_limit").
		Find(&respondentInfos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get my responses: %w", err)
	}

	return respondentInfos, nil
}

// GetRespondentDetail 回答のIDから回答の詳細情報を取得
func (*Respondent) GetRespondentDetail(ctx context.Context, responseID int) (RespondentDetail, error) {
	db, err := getTx(ctx)
	if err != nil {
		return RespondentDetail{}, fmt.Errorf("failed to get tx: %w", err)
	}

	respondent := Respondents{}

	err = db.
		Session(&gorm.Session{}).
		Where("respondents.response_id = ?", responseID).
		Select("QuestionnaireID", "UserTraqid", "ModifiedAt", "SubmittedAt").
		Take(&respondent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return RespondentDetail{}, ErrRecordNotFound
	}
	if err != nil {
		return RespondentDetail{}, fmt.Errorf("failed to get respondent: %w", err)
	}

	questions := []Questions{}
	err = db.
		Where("questionnaire_id = ?", respondent.QuestionnaireID).
		Preload("Responses", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("QuestionID", "Body").
				Where("response_id = ?", responseID)
		}).
		Select("ID", "Type").
		Find(&questions).Error
	if err != nil {
		return RespondentDetail{}, fmt.Errorf("failed to get questions: %w", err)
	}

	respondentDetail := RespondentDetail{
		ResponseID:      responseID,
		TraqID:          respondent.UserTraqid,
		QuestionnaireID: respondent.QuestionnaireID,
		ModifiedAt:      respondent.ModifiedAt,
		SubmittedAt:     respondent.SubmittedAt,
	}

	for _, question := range questions {
		responseBody := ResponseBody{
			QuestionID:   question.ID,
			QuestionType: question.Type,
		}

		switch question.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, response := range question.Responses {
				responseBody.OptionResponse = append(responseBody.OptionResponse, response.Body.String)
			}
		default:
			if len(question.Responses) == 0 {
				responseBody.Body = null.NewString("", false)
			} else {
				responseBody.Body = question.Responses[0].Body
			}
		}

		respondentDetail.Responses = append(respondentDetail.Responses, responseBody)
	}

	return respondentDetail, nil
}

// GetRespondentDetails アンケートの回答の詳細情報一覧の取得
func (*Respondent) GetRespondentDetails(ctx context.Context, questionnaireID int, sort string, onlyMyResponse bool, userID string, isDraft *bool) ([]RespondentDetail, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	respondents := []Respondents{}

	// Note: respondents.submitted_at IS NOT NULLで一時保存の回答を除外している
	query := db.
		Session(&gorm.Session{}).
		Where("respondents.questionnaire_id = ? AND respondents.submitted_at IS NOT NULL", questionnaireID).
		Select("ResponseID", "UserTraqid", "ModifiedAt", "SubmittedAt")
	if onlyMyResponse {
		query = query.Where("user_traqid = ?", userID)
	}
	if isDraft != nil {
		if *isDraft {
			query = query.Where("submitted_at IS NULL")
		} else {
			query = query.Where("submitted_at IS NOT NULL")
		}
	}

	query, sortNum, err := setRespondentsOrder(query, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to set order: %w", err)
	}

	err = query.
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get respondents: %w", err)
	}

	if len(respondents) == 0 {
		return []RespondentDetail{}, nil
	}

	responseIDs := make([]int, 0, len(respondents))
	for _, respondent := range respondents {
		responseIDs = append(responseIDs, respondent.ResponseID)
	}

	isAnonymous, err := NewQuestionnaire().GetResponseIsAnonymousByQuestionnaireID(ctx, questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("failed to get response is anonymous by questionnaire id: %w", err)
	}

	respondentDetails := make([]RespondentDetail, 0, len(respondents))
	respondentDetailMap := make(map[int]*RespondentDetail, len(respondents))
	for i, respondent := range respondents {
		r := RespondentDetail{
			ResponseID:      respondent.ResponseID,
			QuestionnaireID: questionnaireID,
			SubmittedAt:     respondent.SubmittedAt,
			ModifiedAt:      respondent.ModifiedAt,
		}

		if !isAnonymous {
			r.TraqID = respondent.UserTraqid
		} else {
			r.TraqID = ""
		}

		respondentDetails = append(respondentDetails, r)

		respondentDetailMap[respondent.ResponseID] = &respondentDetails[i]
	}

	questions := []Questions{}
	err = db.
		Preload("Responses", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("ResponseID", "QuestionID", "Body").
				Where("response_id IN (?)", responseIDs)
		}).
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Select("ID", "Type").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	for _, question := range questions {
		responseBodyMap := make(map[int][]string, len(respondents))
		for _, response := range question.Responses {
			if response.Body.Valid {
				responseBodyMap[response.ResponseID] = append(responseBodyMap[response.ResponseID], response.Body.String)
			}
		}

		for i := range respondentDetails {
			responseBodies := responseBodyMap[respondentDetails[i].ResponseID]
			responseBody := ResponseBody{
				QuestionID:   question.ID,
				QuestionType: question.Type,
			}

			switch responseBody.QuestionType {
			case "MultipleChoice", "Checkbox", "Dropdown":
				if responseBodies == nil {
					responseBody.OptionResponse = []string{}
				} else {
					responseBody.OptionResponse = responseBodies
				}
			default:
				if len(responseBodies) == 0 {
					responseBody.Body = null.NewString("", false)
				} else {
					responseBody.Body = null.NewString(responseBodies[0], true)
				}
			}

			respondentDetails[i].Responses = append(respondentDetails[i].Responses, responseBody)
		}
	}

	respondentDetails, err = sortRespondentDetail(sortNum, len(questions), respondentDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to sort RespondentDetails: %w", err)
	}

	return respondentDetails, nil
}

type myResponseGroupRow struct {
	QuestionnaireID int       `gorm:"column:questionnaire_id"`
	Title           string    `gorm:"column:title"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	ModifiedAt      time.Time `gorm:"column:modified_at"`
	ResTimeLimit    null.Time `gorm:"column:res_time_limit"`
	IsAnonymous     bool      `gorm:"column:is_anonymous"`
	IsTargetingMe   bool      `gorm:"column:is_targeting_me"`
	FirstResponseID int       `gorm:"column:first_response_id"`
}

func buildMyResponseBaseQuery(db *gorm.DB, userID string, questionnaireIDs []int, isDraft *bool) *gorm.DB {
	query := db.
		Table("respondents").
		Joins("INNER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Where("respondents.deleted_at IS NULL AND questionnaires.deleted_at IS NULL AND respondents.user_traqid = ?", userID)

	if questionnaireIDs != nil {
		query = query.Where("respondents.questionnaire_id IN (?)", questionnaireIDs)
	}

	if isDraft != nil {
		if *isDraft {
			query = query.Where("respondents.submitted_at IS NULL")
		} else {
			query = query.Where("respondents.submitted_at IS NOT NULL")
		}
	}

	return query
}

func setMyResponseGroupOrder(query *gorm.DB, sort string) (*gorm.DB, error) {
	switch sort {
	case "":
		return query.Order("first_response_id"), nil
	case "submitted_at":
		return query.
			Order("MIN(respondents.submitted_at)").
			Order("first_response_id"), nil
	case "-submitted_at":
		return query.
			Order("MAX(respondents.submitted_at) DESC").
			Order("first_response_id"), nil
	case "modified_at":
		return query.
			Order("MIN(respondents.modified_at)").
			Order("first_response_id"), nil
	case "-modified_at":
		return query.
			Order("MAX(respondents.modified_at) DESC").
			Order("first_response_id"), nil
	case "traqid", "-traqid":
		return query.Order("first_response_id"), nil
	default:
		return nil, fmt.Errorf("failed to convert sort param to group order: %w", ErrInvalidSortParam)
	}
}

// GetMyResponseGroups 自分の回答をアンケートごとにまとめて取得
func (*Respondent) GetMyResponseGroups(ctx context.Context, sort string, userID string, questionnaireIDs []int, isDraft *bool, pageNum int) ([]MyResponseGroup, int, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transaction: %w", err)
	}

	baseQuery := buildMyResponseBaseQuery(db, userID, questionnaireIDs, isDraft)

	var count int64
	err = baseQuery.
		Session(&gorm.Session{}).
		Distinct("respondents.questionnaire_id").
		Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count my response questionnaires: %w", err)
	}

	if count == 0 {
		return []MyResponseGroup{}, 0, nil
	}

	pageMax := (int(count) + 19) / 20
	if pageNum > pageMax {
		return nil, 0, ErrTooLargePageNum
	}

	groupRows := []myResponseGroupRow{}
	groupQuery := buildMyResponseBaseQuery(db, userID, questionnaireIDs, isDraft).
		Select(
			"respondents.questionnaire_id, questionnaires.title, questionnaires.created_at, questionnaires.modified_at, questionnaires.res_time_limit, questionnaires.is_anonymous, "+
				"EXISTS(SELECT 1 FROM targets WHERE targets.questionnaire_id = questionnaires.id AND targets.user_traqid = ?) AS is_targeting_me, "+
				"MIN(respondents.response_id) AS first_response_id",
			userID,
		).
		Group("respondents.questionnaire_id, questionnaires.id, questionnaires.title, questionnaires.created_at, questionnaires.modified_at, questionnaires.res_time_limit, questionnaires.is_anonymous")
	groupQuery, err = setMyResponseGroupOrder(groupQuery, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to set my response group order: %w", err)
	}

	err = groupQuery.
		Limit(20).
		Offset((pageNum - 1) * 20).
		Find(&groupRows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get my response groups: %w", err)
	}

	if len(groupRows) == 0 {
		return []MyResponseGroup{}, pageMax, nil
	}

	groups := make([]MyResponseGroup, 0, len(groupRows))
	groupIndexByQuestionnaireID := make(map[int]int, len(groupRows))
	pageQuestionnaireIDs := make([]int, 0, len(groupRows))
	for i, row := range groupRows {
		groups = append(groups, MyResponseGroup{
			QuestionnaireInfo: MyResponseQuestionnaireInfo{
				QuestionnaireID:     row.QuestionnaireID,
				Title:               row.Title,
				CreatedAt:           row.CreatedAt,
				ModifiedAt:          row.ModifiedAt,
				ResponseDueDateTime: row.ResTimeLimit,
				IsAnonymous:         row.IsAnonymous,
				IsTargetingMe:       row.IsTargetingMe,
			},
			Responses: []RespondentDetail{},
		})
		groupIndexByQuestionnaireID[row.QuestionnaireID] = i
		pageQuestionnaireIDs = append(pageQuestionnaireIDs, row.QuestionnaireID)
	}

	respondents := []Respondents{}
	respondentQuery := db.
		Session(&gorm.Session{}).
		Where("respondents.deleted_at IS NULL AND respondents.user_traqid = ? AND respondents.questionnaire_id IN (?)", userID, pageQuestionnaireIDs).
		Select("ResponseID", "QuestionnaireID", "UserTraqid", "ModifiedAt", "SubmittedAt")
	if isDraft != nil {
		if *isDraft {
			respondentQuery = respondentQuery.Where("submitted_at IS NULL")
		} else {
			respondentQuery = respondentQuery.Where("submitted_at IS NOT NULL")
		}
	}
	respondentQuery, _, err = setRespondentsOrder(respondentQuery, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to set respondents order: %w", err)
	}

	err = respondentQuery.Find(&respondents).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get respondents for my response groups: %w", err)
	}

	if len(respondents) == 0 {
		return groups, pageMax, nil
	}

	responseIDs := make([]int, 0, len(respondents))
	respondentDetailMap := make(map[int]*RespondentDetail, len(respondents))
	responseIDsByQuestionnaireID := make(map[int][]int, len(groupRows))
	for _, respondent := range respondents {
		groupIdx := groupIndexByQuestionnaireID[respondent.QuestionnaireID]
		groups[groupIdx].Responses = append(groups[groupIdx].Responses, RespondentDetail{
			ResponseID:      respondent.ResponseID,
			TraqID:          respondent.UserTraqid,
			QuestionnaireID: respondent.QuestionnaireID,
			ModifiedAt:      respondent.ModifiedAt,
			SubmittedAt:     respondent.SubmittedAt,
			Responses:       []ResponseBody{},
		})
		lastIdx := len(groups[groupIdx].Responses) - 1
		respondentDetailMap[respondent.ResponseID] = &groups[groupIdx].Responses[lastIdx]
		responseIDs = append(responseIDs, respondent.ResponseID)
		responseIDsByQuestionnaireID[respondent.QuestionnaireID] = append(responseIDsByQuestionnaireID[respondent.QuestionnaireID], respondent.ResponseID)
	}

	questions := []Questions{}
	err = db.
		Preload("Responses", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("ResponseID", "QuestionID", "Body").
				Where("response_id IN (?)", responseIDs)
		}).
		Where("questionnaire_id IN (?)", pageQuestionnaireIDs).
		Order("questionnaire_id").
		Order("question_num").
		Select("ID", "QuestionnaireID", "QuestionNum", "Type").
		Find(&questions).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get questions for my response groups: %w", err)
	}

	for _, question := range questions {
		responseBodyMap := make(map[int][]string)
		for _, response := range question.Responses {
			if response.Body.Valid {
				responseBodyMap[response.ResponseID] = append(responseBodyMap[response.ResponseID], response.Body.String)
			}
		}

		for _, responseID := range responseIDsByQuestionnaireID[question.QuestionnaireID] {
			respondentDetail := respondentDetailMap[responseID]
			if respondentDetail == nil {
				continue
			}

			responseBodies := responseBodyMap[responseID]
			responseBody := ResponseBody{
				QuestionID:   question.ID,
				QuestionType: question.Type,
			}

			switch question.Type {
			case "MultipleChoice", "Checkbox", "Dropdown":
				if responseBodies == nil {
					responseBody.OptionResponse = []string{}
				} else {
					responseBody.OptionResponse = responseBodies
				}
			default:
				if len(responseBodies) == 0 {
					responseBody.Body = null.NewString("", false)
				} else {
					responseBody.Body = null.NewString(responseBodies[0], true)
				}
			}

			respondentDetail.Responses = append(respondentDetail.Responses, responseBody)
		}
	}

	return groups, pageMax, nil
}

// GetRespondentsUserIDs 回答者のユーザーID取得
func (*Respondent) GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]Respondents, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %w", err)
	}

	respondents := []Respondents{}

	err = db.
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Select("questionnaire_id, user_traqid").
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get respondents:%w", err)
	}

	return respondents, nil
}

// GetMyResponses 自分のすべての回答を取得
func (*Respondent) GetMyResponseIDs(ctx context.Context, sort string, userID string, questionnaireIDs []int, isDraft *bool) ([]int, error) {
	db, err := getTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	responsesID := []int{}
	query := db.Model(&Respondents{}).
		Where("deleted_at IS NULL AND user_traqid = ?", userID).
		Select("response_id")

	if questionnaireIDs != nil {
		query = query.Where("questionnaire_id IN (?)", questionnaireIDs)
	}

	if isDraft != nil {
		if *isDraft {
			query = query.Where("submitted_at IS NULL")
		} else {
			query = query.Where("submitted_at IS NOT NULL")
		}
	}

	query, _, err = setRespondentsOrder(query, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to set respondents order: %w", err)
	}

	err = query.Select("respondents.response_id").Find(&responsesID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responsesID: %w", err)
	}

	return responsesID, nil
}

// CheckRespondent 回答者かどうかの確認
func (*Respondent) CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error) {
	db, err := getTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get tx: %w", err)
	}

	err = db.
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		First(&Respondents{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get response: %w", err)
	}

	return true, nil
}

func setRespondentsOrder(query *gorm.DB, sort string) (*gorm.DB, int, error) {
	var sortNum int
	switch sort {
	case "traqid":
		query = query.Order("user_traqid")
	case "-traqid":
		query = query.Order("user_traqid DESC")
	case "submitted_at":
		query = query.Order("submitted_at")
	case "-submitted_at":
		query = query.Order("submitted_at DESC")
	case "modified_at":
		query = query.Order("modified_at")
	case "-modified_at":
		query = query.Order("modified_at DESC")
	case "":
	default:
		var err error
		sortNum, err = strconv.Atoi(sort)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert sort param to int: %w", err)
		}
	}

	query = query.Order("response_id")

	return query, sortNum, nil
}

func sortRespondentDetail(sortNum int, questionNum int, respondentDetails []RespondentDetail) ([]RespondentDetail, error) {
	if sortNum == 0 {
		return respondentDetails, nil
	}
	sortNumAbs := int(math.Abs(float64(sortNum)))
	if sortNumAbs > questionNum {
		return nil, fmt.Errorf("sort param is too large: %d", sortNum)
	}

	sort.Slice(respondentDetails, func(i, j int) bool {
		bodyI := respondentDetails[i].Responses[sortNumAbs-1]
		bodyJ := respondentDetails[j].Responses[sortNumAbs-1]
		if bodyI.QuestionType == "Number" {
			numi, err := strconv.ParseFloat(bodyI.Body.String, 64)
			if err != nil {
				return true
			}
			numj, err := strconv.ParseFloat(bodyJ.Body.String, 64)
			if err != nil {
				return true
			}
			if sortNum < 0 {
				return numi > numj
			}
			return numi < numj
		}
		if bodyI.QuestionType == "MultipleChoice" {
			choiceI := ""
			if len(bodyI.OptionResponse) > 0 {
				choiceI = bodyI.OptionResponse[0]
			}
			choiceJ := ""
			if len(bodyJ.OptionResponse) > 0 {
				choiceJ = bodyJ.OptionResponse[0]
			}
			if sortNum < 0 {
				return choiceI > choiceJ
			}
			return choiceI < choiceJ
		}
		if bodyI.QuestionType == "Checkbox" {
			selectionsI := strings.Join(bodyI.OptionResponse, ", ")
			selectionsJ := strings.Join(bodyJ.OptionResponse, ", ")
			if sortNum < 0 {
				return selectionsI > selectionsJ
			}
			return selectionsI < selectionsJ
		}
		if sortNum < 0 {
			return bodyI.Body.String > bodyJ.Body.String
		}
		return bodyI.Body.String < bodyJ.Body.String
	})

	return respondentDetails, nil
}
