package model

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

//Respondents respondentsテーブルの構造体
type Respondents struct {
	ResponseID      int       `json:"responseID" gorm:"type:int(11) AUTO_INCREMENT NOT NULL PRIMARY KEY;"`
	QuestionnaireID int       `json:"questionnaireID" gorm:"type:int(11) NOT NULL;"`
	UserTraqid      string    `json:"user_traq_id,omitempty" gorm:"type:char(30) NOT NULL;"`
	ModifiedAt      time.Time `json:"modified_at,omitempty" gorm:"type:timestamp NOT NULL;default:CURRENT_TIMESTAMP;"`
	SubmittedAt     null.Time `json:"submitted_at,omitempty" gorm:"type:timestamp NULL;default:CURRENT_TIMESTAMP;"`
	DeletedAt       null.Time `json:"deleted_at,omitempty" gorm:"type:timestamp NULL;default:NULL;"`
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (*Respondents) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set ModifiedAt: %w", err)
	}

	err = scope.SetColumn("SubmittedAt", time.Now())
	if err != nil {
		return fmt.Errorf("failed to set SubmitedAt: %w", err)
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
	Title        string `json:"questionnaire_title"`
	ResTimeLimit string `json:"res_time_limit"`
	Respondents
}

// RespondentDetail 回答の詳細情報の構造体
type RespondentDetail struct {
	ResponseID      int            `json:"responseID,omitempty"`
	TraqID          string         `json:"traqID,omitempty"`
	QuestionnaireID int            `json:"questionnaireID,omitempty"`
	SubmittedAt     time.Time      `json:"submitted_at,omitempty"`
	ModifiedAt      time.Time      `json:"modified_at,omitempty"`
	Responses       []ResponseBody `json:"body"`
}

//InsertRespondent 回答の追加
func InsertRespondent(c echo.Context, questionnaireID int, submitedAt null.Time) (int, error) {
	userID := GetUserID(c)

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
			c.Logger().Error(fmt.Errorf("failed to insert a respondent record: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		err = tx.Select("response_id").Order("response_id DESC").Last(&respondent).Error
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to get the last respondent record: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return nil
	})
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed in transaction: %w", err))
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return respondent.ResponseID, nil
}

// DeleteRespondent 回答の削除
func DeleteRespondent(c echo.Context, responseID int) error {
	userID := GetUserID(c)

	err := db.Exec("UPDATE `respondents` INNER JOIN administrators ON administrators.questionnaire_id = respondents.questionnaire_id SET `respondents`.`deleted_at` = ? WHERE (respondents.response_id = ? AND (administrators.user_traqid = ? OR respondents.user_traqid = ?))", time.Now(), responseID, userID, userID).Error
	if gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = db.
		Where("response_id = ?", responseID).
		Delete(&Response{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// GetRespondentInfos ユーザーの回答とその周辺情報一覧の取得
func GetRespondentInfos(c echo.Context, userID string, questionnaireIDs ...int) ([]RespondentInfo, error) {
	respondentInfos := []RespondentInfo{}

	query := db.
		Table("respondents").
		Joins("LEFT OUTER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Order("respondents.submitted_at DESC").
		Where("user_traqid = ? AND respondents.deleted_at IS NULL", userID)

	if len(questionnaireIDs) != 0 {
		questionnaireID := questionnaireIDs[0]
		query = query.Where("questionnaire_id = ?", questionnaireID)
	}

	rows, err := query.
		Select("respondents.questionnaire_id, respondents.response_id, respondents.modified_at, respondents.submitted_at, questionnaires.title, " +
			"questionnaires.res_time_limit").
		Rows()
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get my responses: %w", err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	for rows.Next() {
		respondentInfo := RespondentInfo{
			Respondents: Respondents{},
		}

		err := db.ScanRows(rows, &respondentInfo)
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to scan responses: %w", err))
			return nil, echo.NewHTTPError(http.StatusInternalServerError)
		}

		respondentInfos = append(respondentInfos, respondentInfo)
	}

	return respondentInfos, nil
}

// GetRespondentDetail 回答のIDから回答の詳細情報を取得
func GetRespondentDetail(c echo.Context, responseID int) (RespondentDetail, error) {
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
			if !res.Respondents.SubmittedAt.Valid {
				return RespondentDetail{}, fmt.Errorf("unexpected null submited_at(response_id: %d)", res.ResponseID)
			}
			respondentDetail.SubmittedAt = res.Respondents.SubmittedAt.Time
			respondentDetail.ModifiedAt = res.Respondents.ModifiedAt
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
func GetRespondentDetails(c echo.Context, questionnaireID int, sort string) ([]RespondentDetail, error) {
	query := db.
		Table("respondents").
		Joins("LEFT OUTER JOIN question ON respondents.questionnaire_id = question.questionnaire_id").
		Joins("LEFT OUTER JOIN response ON respondents.response_id = response.response_id AND question.id = response.question_id")
	query, sortNum, err := setRespondentsOrder(query, sort)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to set order: %w", err))
	}

	rows, err := query.
		Where("respondents.questionnaire_id = ? AND respondents.deleted_at IS NULL", questionnaireID).
		Select("respondents.response_id, respondents.user_traqid, respondents.modified_at, respondents.submitted_at, question.id, question.type, response.body").
		Rows()
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get respondents: %w", err))
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
			return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to scan response detail: %w", err))
		}

		log.Printf("%+v", res)
		if _, ok := responseBodyMap[res.ResponseID]; !ok {
			if !res.Respondents.SubmittedAt.Valid {
				return nil, fmt.Errorf("unexpected null submited_at(response_id: %d)", res.ResponseID)
			}
			respondentDetails = append(respondentDetails, RespondentDetail{
				ResponseID:      res.Respondents.ResponseID,
				TraqID:          res.UserTraqid,
				QuestionnaireID: res.Respondents.QuestionnaireID,
				SubmittedAt:     res.Respondents.SubmittedAt.Time,
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
			body, ok := bodyMap[responseBody.QuestionID]
			switch responseBody.QuestionType {
			case "MultipleChoice", "Checkbox", "Dropdown":
				if !ok {
					return nil, errors.New("unexpected no response")
				}
				responseBody.OptionResponse = body
			default:
				if !ok || len(body) == 0 {
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
		return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to sort RespondentDetails: %w", err))
	}

	return respondentDetails, nil
}

// CheckRespondent 回答者かどうかの確認
func CheckRespondent(userID string, questionnaireID int) (bool, error) {
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
func CheckRespondentByResponseID(userID string, responseID int) (bool, error) {
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
			numi, err := strconv.Atoi(bodyI.Body.String)
			if err != nil {
				return true
			}
			numj, err := strconv.Atoi(bodyJ.Body.String)
			if err != nil {
				return true
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
