package model

import (
	"errors"
	"fmt"
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
	ResponseID      int       `json:"responseID" gorm:"type:int(11);NOT NULL;PRIMARY_KEY;AUTO_INCREMENT"`
	QuestionnaireID int       `json:"questionnaireID" gorm:"type:int(11);NOT NULL;"`
	UserTraqid      string    `json:"user_traq_id,omitempty" gorm:"type:char(30);NOT NULL;"`
	ModifiedAt      time.Time `json:"modified_at,omitempty" gorm:"type:timestamp;NOT NULL;DEFAULT CURRENT_TIMESTAMP;"`
	SubmittedAt     null.Time `json:"submitted_at,omitempty" gorm:"type:timestamp;"`
	DeletedAt       null.Time `json:"deleted_at,omitempty" gorm:"type:timestamp;"`
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (*Respondents) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedAt", time.Now())
	scope.SetColumn("SubmittedAt", time.Now())

	return nil
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (*Respondents) BeforeUpdate(scope *gorm.Scope) (err error) {
	scope.SetColumn("ModifiedAt", time.Now())

	return nil
}

type RespondentInfo struct {
	Title        string `json:"questionnaire_title"`
	ResTimeLimit string `json:"res_time_limit"`
	Respondents
}

type RespondentDetail struct {
	Respondents `gorm:"embedded"`
	Responses []ResponseBody `json:"body"`
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
		return nil, errors.New("invalid sort num")
	}
	sortNumAbs := int(math.Abs(float64(sortNum)))
	sort.Slice(respondentDetails, func(i, j int) bool {
		bodyI := respondentDetails[i].Responses[sortNumAbs-1]
		bodyJ := respondentDetails[j].Responses[sortNumAbs-1]
		if bodyI.Question.Type == "Number" {
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

//InsertRespondent 回答者の追加
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

	err := gormDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&respondent).Error
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to insert a respondent record: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		err = tx.Select("response_id").Last(&respondent).Error
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

func UpdateRespondents(c echo.Context, questionnaireID int, responseID int) error {
	userID := GetUserID(c)

	err := gormDB.
		Model(&Respondents{}).
		Where("user_traqid = ? AND response_id = ?", userID, responseID).
		Update(map[string]interface{}{
			"questionnaire_id": questionnaireID,
		}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func DeleteRespondent(c echo.Context, responseID int) error {
	userID := GetUserID(c)

	err := gormDB.
		Model(&Respondents{}).
		Joins("INNER JOIN administrators ON administrators.questionnaire_id = respondents.questionnaire_id").
		Where("respondents.response_id = ? AND administrators.user_traqid = ? OR respondents.user_traqid = ?", responseID, userID, userID).
		Update(map[string]interface{}{
			"respondents.deleted_at": time.Now(),
		}).Error
	if gorm.IsRecordNotFoundError(err) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = gormDB.
		Where("response_id = ?", responseID).
		Delete(&Response{}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func IsRespondent(c echo.Context, questionnaireID int) (bool, error) {
	userID := GetUserID(c)

	err := gormDB.
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		First(&Respondents{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get response: %w", err))
		return false, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return true, nil
}

func GetRespondentInfos(c echo.Context, userID string, questionnaireIDs ...int) ([]RespondentInfo, error) {
	respondentInfos := []RespondentInfo{}

	query := gormDB.
		Table("respondents").
		Joins("LEFT OUTER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Where("user_traqid = ?", userID)

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

		err := gormDB.ScanRows(rows, respondentInfo)
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to scan responses: %w", err))
			return nil, echo.NewHTTPError(http.StatusInternalServerError)
		}

		respondentInfos = append(respondentInfos, respondentInfo)
	}

	return respondentInfos, nil
}

func GetRespondentDetail(c echo.Context, responseID int) (RespondentDetail, error) {
	userID := GetUserID(c)

	rows, err := gormDB.
		Table("respondents").
		Joins("LEFT OUTER JOIN question ON respondents.questionnaire_id = question.questionnaire_id").
		Joins("LEFT OUTER JOIN response ON respondents.response_id = response.response_id AND respondents.question_id = response.question_id").
		Where("respondents.response_id = ? AND respondents.user_traqid = ?", responseID, userID).
		Select("respondents.questionnaire_id, respondents.modified_at, respondents.submitted_at, question.id, question.type, response.body").
		Rows()
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			c.Logger().Error(err)
			return RespondentDetail{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	isRespondentSetted := false
	respondentDetail := RespondentDetail{}
	responseBodyMap := map[int][]string{}
	for rows.Next() {
		res := struct {
			Respondents `gorm:"embedded"`
			ResponseBody `gorm:"embedded"`
		}{}
		err := gormDB.ScanRows(rows, &res)
		if err != nil {
			return RespondentDetail{}, fmt.Errorf("failed to scan response detail: %w", err)
		}
		if !isRespondentSetted {
			respondentDetail.Respondents = res.Respondents
		}

		respondentDetail.Responses = append(respondentDetail.Responses, ResponseBody{
			Question: res.ResponseBody.Question,
		})

		if res.ResponseBody.Body.Valid {
			responseBodyMap[res.ResponseBody.Question.ID] = append(responseBodyMap[res.ResponseBody.Question.ID], res.ResponseBody.Body.String)
		}
	}

	for i := range respondentDetail.Responses {
		response := &respondentDetail.Responses[i]
		responseBody, ok := responseBodyMap[response.Question.ID]
		switch response.Question.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			if !ok {
				return RespondentDetail{}, errors.New("unexpected no response")
			}
			response.OptionResponse = responseBody
		default:
			if !ok || len(responseBody)==0 {
				response.Body = null.NewString("", false)
			}
			response.Body = null.NewString(responseBody[0], true)
		}
	}

	return respondentDetail, nil
}

func GetRespondentDetails(c echo.Context, questionnaireID int, sort string) ([]RespondentDetail, error) {
	query := gormDB.
		Table("respondents").
		Joins("LEFT OUTER JOIN question ON respondents.questionnaire_id = question.questionnaire_id").
		Joins("LEFT OUTER JOIN response ON respondents.response_id = response.response_id AND respondents.question_id = response.question_id")
	query, sortNum, err := setRespondentsOrder(query, sort)

	rows, err := query.
		Where("respondents.questionnaire_id = ?", questionnaireID).
		Select("respondents.response_id, respondents.modified_at, respondents.submitted_at, question.id, question.type, response.body").
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
			Respondents `gorm:"embedded"`
			ResponseBody `gorm:"embedded"`
		}{}
		err := gormDB.ScanRows(rows, &res)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to scan response detail: %w", err))
		}

		if _, ok := responseBodyMap[res.ResponseID];!ok {
			respondentDetails = append(respondentDetails, RespondentDetail{
				Respondents: res.Respondents,
			})
		}

		responseBodyMap[res.ResponseID] = append(responseBodyMap[res.ResponseID], res.ResponseBody)
	}

	for i := range respondentDetails {
		responseDetail := &respondentDetails[i]
		responseBodies := responseBodyMap[responseDetail.ResponseID]

		responseBodyList := []ResponseBody{}
		bodyMap := map[int][]string{}
		for _,v := range responseBodies {
			if _,ok := bodyMap[v.Question.ID];!ok {
				responseBodyList = append(responseBodyList, ResponseBody{
					Question: v.Question,
				})
			}

			if v.Body.Valid {
				bodyMap[v.Question.ID] = append(bodyMap[v.Question.ID], v.Body.String)
			}
		}

		for i := range responseBodyList {
			responseBody := &responseBodyList[i]
			body, ok := bodyMap[responseBody.Question.ID]
			switch responseBody.Question.Type {
			case "MultipleChoice", "Checkbox", "Dropdown":
				if !ok {
					return nil, errors.New("unexpected no response")
				}
				responseBody.OptionResponse = body
			default:
				if !ok || len(body)==0 {
					responseBody.Body = null.NewString("", false)
				}
				responseBody.Body = null.NewString(body[0], true)
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
