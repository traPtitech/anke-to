package model

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

// Respondent RespondentRepositoryの実装
type Respondent struct{}

// NewRespondent Respondentのコンストラクター
func NewRespondent() *Respondent {
	return new(Respondent)
}

//Respondents respondentsテーブルの構造体
type Respondents struct {
	ResponseID      int       `json:"responseID" gorm:"type:int(11) AUTO_INCREMENT NOT NULL PRIMARY KEY;"`
	QuestionnaireID int       `json:"questionnaireID" gorm:"type:int(11) NOT NULL;"`
	UserTraqid      string    `json:"user_traq_id,omitempty" gorm:"type:char(30) NOT NULL;"`
	ModifiedAt      time.Time `json:"modified_at,omitempty" gorm:"type:timestamp NOT NULL;default:CURRENT_TIMESTAMP;"`
	SubmittedAt     null.Time `json:"submitted_at,omitempty" gorm:"type:timestamp NULL;default:NULL;"`
	DeletedAt       null.Time `json:"deleted_at,omitempty" gorm:"type:timestamp NULL;default:NULL;"`
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (*Respondents) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set ModifiedAt: %w", err)
	}

	return nil
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (*Respondents) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set ModifiedAt: %w", err)
	}

	return nil
}

// RespondentInfo 回答とその周辺情報の構造体
type RespondentInfo struct {
	Title        string    `json:"questionnaire_title"`
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

//InsertRespondent 回答の追加
func (*Respondent) InsertRespondent(userID string, questionnaireID int, submitedAt null.Time) (int, error) {
	var respondent Respondents
	if submitedAt.Valid {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
			SubmittedAt:     submitedAt,
		}
	} else {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&respondent).Error
		if err != nil {
			return fmt.Errorf("failed to insert a respondent record: %w", err)
		}

		err = tx.Select("response_id").Order("response_id DESC").Last(&respondent).Error
		if err != nil {
			return fmt.Errorf("failed to get the last respondent record: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed in transaction: %w", err)
	}

	return respondent.ResponseID, nil
}

// UpdateSubmittedAt 投稿日時更新
func (*Respondent) UpdateSubmittedAt(responseID int) error {
	err := db.
		Model(&Respondents{}).
		Where("response_id = ?", responseID).
		Update("submitted_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("failed to update response's submitted_at: %w", err)
	}

	return nil
}

// DeleteRespondent 回答の削除
func (*Respondent) DeleteRespondent(userID string, responseID int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec("UPDATE `respondents` INNER JOIN administrators ON administrators.questionnaire_id = respondents.questionnaire_id SET `respondents`.`deleted_at` = ? WHERE (respondents.response_id = ? AND (administrators.user_traqid = ? OR respondents.user_traqid = ?))", time.Now(), responseID, userID, userID)
		err := result.Error
		if err != nil {
			return fmt.Errorf("failed to delete respondents: %w", err)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("failed to delete respondents : %w", ErrNoRecordDeleted)
		}

		err = tx.
			Where("response_id = ?", responseID).
			Delete(&Responses{}).Error
		if err != nil {
			return fmt.Errorf("failed to delete response: %w", err)
		}
		return nil
	})
}

// GetRespondentInfos ユーザーの回答とその周辺情報一覧の取得
func (*Respondent) GetRespondentInfos(userID string, questionnaireIDs ...int) ([]RespondentInfo, error) {
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
		return nil, fmt.Errorf("ilegal function usase")
	}

	rows, err := query.
		Select("respondents.questionnaire_id, respondents.response_id, respondents.modified_at, respondents.submitted_at, questionnaires.title, " +
			"questionnaires.res_time_limit").
		Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get my responses: %w", err)
	}

	for rows.Next() {
		respondentInfo := RespondentInfo{
			Respondents: Respondents{},
		}

		err := db.ScanRows(rows, &respondentInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to scan responses: %w", err)
		}

		respondentInfos = append(respondentInfos, respondentInfo)
	}

	return respondentInfos, nil
}

// GetRespondentDetail 回答のIDから回答の詳細情報を取得
func (*Respondent) GetRespondentDetail(responseID int) (RespondentDetail, error) {
	rows, err := db.
		Table("respondents").
		Joins("LEFT OUTER JOIN question ON respondents.questionnaire_id = question.questionnaire_id").
		Joins("LEFT OUTER JOIN response ON respondents.response_id = response.response_id AND question.id = response.question_id AND response.deleted_at IS NULL").
		Where("respondents.response_id = ? AND respondents.deleted_at IS NULL", responseID).
		Select("respondents.questionnaire_id, respondents.modified_at, respondents.submitted_at, question.id, question.type, response.body").
		Rows()
	if err != nil {
		return RespondentDetail{}, fmt.Errorf("failed to get respondents: %w", err)
	}

	isNoRows := true
	isRespondentSetted := false
	respondentDetail := RespondentDetail{}
	responseBodyMap := map[int][]string{}
	for rows.Next() {
		isNoRows = false
		res := struct {
			Respondents  `gorm:"embedded"`
			ResponseBody `gorm:"embedded"`
		}{}
		err := db.ScanRows(rows, &res)
		if err != nil {
			return RespondentDetail{}, fmt.Errorf("failed to scan response detail: %w", err)
		}
		if !isRespondentSetted {
			respondentDetail.QuestionnaireID = res.Respondents.QuestionnaireID
			respondentDetail.SubmittedAt = res.Respondents.SubmittedAt
			respondentDetail.ModifiedAt = res.Respondents.ModifiedAt
			isRespondentSetted = true
		}

		respondentDetail.Responses = append(respondentDetail.Responses, ResponseBody{
			QuestionID:   res.ResponseBody.QuestionID,
			QuestionType: res.ResponseBody.QuestionType,
		})

		if res.ResponseBody.Body.Valid {
			responseBodyMap[res.ResponseBody.QuestionID] = append(responseBodyMap[res.ResponseBody.QuestionID], res.ResponseBody.Body.String)
		}
	}
	if isNoRows {
		return RespondentDetail{}, fmt.Errorf("failed to get respondents: %w", gorm.ErrRecordNotFound)
	}

	for i := range respondentDetail.Responses {
		response := &respondentDetail.Responses[i]
		responseBody := responseBodyMap[response.QuestionID]
		switch response.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			response.OptionResponse = responseBody
		default:
			if len(responseBody) == 0 {
				response.Body = null.NewString("", false)
			} else {
				response.Body = null.NewString(responseBody[0], true)
			}
		}
	}

	return respondentDetail, nil
}

// GetRespondentDetails アンケートの回答の詳細情報一覧の取得
func (*Respondent) GetRespondentDetails(questionnaireID int, sort string) ([]RespondentDetail, error) {
	query := db.
		Table("respondents").
		Joins("LEFT OUTER JOIN question ON respondents.questionnaire_id = question.questionnaire_id").
		Joins("LEFT OUTER JOIN response ON respondents.response_id = response.response_id AND question.id = response.question_id")
	query, sortNum, err := setRespondentsOrder(query, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to set order: %w", err)
	}

	rows, err := query.
		Where("respondents.questionnaire_id = ? AND respondents.deleted_at IS NULL AND respondents.submitted_at IS NOT NULL AND question.deleted_at IS NULL AND response.deleted_at IS NULL", questionnaireID).
		Select("respondents.response_id, respondents.user_traqid, respondents.modified_at, respondents.submitted_at, question.id, question.type, response.body").
		Order("respondents.response_id, question.question_num").
		Rows()
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return []RespondentDetail{}, nil
		}
	}

	respondentDetails := []RespondentDetail{}
	responseBodyMap := map[int][]ResponseBody{}
	for rows.Next() {
		res := struct {
			Respondents  `gorm:"embedded"`
			ResponseBody `gorm:"embedded"`
		}{}
		err := db.ScanRows(rows, &res)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response detail: %w", err)
		}

		if _, ok := responseBodyMap[res.ResponseID]; !ok {
			respondentDetails = append(respondentDetails, RespondentDetail{
				ResponseID:      res.Respondents.ResponseID,
				TraqID:          res.UserTraqid,
				QuestionnaireID: res.Respondents.QuestionnaireID,
				SubmittedAt:     res.Respondents.SubmittedAt,
				ModifiedAt:      res.ModifiedAt,
			})
		}

		responseBodyMap[res.ResponseID] = append(responseBodyMap[res.ResponseID], res.ResponseBody)
	}

	for i := range respondentDetails {
		responseDetail := &respondentDetails[i]
		responseBodies := responseBodyMap[responseDetail.ResponseID]

		responseBodyList := []ResponseBody{}
		bodyMap := map[int][]string{}
		for _, v := range responseBodies {
			if _, ok := bodyMap[v.QuestionID]; !ok {
				responseBodyList = append(responseBodyList, ResponseBody{
					QuestionID:   v.QuestionID,
					QuestionType: v.QuestionType,
				})
			}

			if v.Body.Valid {
				bodyMap[v.QuestionID] = append(bodyMap[v.QuestionID], v.Body.String)
			}
		}

		for i := range responseBodyList {
			responseBody := &responseBodyList[i]
			body := bodyMap[responseBody.QuestionID]
			switch responseBody.QuestionType {
			case "MultipleChoice", "Checkbox", "Dropdown":
				if body == nil {
					responseBody.OptionResponse = []string{}
				} else {
					responseBody.OptionResponse = body
				}
			default:
				if len(body) == 0 {
					responseBody.Body = null.NewString("", false)
				} else {
					responseBody.Body = null.NewString(body[0], true)
				}
			}
		}
		responseDetail.Responses = responseBodyList
	}

	respondentDetails, err = sortRespondentDetail(sortNum, respondentDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to sort RespondentDetails: %w", err)
	}

	return respondentDetails, nil
}

// GetRespondentsUserIDs 回答者のユーザーID取得
func (*Respondent) GetRespondentsUserIDs(questionnaireIDs []int) ([]Respondents, error) {
	respondents := []Respondents{}
	err := db.
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Select("questionnaire_id, user_traqid").
		Find(&respondents).Error
	if err != nil {
		return []Respondents{}, nil
	}

	return respondents, nil
}

// CheckRespondent 回答者かどうかの確認
func (*Respondent) CheckRespondent(userID string, questionnaireID int) (bool, error) {
	err := db.
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		First(&Respondents{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get response: %w", err)
	}

	return true, nil
}

// CheckRespondentByResponseID 回答者かどうかの確認
func (*Respondent) CheckRespondentByResponseID(userID string, responseID int) (bool, error) {
	err := db.
		Where("user_traqid = ? AND response_id = ?", userID, responseID).
		First(&Respondents{}).Error
	if gorm.IsRecordNotFoundError(err) {
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
		query = query.Order("respondents.user_traqid")
	case "-traqid":
		query = query.Order("respondents.user_traqid DESC")
	case "submitted_at":
		query = query.Order("respondents.submitted_at")
	case "-submitted_at":
		query = query.Order("respondents.submitted_at DESC")
	case "":
	default:
		var err error
		sortNum, err = strconv.Atoi(sort)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert sort param to int: %w", err)
		}
	}

	return query, sortNum, nil
}

func sortRespondentDetail(sortNum int, respondentDetails []RespondentDetail) ([]RespondentDetail, error) {
	if sortNum == 0 {
		return respondentDetails, nil
	}
	sortNumAbs := int(math.Abs(float64(sortNum)))
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
		if sortNum < 0 {
			return bodyI.Body.String > bodyJ.Body.String
		}
		return bodyI.Body.String < bodyJ.Body.String
	})

	return respondentDetails, nil
}
