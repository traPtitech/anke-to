package model

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"
)

type Questionnaires struct {
	ID           int            `json:"questionnaireID" gorm:"type:int(11);AUTO_INCREMENT;NOT NULL;"`
	Title        string         `json:"title"           gorm:"type:char(50);NOT NULL;UNIQUE;"`
	Description  string         `json:"description"     gorm:"type:text;NOT NULL;"`
	ResTimeLimit null.Time `json:"res_time_limit"  gorm:"type:timestamp;DEFAULT:NULL;"`
	DeletedAt    null.Time `json:"deleted_at"      gorm:"type:timestamp;DEFAULT:NULL;"`
	ResSharedTo  string         `json:"res_shared_to"   gorm:"type:char(30);NOT NULL;DEFAULT:administrators;"`
	CreatedAt    time.Time      `json:"created_at"      gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;"`
	ModifiedAt   time.Time      `json:"modified_at"     gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;"`
}

type QuestionnairesInfo struct {
	Questionnaires
	IsTargeted   bool   `json:"is_targeted" gorm:"type:boolean"`
}

type TargettedQuestionnaires struct {
	ID           int    `json:"questionnaireID"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ResTimeLimit string `json:"res_time_limit"`
	ResSharedTo  string `json:"res_shared_to"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
	RespondedAt  string `json:"responded_at"`
}

func SetQuestionnairesOrder(query *gorm.DB, sort string) (*gorm.DB, error) {
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

// エラーが起きれば(nil, err)
// 起こらなければ(allquestionnaires, nil)を返す
func GetAllQuestionnaires(c echo.Context) ([]Questionnaires, error) {
	// query parametar
	sort := c.QueryParam("sort")

	// アンケート一覧の配列
	allquestionnaires := []Questionnaires{}

	query := gormDB.Where("deleted_at IS NULL")

	query, err := SetQuestionnairesOrder(query, sort)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := query.Find(&allquestionnaires).Error; err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return allquestionnaires, nil
}

/*
アンケートの一覧
2つ目の戻り値はページ数の最大
*/
func GetQuestionnaires(c echo.Context, nontargeted bool) ([]QuestionnairesInfo, int, error) {
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

	questionnaires := make([]QuestionnairesInfo, 0, 20)

	query := gormDB.Table("questionnaires").Joins("LEFT OUTER JOIN targets ON questionnaires.id = targets.questionnaire_id")

	query, err = SetQuestionnairesOrder(query, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to set the order of the questionnaire table: %w", err)
	}

	if nontargeted {
		query = query.Where("targets.user_traqid != ? AND targets.user_traqid != 'traP'", userID)
	}

	count := 0
	err = query.Count(&count).Error
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to retrieve the number of questionnaires: %w", err))
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError)
	}
	pageMax := (count + 19)/20

	if pageNum > pageMax {
		c.Logger().Error("too large page number")
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest)
	}

	offset := (pageNum - 1) * 20
	query = query.Limit(20).Offset(offset)

	err = query.Select("questionnaires.*, (targets.user_traqid != ? AND targets.user_traqid != 'traP') AS is_targeted", userID).Find(&questionnaires).Error
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

		retQuestionnaires := make([]QuestionnairesInfo, 0, len(questionnaires))
		for _, q := range questionnaires {
			if search != "" && !r.MatchString(strings.ToLower(q.Title)) {
				continue
			}

			retQuestionnaires = append(retQuestionnaires,q)
		}

		questionnaires = retQuestionnaires
	}

	return questionnaires, pageMax, nil
}

func GetQuestionnaire(c echo.Context, questionnaireID int) (Questionnaires, error) {
	questionnaire := Questionnaires{}

	err := gormDB.Where("id = ? AND deleted_at IS NULL", questionnaireID).First(&questionnaire).Error
	if err != nil {
		c.Logger().Error(err)
		if gorm.IsRecordNotFoundError(err) {
			return Questionnaires{}, echo.NewHTTPError(http.StatusNotFound)
		} else {
			return Questionnaires{}, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return questionnaire, nil
}

func GetQuestionnaireInfo(c echo.Context, questionnaireID int) (Questionnaires, []string, []string, []string, error) {
	questionnaire, err := GetQuestionnaire(c, questionnaireID)
	if err != nil {
		return Questionnaires{}, nil, nil, nil, err
	}

	targets, err := GetTargets(c, questionnaireID)
	if err != nil {
		return Questionnaires{}, nil, nil, nil, err
	}

	administrators, err := GetAdministrators(c, questionnaireID)
	if err != nil {
		return Questionnaires{}, nil, nil, nil, err
	}

	respondents, err := GetRespondents(c, questionnaireID)
	if err != nil {
		return Questionnaires{}, nil, nil, nil, err
	}
	return questionnaire, targets, administrators, respondents, nil
}

func GetQuestionnaireLimit(c echo.Context, questionnaireID int) (string, error) {
	res := struct {
		Title        string         `db:"title"`
		ResTimeLimit null.Time `db:"res_time_limit"`
	}{}
	if err := db.Get(&res,
		"SELECT res_time_limit FROM questionnaires WHERE id = ? AND deleted_at IS NULL",
		questionnaireID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} else {
			c.Logger().Error(err)
			return "", echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return NullTimeToString(res.ResTimeLimit), nil
}

func GetTitleAndLimit(c echo.Context, questionnaireID int) (string, string, error) {
	res := struct {
		Title        string         `db:"title"`
		ResTimeLimit null.Time `db:"res_time_limit"`
	}{}
	if err := db.Get(&res,
		"SELECT title, res_time_limit FROM questionnaires WHERE id = ? AND deleted_at IS NULL",
		questionnaireID); err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil
		} else {
			c.Logger().Error(err)
			return "", "", echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return res.Title, NullTimeToString(res.ResTimeLimit), nil
}

func InsertQuestionnaire(c echo.Context, title string, description string, resTimeLimit string, resSharedTo string) (int, error) {
	var result sql.Result

	if resTimeLimit == "" || resTimeLimit == "NULL" {
		resTimeLimit = "NULL"
		var err error
		result, err = db.Exec(
			`INSERT INTO questionnaires (title, description, res_shared_to, created_at, modified_at)
			VALUES (?, ?, ?, ?, ?)`,
			title, description, resSharedTo, time.Now(), time.Now())
		if err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		var err error
		result, err = db.Exec(
			`INSERT INTO questionnaires (title, description, res_time_limit, res_shared_to, created_at, modified_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
			title, description, resTimeLimit, resSharedTo, time.Now(), time.Now())
		if err != nil {
			c.Logger().Error(err)
			return 0, echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		c.Logger().Error(err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return int(lastID), nil
}

func UpdateQuestionnaire(c echo.Context, title string, description string, resTimeLimit string, resSharedTo string, questionnaireID int) error {
	if resTimeLimit == "" || resTimeLimit == "NULL" {
		resTimeLimit = "NULL"
		if _, err := db.Exec(
			`UPDATE questionnaires SET title = ?, description = ?, res_time_limit = NULL,
			res_shared_to = ?, modified_at = ? WHERE id = ?`,
			title, description, resSharedTo, time.Now(), questionnaireID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := db.Exec(
			`UPDATE questionnaires SET title = ?, description = ?, res_time_limit = ?,
			res_shared_to = ?, modified_at = ? WHERE id = ?`,
			title, description, resTimeLimit, resSharedTo, time.Now(), questionnaireID); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

func DeleteQuestionnaire(c echo.Context, questionnaireID int) error {
	if _, err := db.Exec(
		"UPDATE questionnaires SET deleted_at = ? WHERE id = ?",
		time.Now(), questionnaireID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func GetResShared(c echo.Context, questionnaireID int) (string, error) {
	resSharedTo := ""
	if err := db.Get(&resSharedTo,
		`SELECT res_shared_to FROM questionnaires WHERE deleted_at IS NULL AND id = ?`,
		questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return "", echo.NewHTTPError(http.StatusNotFound)
		} else {
			return "", echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return resSharedTo, nil
}

func GetTargettedQuestionnaires(c echo.Context) ([]TargettedQuestionnaires, error) {
	// 全てのアンケート
	allquestionnaires, err := GetAllQuestionnaires(c)
	if err != nil {
		return nil, err
	}

	// 自分がtargetになっているアンケート
	targetedQuestionnaireID, err := GetTargettedQuestionnaireID(c)
	if err != nil {
		return nil, err
	}

	ret := []TargettedQuestionnaires{}
	for _, v := range allquestionnaires {
		var targeted = false
		for _, w := range targetedQuestionnaireID {
			if w == v.ID {
				targeted = true
			}
		}
		if !targeted {
			continue
		}
		respondedAt, err := RespondedAt(c, v.ID)
		if err != nil {
			return nil, err
		}

		// アンケートの期限がNULLでなく期限を過ぎていたら次へ
		if v.ResTimeLimit.Valid && time.Now().After(v.ResTimeLimit.Time) {
			continue
		}

		ret = append(ret,
			TargettedQuestionnaires{
				ID:           v.ID,
				Title:        v.Title,
				Description:  v.Description,
				ResTimeLimit: NullTimeToString(v.ResTimeLimit),
				ResSharedTo:  v.ResSharedTo,
				CreatedAt:    v.CreatedAt.Format(time.RFC3339),
				ModifiedAt:   v.ModifiedAt.Format(time.RFC3339),
				RespondedAt:  respondedAt,
			})
	}

	// アンケートが1つも無い場合
	if len(ret) == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound)
	}

	// 回答期限が近い順に
	sort.Slice(ret, func(i, j int) bool {
		if ret[i].ResTimeLimit == "NULL" && ret[j].ResTimeLimit == "NULL" {
			return ret[i].ModifiedAt > ret[j].ModifiedAt
		} else if ret[i].ResTimeLimit == "NULL" {
			return false
		} else if ret[j].ResTimeLimit == "NULL" {
			return true
		} else {
			return ret[i].ResTimeLimit < ret[j].ResTimeLimit
		}
	})

	return ret, nil
}

func GetTargettedQuestionnairesBytraQID(c echo.Context, traQID string) ([]TargettedQuestionnaires, error) {
	// 全てのアンケート
	allquestionnaires, err := GetAllQuestionnaires(c)
	if err != nil {
		return nil, err
	}

	// 指定したtraQIDがtargetになっているアンケート
	targetedQuestionnaireID, err := GetTargettedQuestionnaireIDBytraQID(c, traQID)
	if err != nil {
		return nil, err
	}

	ret := []TargettedQuestionnaires{}
	for _, v := range allquestionnaires {
		var targeted = false
		for _, w := range targetedQuestionnaireID {
			if w == v.ID {
				targeted = true
			}
		}
		if !targeted {
			continue
		}

		// 回答済みなら次へ
		respondedAt, err := RespondedAtBytraQID(c, v.ID, traQID)
		if err != nil {
			return nil, err
		}
		if respondedAt != "NULL" {
			continue
		}

		// アンケートの期限がNULLまたは期限を過ぎていたら次へ
		if !v.ResTimeLimit.Valid || time.Now().After(v.ResTimeLimit.Time) {
			continue
		}

		ret = append(ret,
			TargettedQuestionnaires{
				ID:           v.ID,
				Title:        v.Title,
				Description:  v.Description,
				ResTimeLimit: NullTimeToString(v.ResTimeLimit),
				ResSharedTo:  v.ResSharedTo,
				CreatedAt:    v.CreatedAt.Format(time.RFC3339),
				ModifiedAt:   v.ModifiedAt.Format(time.RFC3339),
				RespondedAt:  respondedAt,
			})
	}

	// アンケートが1つも無い場合
	if len(ret) == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound)
	}

	// 回答期限が近い順に
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ResTimeLimit < ret[j].ResTimeLimit
	})

	return ret, nil
}
