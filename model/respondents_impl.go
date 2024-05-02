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
func (r *Respondents) BeforeCreate(tx *gorm.DB) error {
	r.ModifiedAt = time.Now()

	return nil
}

// BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (r *Respondents) BeforeUpdate(tx *gorm.DB) error {
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

	if !questionnaire.IsDuplicateAnswerAllowed && submittedAt.Valid {
		// delete old answers
		err = db.
			Where("questionnaire_id = ? AND user_traqid = ?", questionnaireID, userID).
			Delete(&Respondents{}).Error
		// 既存の回答がなかった場合はそのまま進む
		if errors.Is(err, ErrNoRecordDeleted) {
			err = nil
		}
		if err != nil {
			return 0, fmt.Errorf("failed to delete old answers: %w", err)
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

	if len(questionnaireIDs) != 0 {
		questionnaireID := questionnaireIDs[0]
		query = query.Where("questionnaire_id = ?", questionnaireID)
	} else if len(questionnaireIDs) > 1 {
		// 空配列か1要素の取得にしか用いない
		return nil, errors.New("illegal function usage")
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
func (*Respondent) GetRespondentDetails(ctx context.Context, questionnaireID int, sort string) ([]RespondentDetail, error) {
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

	respondentDetails := make([]RespondentDetail, 0, len(respondents))
	respondentDetailMap := make(map[int]*RespondentDetail, len(respondents))
	for i, respondent := range respondents {
		respondentDetails = append(respondentDetails, RespondentDetail{
			ResponseID:      respondent.ResponseID,
			TraqID:          respondent.UserTraqid,
			QuestionnaireID: questionnaireID,
			SubmittedAt:     respondent.SubmittedAt,
			ModifiedAt:      respondent.ModifiedAt,
		})

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
