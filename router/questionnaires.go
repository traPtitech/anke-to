package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/guregu/null.v3"

	"github.com/traPtitech/anke-to/model"
)

// GetQuestionnaires GET /questionnaires
func GetQuestionnaires(c echo.Context) error {
	questionnaires, pageMax, err := model.GetQuestionnaires(c, c.QueryParam("nontargeted") == "true")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"page_max":       pageMax,
		"questionnaires": questionnaires,
	})
}

// PostQuestionnaire POST /questionnaires
func PostQuestionnaire(c echo.Context) error {
	req := struct {
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		ResTimeLimit   null.Time `json:"res_time_limit"`
		ResSharedTo    string    `json:"res_shared_to"`
		Targets        []string  `json:"targets"`
		Administrators []string  `json:"administrators"`
	}{}

	// JSONを構造体につける
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	lastID, err := model.InsertQuestionnaire(c, req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo)
	if err != nil {
		return err
	}

	if err := model.InsertTargets(lastID, req.Targets); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := model.InsertAdministrators(c, lastID, req.Administrators); err != nil {
		return err
	}

	timeLimit := "なし"
	if req.ResTimeLimit.Valid {
		timeLimit = req.ResTimeLimit.Time.Format("2006/01/02 15:04")
	}

	targetsMentionText := "なし"
	if len(req.Targets) != 0 {
		targetsMentionText = "@" + strings.Join(req.Targets, " @")
	}

	if err := model.PostMessage(c,
		"### アンケート『"+"["+req.Title+"](https://anke-to.trap.jp/questionnaires/"+
			strconv.Itoa(lastID)+")"+"』が作成されました\n"+
			"#### 管理者\n"+strings.Join(req.Administrators, ",")+"\n"+
			"#### 説明\n"+req.Description+"\n"+
			"#### 回答期限\n"+timeLimit+"\n"+
			"#### 対象者\n"+targetsMentionText+"\n"+
			"#### 回答リンク\n"+
			"https://anke-to.trap.jp/responses/new/"+strconv.Itoa(lastID)); err != nil {
		c.Logger().Error(err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": lastID,
		"title":           req.Title,
		"description":     req.Description,
		"res_time_limit":  req.ResTimeLimit,
		"deleted_at":      "NULL",
		"created_at":      time.Now().Format(time.RFC3339),
		"modified_at":     time.Now().Format(time.RFC3339),
		"res_shared_to":   req.ResSharedTo,
		"targets":         req.Targets,
		"administrators":  req.Administrators,
	})
}

// GetQuestionnaire GET /questionnaires/:questionnaireID
func GetQuestionnaire(c echo.Context) error {
	strQuestionnaireID := c.Param("questionnaireID")
	questionnaireID, err := strconv.Atoi(strQuestionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
	}

	questionnaire, targets, administrators, respondents, err := model.GetQuestionnaireInfo(c, questionnaireID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionnaireID": questionnaire.ID,
		"title":           questionnaire.Title,
		"description":     questionnaire.Description,
		"res_time_limit":  model.NullTimeToString(questionnaire.ResTimeLimit),
		"created_at":      questionnaire.CreatedAt.Format(time.RFC3339),
		"modified_at":     questionnaire.ModifiedAt.Format(time.RFC3339),
		"res_shared_to":   questionnaire.ResSharedTo,
		"targets":         targets,
		"administrators":  administrators,
		"respondents":     respondents,
	})
}

// EditQuestionnaire PATCH /questonnaires/:questionnaireID
func EditQuestionnaire(c echo.Context) error {
	questionnaireID, err := getQuestionnaireID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaireID: %w", err))
	}

	req := struct {
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		ResTimeLimit   null.Time `json:"res_time_limit"`
		ResSharedTo    string    `json:"res_shared_to"`
		Targets        []string  `json:"targets"`
		Administrators []string  `json:"administrators"`
	}{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if req.ResSharedTo == "" {
		req.ResSharedTo = "administrators"
	}

	if err := model.UpdateQuestionnaire(
		c, req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo, questionnaireID); err != nil {
		return err
	}

	if err := model.DeleteTargets(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := model.InsertTargets(questionnaireID, req.Targets); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := model.DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	if err := model.InsertAdministrators(c, questionnaireID, req.Administrators); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestionnaire DELETE /questonnaires/:questionnaireID
func DeleteQuestionnaire(c echo.Context) error {
	questionnaireID, err := getQuestionnaireID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaireID: %w", err))
	}

	if err := model.DeleteQuestionnaire(c, questionnaireID); err != nil {
		return err
	}

	if err := model.DeleteTargets(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := model.DeleteAdministrators(c, questionnaireID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// GetQuestions GET /questionnaires/:questionnaireID/questions
func GetQuestions(c echo.Context) error {
	strQuestionnaireID := c.Param("questionnaireID")
	questionnaireID, err := strconv.Atoi(strQuestionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
	}

	allquestions, err := model.GetQuestions(c, questionnaireID)
	if err != nil {
		return err
	}

	if len(allquestions) == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	type questionInfo struct {
		QuestionID      int      `json:"questionID"`
		PageNum         int      `json:"page_num"`
		QuestionNum     int      `json:"question_num"`
		QuestionType    string   `json:"question_type"`
		Body            string   `json:"body"`
		IsRequired      bool     `json:"is_required"`
		CreatedAt       string   `json:"created_at"`
		Options         []string `json:"options"`
		ScaleLabelRight string   `json:"scale_label_right"`
		ScaleLabelLeft  string   `json:"scale_label_left"`
		ScaleMin        int      `json:"scale_min"`
		ScaleMax        int      `json:"scale_max"`
		RegexPattern    string   `json:"regex_pattern"`
		MinBound        string   `json:"min_bound"`
		MaxBound        string   `json:"max_bound"`
	}
	var ret []questionInfo

	for _, v := range allquestions {
		options := []string{}
		scalelabel := model.ScaleLabels{}
		validation := model.Validations{}
		var err error
		switch v.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			options, err = model.GetOptions(c, v.ID)
		case "LinearScale":
			scalelabel, err = model.GetScaleLabel(v.ID)
		case "Text", "Number":
			validation, err = model.GetValidation(v.ID)
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		ret = append(ret,
			questionInfo{
				QuestionID:      v.ID,
				PageNum:         v.PageNum,
				QuestionNum:     v.QuestionNum,
				QuestionType:    v.Type,
				Body:            v.Body,
				IsRequired:      v.IsRequired,
				CreatedAt:       v.CreatedAt.Format(time.RFC3339),
				Options:         options,
				ScaleLabelRight: scalelabel.ScaleLabelRight,
				ScaleLabelLeft:  scalelabel.ScaleLabelLeft,
				ScaleMin:        scalelabel.ScaleMin,
				ScaleMax:        scalelabel.ScaleMax,
				RegexPattern:    validation.RegexPattern,
				MinBound:        validation.MinBound,
				MaxBound:        validation.MaxBound,
			})
	}

	return c.JSON(http.StatusOK, ret)
}
