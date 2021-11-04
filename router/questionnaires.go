package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

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
	model.ITransaction
	traq.IWebhook
}

const MaxTitleLength = 50

// NewQuestionnaire Questionnaireのコンストラクタ
func NewQuestionnaire(
	questionnaire model.IQuestionnaire,
	target model.ITarget,
	administrator model.IAdministrator,
	question model.IQuestion,
	option model.IOption,
	scaleLabel model.IScaleLabel,
	validation model.IValidation,
	transaction model.ITransaction,
	webhook traq.IWebhook,
) *Questionnaire {
	return &Questionnaire{
		IQuestionnaire: questionnaire,
		ITarget:        target,
		IAdministrator: administrator,
		IQuestion:      question,
		IOption:        option,
		IScaleLabel:    scaleLabel,
		IValidation:    validation,
		ITransaction:   transaction,
		IWebhook:       webhook,
	}
}

type GetQuestionnairesQueryParam struct {
	Sort        string `validate:"omitempty,oneof=created_at -created_at title -title modified_at -modified_at"`
	Search      string `validate:"omitempty"`
	Page        string `validate:"omitempty,number,min=0"`
	Nontargeted string `validate:"omitempty,boolean"`
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
	nontargeted := c.QueryParam("nontargeted")

	p := GetQuestionnairesQueryParam{
		Sort:        sort,
		Search:      search,
		Page:        page,
		Nontargeted: nontargeted,
	}

	validate, err := getValidator(c)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get validator:%w", err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(), p)
	if err != nil {
		c.Logger().Info(fmt.Errorf("failed to validate:%w", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

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

	var nontargetedBool bool
	if len(nontargeted) != 0 {
		nontargetedBool, err = strconv.ParseBool(nontargeted)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to convert the string query parameter 'nontargeted'(%s) to bool: %w", nontargeted, err))
		}
	} else {
		nontargetedBool = false
	}

	questionnaires, pageMax, err := q.IQuestionnaire.GetQuestionnaires(c.Request().Context(), userID, sort, search, pageNum, nontargetedBool)
	if err != nil {
		if errors.Is(err, model.ErrTooLargePageNum) || errors.Is(err, model.ErrInvalidRegex) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if errors.Is(err, model.ErrDeadlineExceeded) {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "deadline exceeded")
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
	err := c.Bind(&req)
	if err != nil {
		c.Logger().Infof("invalid request body: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	validate, err := getValidator(c)
	if err != nil {
		c.Logger().Errorf("failed to get validator: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(), req)
	if err != nil {
		c.Logger().Infof("failed to validate: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.ResTimeLimit.Valid {
		isBefore := req.ResTimeLimit.ValueOrZero().Before(time.Now())
		if isBefore {
			c.Logger().Infof("invalid resTimeLimit: %+v", req.ResTimeLimit)
			return echo.NewHTTPError(http.StatusBadRequest, "res time limit is before now")
		}
	}

	var questionnaireID int
	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		questionnaireID, err = q.InsertQuestionnaire(ctx, req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo)
		if err != nil {
			c.Logger().Errorf("failed to insert a questionnaire: %w", err)
			return err
		}

		err := q.InsertTargets(ctx, questionnaireID, req.Targets)
		if err != nil {
			c.Logger().Errorf("failed to insert targets: %w", err)
			return err
		}

		err = q.InsertAdministrators(ctx, questionnaireID, req.Administrators)
		if err != nil {
			c.Logger().Errorf("failed to insert administrators: %w", err)
			return err
		}

		message := createQuestionnaireMessage(
			questionnaireID,
			req.Title,
			req.Description,
			req.Administrators,
			req.ResTimeLimit,
			req.Targets,
		)
		err = q.PostMessage(message)
		if err != nil {
			c.Logger().Errorf("failed to post message: %w", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to post message to traQ")
		}

		return nil
	})
	if err != nil {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			return httpError
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create a questionnaire")
	}

	now := time.Now()
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": questionnaireID,
		"title":           req.Title,
		"description":     req.Description,
		"res_time_limit":  req.ResTimeLimit,
		"deleted_at":      "NULL",
		"created_at":      now.Format(time.RFC3339),
		"modified_at":     now.Format(time.RFC3339),
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

	questionnaire, targets, administrators, respondents, err := q.GetQuestionnaireInfo(c.Request().Context(), questionnaireID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
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

// PostByQuestionnaireID POST /questionnaires/:questionnaireID/questions
func (q *Questionnaire) PostByQuestionnaireID(c echo.Context) error {
	strQuestionnaireID := c.Param("questionnaireID")
	questionnaireID, err := strconv.Atoi(strQuestionnaireID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
	}
	req := PostAndEditQuestionRequest{}
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

	switch req.QuestionType {
	case "Text":
		//正規表現のチェック
		if _, err := regexp.Compile(req.RegexPattern); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	case "Number":
		//数字か，min<=maxになってるか
		if err := q.CheckNumberValid(req.MinBound, req.MaxBound); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	lastID, err := q.InsertQuestion(c.Request().Context(), questionnaireID, req.PageNum, req.QuestionNum, req.QuestionType, req.Body, req.IsRequired)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch req.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		for i, v := range req.Options {
			if err := q.InsertOption(c.Request().Context(), lastID, i+1, v); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	case "LinearScale":
		if err := q.InsertScaleLabel(c.Request().Context(), lastID,
			model.ScaleLabels{
				ScaleLabelLeft:  req.ScaleLabelLeft,
				ScaleLabelRight: req.ScaleLabelRight,
				ScaleMax:        req.ScaleMax,
				ScaleMin:        req.ScaleMin,
			}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "Text", "Number":
		if err := q.InsertValidation(c.Request().Context(), lastID,
			model.Validations{
				RegexPattern: req.RegexPattern,
				MinBound:     req.MinBound,
				MaxBound:     req.MaxBound,
			}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionID":        int(lastID),
		"questionnaireID":   questionnaireID,
		"question_type":     req.QuestionType,
		"question_num":      req.QuestionNum,
		"page_num":          req.PageNum,
		"body":              req.Body,
		"is_required":       req.IsRequired,
		"options":           req.Options,
		"scale_label_right": req.ScaleLabelRight,
		"scale_label_left":  req.ScaleLabelLeft,
		"scale_max":         req.ScaleMax,
		"scale_min":         req.ScaleMin,
		"regex_pattern":     req.RegexPattern,
		"min_bound":         req.MinBound,
		"max_bound":         req.MaxBound,
	})
}

// EditQuestionnaire PATCH /questionnaires/:questionnaireID
func (q *Questionnaire) EditQuestionnaire(c echo.Context) error {
	questionnaireID, err := getQuestionnaireID(c)
	if err != nil {
		c.Logger().Errorf("failed to get questionnaireID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	req := PostAndEditQuestionnaireRequest{}

	err = c.Bind(&req)
	if err != nil {
		c.Logger().Infof("failed to bind request: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	validate, err := getValidator(c)
	if err != nil {
		c.Logger().Errorf("failed to get validator: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = validate.StructCtx(c.Request().Context(), req)
	if err != nil {
		c.Logger().Infof("failed to validate: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		err = q.UpdateQuestionnaire(ctx, req.Title, req.Description, req.ResTimeLimit, req.ResSharedTo, questionnaireID)
		if err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			c.Logger().Errorf("failed to update questionnaire: %w", err)
			return err
		}

		err = q.DeleteTargets(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete targets: %w", err)
			return err
		}

		err = q.InsertTargets(ctx, questionnaireID, req.Targets)
		if err != nil {
			c.Logger().Errorf("failed to insert targets: %w", err)
			return err
		}

		err = q.DeleteAdministrators(ctx, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrators: %w", err)
			return err
		}

		err = q.InsertAdministrators(ctx, questionnaireID, req.Administrators)
		if err != nil {
			c.Logger().Errorf("failed to insert administrators: %w", err)
			return err
		}

		return nil
	})
	if err != nil {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			return httpError
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update a questionnaire")
	}

	return c.NoContent(http.StatusOK)
}

// DeleteQuestionnaire DELETE /questionnaires/:questionnaireID
func (q *Questionnaire) DeleteQuestionnaire(c echo.Context) error {
	questionnaireID, err := getQuestionnaireID(c)
	if err != nil {
		c.Logger().Errorf("failed to get questionnaireID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		err = q.IQuestionnaire.DeleteQuestionnaire(c.Request().Context(), questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete questionnaire: %w", err)
			return err
		}

		err = q.DeleteTargets(c.Request().Context(), questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete targets: %w", err)
			return err
		}

		err = q.DeleteAdministrators(c.Request().Context(), questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to delete administrators: %w", err)
			return err
		}

		return nil
	})
	if err != nil {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			return httpError
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete a questionnaire")
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

	allquestions, err := q.IQuestion.GetQuestions(c.Request().Context(), questionnaireID)
	if err != nil {
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

	options, err := q.GetOptions(c.Request().Context(), optionIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	optionMap := make(map[int][]string, len(options))
	for _, option := range options {
		optionMap[option.QuestionID] = append(optionMap[option.QuestionID], option.Body)
	}

	scaleLabels, err := q.GetScaleLabels(c.Request().Context(), scaleLabelIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, label := range scaleLabels {
		scaleLabelMap[label.QuestionID] = label
	}

	validations, err := q.GetValidations(c.Request().Context(), validationIDs)
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

func createQuestionnaireMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit null.Time, targets []string) string {
	var resTimeLimitText string
	if resTimeLimit.Valid {
		resTimeLimitText = resTimeLimit.Time.Local().Format("2006/01/02 15:04")
	} else {
		resTimeLimitText = "なし"
	}

	var targetsMentionText string
	if len(targets) == 0 {
		targetsMentionText = "なし"
	} else {
		targetsMentionText = "@" + strings.Join(targets, " @")
	}

	return fmt.Sprintf(
		`### アンケート『[%s](https://anke-to.trap.jp/questionnaires/%d)』が作成されました
#### 管理者
%s
#### 説明
%s
#### 回答期限
%s
#### 対象者
%s
#### 回答リンク
https://anke-to.trap.jp/responses/new/%d`,
		title,
		questionnaireID,
		strings.Join(administrators, ","),
		description,
		resTimeLimitText,
		targetsMentionText,
		questionnaireID,
	)
}
