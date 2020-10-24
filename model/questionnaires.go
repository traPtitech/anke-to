package model

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
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

//TargettedQuestionnaire targetになっているアンケートの情報
type TargettedQuestionnaire struct {
	Questionnaires
	RespondedAt string `json:"responded_at"`
}

//InsertQuestionnaire アンケートの追加
func InsertQuestionnaire(c echo.Context, title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error) {
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
		c.Logger().Error(fmt.Errorf("failed in the transaction: %w", err))
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return questionnaire.ID, nil
}

//UpdateQuestionnaire アンケートの更新
func UpdateQuestionnaire(c echo.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error {
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
			c.Logger().Error(fmt.Errorf("failed to update a questionnaire record: %w", err))
			return echo.NewHTTPError(http.StatusInternalServerError)
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
		c.Logger().Error(fmt.Errorf("failed to update a questionnaire record: %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

//DeleteQuestionnaire アンケートの削除
func DeleteQuestionnaire(c echo.Context, questionnaireID int) error {
	err := db.Delete(&Questionnaires{ID: questionnaireID}).Error
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

/*GetQuestionnaires アンケートの一覧
2つ目の戻り値はページ数の最大値*/
func GetQuestionnaires(c echo.Context, nontargeted bool) ([]QuestionnaireInfo, int, error) {
	userID := GetUserID(c)
	sort := c.QueryParam("sort")
	search := c.QueryParam("search")
	page := c.QueryParam("page")
	if len(page) == 0 {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to convert the string query parameter 'page'(%s) to integer: %w", page, err))
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest)
	}
	if pageNum <= 0 {
		c.Logger().Error(errors.New("page cannot be less than 0"))
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest)
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

	count := 0
	err = query.Count(&count).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to retrieve the number of questionnaires: %w", err))
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	if count == 0 {
		c.Logger().Error(fmt.Errorf("failed to get the targeted questionnaires: %w", err))
		return nil, 0, echo.NewHTTPError(http.StatusNotFound)
	}
	pageMax := (count + 19) / 20

	if pageNum > pageMax {
		c.Logger().Error("too large page number")
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest)
	}

	offset := (pageNum - 1) * 20
	query = query.Limit(20).Offset(offset)

	err = query.
		Group("questionnaires.id").
		Select("questionnaires.*, (targets.user_traqid = ? OR targets.user_traqid = 'traP') AS is_targeted", userID).
		Find(&questionnaires).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Logger().Error(fmt.Errorf("failed to get the targeted questionnaires: %w", err))
		return nil, 0, echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get the targeted questionnaires: %w", err))
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	if len(search) != 0 {
		r, err := regexp.Compile(strings.ToLower(search))
		if err != nil {
			c.Logger().Error("invalid search param regexp")
			return nil, 0, echo.NewHTTPError(http.StatusBadRequest)
		}

		retQuestionnaires := make([]QuestionnaireInfo, 0, len(questionnaires))
		for _, q := range questionnaires {
			if search != "" && !r.MatchString(strings.ToLower(q.Title)) {
				continue
			}

			retQuestionnaires = append(retQuestionnaires, q)
		}

		questionnaires = retQuestionnaires
	}

	return questionnaires, pageMax, nil
}

//GetQuestionnaireInfo アンケートの詳細な情報取得
func GetQuestionnaireInfo(c echo.Context, questionnaireID int) (*Questionnaires, []string, []string, []string, error) {
	questionnaire := Questionnaires{}
	targets := []string{}
	administrators := []string{}
	respondents := []string{}

	err := db.
		Model(&Questionnaires{}).
		Where("questionnaires.id = ?", questionnaireID).
		First(&questionnaire).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get a questionnaire: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, echo.NewHTTPError(http.StatusNotFound)
		}
		return nil, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = db.
		Table("targets").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &targets).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get targets: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, echo.NewHTTPError(http.StatusNotFound)
		}
		return nil, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = db.
		Table("administrators").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &administrators).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get administrators: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, echo.NewHTTPError(http.StatusNotFound)
		}
		return nil, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = db.
		Table("respondents").
		Where("questionnaire_id = ?", questionnaire.ID).
		Pluck("user_traqid", &respondents).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get respondents: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, nil, nil, echo.NewHTTPError(http.StatusNotFound)
		}
		return nil, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return &questionnaire, targets, administrators, respondents, nil
}

//GetTargettedQuestionnaires targetになっているアンケートの取得
func GetTargettedQuestionnaires(c echo.Context, userID string, answered string) ([]TargettedQuestionnaire, error) {
	sort := c.QueryParam("sort")
	query := db.
		Table("questionnaires").
		Where("questionnaires.res_time_limit > ? OR questionnaires.res_time_limit IS NULL", time.Now()).
		Joins("INNER JOIN targets ON questionnaires.id = targets.questionnaire_id").
		Where("targets.user_traqid = ? OR targets.user_traqid = 'traP'", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id").
		Where("respondents.user_traqid = ? OR respondents.user_traqid IS NULL", userID).
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
		c.Logger().Error(fmt.Errorf("failed to get the targeted questionnaires: %w", err))
		return nil, echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get the targeted questionnaires: %w", err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return questionnaires, nil
}

//GetQuestionnaireLimit アンケートの回答期限の取得
func GetQuestionnaireLimit(c echo.Context, questionnaireID int) (string, error) {
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
		c.Logger().Error(err)
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}

	return NullTimeToString(res.ResTimeLimit), nil
}

//GetResShared アンケートの回答の公開範囲の取得
func GetResShared(c echo.Context, questionnaireID int) (string, error) {
	res := Questionnaires{}

	err := db.
		Model(Questionnaires{}).
		Where("id = ?", questionnaireID).
		Select("res_shared_to").
		Scan(&res).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get resShared: %w", err))
		if gorm.IsRecordNotFoundError(err) {
			return "", echo.NewHTTPError(http.StatusNotFound)
		}
		return "", echo.NewHTTPError(http.StatusInternalServerError)
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
