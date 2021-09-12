package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v3"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
)

// Questionnaire Questionnaireの構造体
type Questionnaire struct {
	model.IQuestionnaire
	model.ITarget
	model.IAdministrator
	model.IQuestion
	model.IOption
	model.IScaleLabel
	model.IValidation
	traq.IWebhook
}

const MaxTitleLength = 50

// NewQuestionnaire Questionnaireのコンストラクタ
func NewQuestionnaire(questionnaire model.IQuestionnaire, target model.ITarget, administrator model.IAdministrator, question model.IQuestion, option model.IOption, scaleLabel model.IScaleLabel, validation model.IValidation, webhook traq.IWebhook) *Questionnaire {
	return &Questionnaire{
		IQuestionnaire: questionnaire,
		ITarget:        target,
		IAdministrator: administrator,
		IQuestion:      question,
		IOption:        option,
		IScaleLabel:    scaleLabel,
		IValidation:    validation,
		IWebhook:       webhook,
	}
}

// GetQuestionnaires GET /questionnaires
func (q *Questionnaire) GetQuestionnaires(c echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	sort := c.QueryParam("sort")
	search := c.QueryParam("search")
	page := c.QueryParam("page")
	if len(page) == 0 {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to convert the string query parameter 'page'(%s) to integer: %w", page, err))
	}
	if pageNum <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("page cannot be less than 0"))
	}
	questionnaires, pageMax, err := q.IQuestionnaire.GetQuestionnaires(userID, sort, search, pageNum, c.QueryParam("nontargeted") == "true")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		} else if errors.Is(err, model.ErrTooLargePageNum) || errors.Is(err, model.ErrInvalidRegex) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"page_max":       pageMax,
		"questionnaires": questionnaires,
	})
}

type PostAndEditQuestionnaireRequest struct {
	Title          string    `json:"title" validate:"required,max=50"`
	Description    string    `json:"description"`
	ResTimeLimit   null.Time `json:"res_time_limit"`
	ResSharedTo    string    `json:"res_shared_to" validate:"required,oneof=administrators respondents public"`
	Targets        []string  `json:"targets" validate:"dive,max=32"`
	Administrators []string  `json:"administrators" validate:"required,min=1,dive,max=32"`
}

// PostQuestionnaire POST /questionnaires
func (q *Questionnaire) PostQuestionnaire(c echo.Context) error {
	req := PostAndEditQuestionnaireRequest{}

	// JSONを構造体につける
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	validate, err := getValidator(c)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get validator: %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(), req)
	if err != nil {
		c.Logger().Info(fmt.Errorf("failed to validate: %w", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	lastID, err := q.InsertQuestionnaire(req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo)
	if err != nil {
		return err
	}

	if err := q.InsertTargets(lastID, req.Targets); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.InsertAdministrators(lastID, req.Administrators); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	timeLimit := "なし"
	if req.ResTimeLimit.Valid {
		timeLimit = req.ResTimeLimit.Time.Local().Format("2006/01/02 15:04")
	}

	targetsMentionText := "なし"
	if len(req.Targets) != 0 {
		targetsMentionText = "@" + strings.Join(req.Targets, " @")
	}

	if err := q.PostMessage(
		"### アンケート『" + "[" + req.Title + "](https://anke-to.trap.jp/questionnaires/" +
			strconv.Itoa(lastID) + ")" + "』が作成されました\n" +
			"#### 管理者\n" + strings.Join(req.Administrators, ",") + "\n" +
			"#### 説明\n" + req.Description + "\n" +
			"#### 回答期限\n" + timeLimit + "\n" +
			"#### 対象者\n" + targetsMentionText + "\n" +
			"#### 回答リンク\n" +
			"https://anke-to.trap.jp/responses/new/" + strconv.Itoa(lastID)); err != nil {
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
func (q *Questionnaire) GetQuestionnaire(c echo.Context) error {
	strQuestionnaireID := c.Param("questionnaireID")
	questionnaireID, err := strconv.Atoi(strQuestionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
	}

	questionnaire, targets, administrators, respondents, err := q.GetQuestionnaireInfo(questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"questionnaireID": questionnaire.ID,
		"title":           questionnaire.Title,
		"description":     questionnaire.Description,
		"res_time_limit":  questionnaire.ResTimeLimit,
		"created_at":      questionnaire.CreatedAt.Format(time.RFC3339),
		"modified_at":     questionnaire.ModifiedAt.Format(time.RFC3339),
		"res_shared_to":   questionnaire.ResSharedTo,
		"targets":         targets,
		"administrators":  administrators,
		"respondents":     respondents,
	})
}

// EditQuestionnaire PATCH /questonnaires/:questionnaireID
func (q *Questionnaire) EditQuestionnaire(c echo.Context) error {
	questionnaireID, err := getQuestionnaireID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaireID: %w", err))
	}

	req := PostAndEditQuestionnaireRequest{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	validate, err := getValidator(c)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get validator: %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(), req)
	if err != nil {
		c.Logger().Info(fmt.Errorf("failed to validate: %w", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.ResTimeLimit.Valid {
		isBefore := req.ResTimeLimit.ValueOrZero().Before(time.Now())
		if isBefore {
			c.Logger().Info(fmt.Sprintf(": %+v", req.ResTimeLimit))
			return echo.NewHTTPError(http.StatusBadRequest, "res time limit is before now")
		}
	}

	if err := q.UpdateQuestionnaire(
		req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo, questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteTargets(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.InsertTargets(questionnaireID, req.Targets); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteAdministrators(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.InsertAdministrators(questionnaireID, req.Administrators); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestionnaire DELETE /questonnaires/:questionnaireID
func (q *Questionnaire) DeleteQuestionnaire(c echo.Context) error {
	questionnaireID, err := getQuestionnaireID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaireID: %w", err))
	}

	if err := q.IQuestionnaire.DeleteQuestionnaire(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteTargets(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := q.DeleteAdministrators(questionnaireID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetQuestions GET /questionnaires/:questionnaireID/questions
func (q *Questionnaire) GetQuestions(c echo.Context) error {
	strQuestionnaireID := c.Param("questionnaireID")
	questionnaireID, err := strconv.Atoi(strQuestionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
	}

	allquestions, err := q.IQuestion.GetQuestions(questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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

	optionIDs := []int{}
	scaleLabelIDs := []int{}
	validationIDs := []int{}
	for _, question := range allquestions {
		switch question.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			optionIDs = append(optionIDs, question.ID)
		case "LinearScale":
			scaleLabelIDs = append(scaleLabelIDs, question.ID)
		case "Text", "Number":
			validationIDs = append(validationIDs, question.ID)
		}
	}

	options, err := q.GetOptions(optionIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	optionMap := make(map[int][]string, len(options))
	for _, option := range options {
		optionMap[option.QuestionID] = append(optionMap[option.QuestionID], option.Body)
	}

	scaleLabels, err := q.GetScaleLabels(scaleLabelIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, label := range scaleLabels {
		scaleLabelMap[label.QuestionID] = label
	}

	validations, err := q.GetValidations(validationIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	validationMap := make(map[int]model.Validations, len(validations))
	for _, validation := range validations {
		validationMap[validation.QuestionID] = validation
	}

	for _, v := range allquestions {
		options := []string{}
		scalelabel := model.ScaleLabels{}
		validation := model.Validations{}
		switch v.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			var ok bool
			options, ok = optionMap[v.ID]
			if !ok {
				options = []string{}
			}
		case "LinearScale":
			var ok bool
			scalelabel, ok = scaleLabelMap[v.ID]
			if !ok {
				scalelabel = model.ScaleLabels{}
			}
		case "Text", "Number":
			var ok bool
			validation, ok = validationMap[v.ID]
			if !ok {
				validation = model.Validations{}
			}
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
			},
		)
	}

	return c.JSON(http.StatusOK, ret)
}
