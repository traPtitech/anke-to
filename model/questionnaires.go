package model

import (
	"net/http"
	"strconv"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Questionnaires struct {
	ID           int            `json:"questionnaireID" db:"id"`
	Title        string         `json:"title"           db:"title"`
	Description  string         `json:"description"     db:"description"`
	ResTimeLimit mysql.NullTime `json:"res_time_limit"  db:"res_time_limit"`
	DeletedAt    mysql.NullTime `json:"deleted_at"      db:"deleted_at"`
	ResSharedTo  string         `json:"res_shared_to"   db:"res_shared_to"`
	CreatedAt    time.Time      `json:"created_at"      db:"created_at"`
	ModifiedAt   time.Time      `json:"modified_at"     db:"modified_at"`
}

// エラーが起きれば(nil, err)
// 起こらなければ(allquestions, nil)を返す
func GetAllQuestionnaires(c echo.Context) ([]Questionnaires, error) {
	// query parametar
	sort := c.QueryParam("sort")
	page := c.QueryParam("page")

	if page == "" {
		page = "1"
	}
	num, err := strconv.Atoi(page)
	if err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}

	var list = map[string]string{
		"":             "",
		"created_at":   "ORDER BY created_at",
		"-created_at":  "ORDER BY created_at DESC",
		"title":        "ORDER BY title",
		"-title":       "ORDER BY title DESC",
		"modified_at":  "ORDER BY modified_at",
		"-modified_at": "ORDER BY modified_at DESC",
	}
	// アンケート一覧の配列
	allquestionnaires := []Questionnaires{}

	if err := DB.Select(&allquestionnaires,
		"SELECT * FROM questionnaires WHERE deleted_at IS NULL "+list[sort]+" lIMIT 20 OFFSET "+strconv.Itoa(20*(num-1))); err != nil {
		c.Logger().Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return allquestionnaires, nil
}

func GetQuestionnaires(c echo.Context, targettype TargetType) error {
	allquestionnaires, err := GetAllQuestionnaires(c)
	if err != nil {
		return err
	}

	userID := GetUserID(c)

	targetedQuestionnaireID := []int{}
	if err := DB.Select(&targetedQuestionnaireID,
		"SELECT questionnaire_id FROM targets WHERE user_traqid = ?", userID); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	type questionnairesInfo struct {
		ID           int       `json:"questionnaireID"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		ResTimeLimit string    `json:"res_time_limit"`
		ResSharedTo  string    `json:"res_shared_to"`
		CreatedAt    time.Time `json:"created_at"`
		ModifiedAt   time.Time `json:"modified_at"`
		IsTargeted   bool      `json:"is_targeted"`
	}
	var ret []questionnairesInfo

	for _, v := range allquestionnaires {
		var targeted = false
		for _, w := range targetedQuestionnaireID {
			if w == v.ID {
				targeted = true
			}
		}
		if (targettype == TargetType(Targeted) && !targeted) || (targettype == TargetType(Nontargeted) && targeted) {
			continue
		}
		ret = append(ret,
			questionnairesInfo{
				ID:           v.ID,
				Title:        v.Title,
				Description:  v.Description,
				ResTimeLimit: TimeConvert(v.ResTimeLimit),
				ResSharedTo:  v.ResSharedTo,
				CreatedAt:    v.CreatedAt,
				ModifiedAt:   v.ModifiedAt,
				IsTargeted:   targeted})
	}

	if len(ret) == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// 構造体の定義で書いたJSONのキーで変換される
	return c.JSON(http.StatusOK, ret)
}

func GetTitleAndLimit(c echo.Context, questionnaireID int) (string, string, error) {
	res := struct {
		Title        string `db:"title"`
		ResTimeLimit string `db:"res_time_limit"`
	}{}
	if err := DB.Get(&res,
		"SELECT title, res_time_limit FROM questionnaires WHERE id = ? AND deleted_at IS NULL",
		questionnaireID); err != nil {
		c.Logger().Error(err)
		if err == sql.ErrNoRows {
			return "", "", echo.NewHTTPError(http.StatusNotFound)
		} else {
			return "", "", echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return res.Title, res.ResTimeLimit, nil
}
