package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

//Respondents respondentsテーブルの構造体
type Respondents struct {
	ResponseID int `json:"responseID" gorm:"type:int(11);NOT NULL;PRIMARY_KEY;AUTO_INCREMENT"`
	QuestionnaireID int `json:"questionnaireID" gorm:"type:int(11);NOT NULL;"`
	UserTraqid string `json:"user_traq_id,omitempty" gorm:"type:char(30);NOT NULL;"`
	ModifiedAt time.Time `json:"modified_at" gorm:"type:timestamp;NOT NULL;DEFAULT CURRENT_TIMESTAMP;"`
	SubmittedAt null.Time `json:"submitted_at" gorm:"type:timestamp;"`
	DeletedAt null.Time `gorm:"type:timestamp;"`
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (*Respondents) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedAt", time.Now())

	return nil
}

type RespondentDetail struct {
	Title           string `json:"questionnaire_title"`
	ResTimeLimit    string `json:"res_time_limit"`
	Respondents
}

//InsertRespondent 回答者の追加
func InsertRespondent(c echo.Context, questionnaireID int, submitedAt null.Time) (int, error) {
	userID := GetUserID(c)

	var respondent Respondents
	if submitedAt.Valid {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid: userID,
			SubmittedAt: submitedAt,
		}
	} else {
		respondent = Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid: userID,
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

func GetRespondentDetails(c echo.Context, questionnaireIDs... int) ([]RespondentDetail, error) {
	userID := GetUserID(c)
	respondentDetails := []RespondentDetail{}

	query := gormDB.
		Table("respondents").
		Joins("LEFT OUTER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Where("user_traqid = ?", userID)
	
	if len(questionnaireIDs) != 0 {
		questionnaireID := questionnaireIDs[0]
		query = query.Where("questionnaire_id = ?", questionnaireID)
	}

	rows,err := query.
		Select("respondents.questionnaire_id, respondents.response_id, respondents.modified_at, respondents.submitted_at, questionnaires.title, "+
		"questionnaires.res_time_limit").
		Rows()
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get my responses: %w", err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	for rows.Next() {
		respondentDetail := RespondentDetail{
			Respondents: Respondents{},
		}

		err := gormDB.ScanRows(rows, respondentDetail)
		if err != nil {
			c.Logger().Error(fmt.Errorf("failed to scan responses: %w", err))
			return nil, echo.NewHTTPError(http.StatusInternalServerError)
		}

		respondentDetails = append(respondentDetails, respondentDetail)
	}

	return respondentDetails, nil
}
