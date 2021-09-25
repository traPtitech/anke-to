package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/anke-to/model"
)

// Middleware Middlewareの構造体
type Middleware struct {
	model.IAdministrator
	model.IRespondent
	model.IQuestion
	model.IQuestionnaire
}

// NewMiddleware Middlewareのコンストラクタ
func NewMiddleware(administrator model.IAdministrator, respondent model.IRespondent, question model.IQuestion, questionnaire model.IQuestionnaire) *Middleware {
	return &Middleware{
		IAdministrator: administrator,
		IRespondent:    respondent,
		IQuestion:      question,
		IQuestionnaire: questionnaire,
	}
}

const (
	validatorKay       = "validator"
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

func (*Middleware) SetValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		c.Set(validatorKay, validate)

		return next(c)
	}
}

/* 消せないアンケートの発生を防ぐための管理者
暫定的にハードコーディングで対応*/
var adminUserIDs = []string{"temma", "sappi_red", "ryoha", "mazrean", "YumizSui", "pure_white_404"}

// SetUserIDMiddleware X-Showcase-UserからユーザーIDを取得しセットする
func (*Middleware) SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Showcase-User")
		if userID == "" {
			userID = "mds_boy"
		}

		c.Set(userIDKey, userID)

		return next(c)
	}
}

// TraPMemberAuthenticate traP部員かの認証
func (*Middleware) TraPMemberAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		// トークンを持たないユーザはアクセスできない
		if userID == "-" {
			return echo.NewHTTPError(http.StatusUnauthorized, "You are not logged in")
		}

		return next(c)
	}
}

// TrapReteLimitMiddleware traP IDベースのリクエスト制限
func (*Middleware) TrapReteLimitMiddlewareFunc() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(5),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			userID, err := getUserID(c)
			if err != nil {
				c.Logger().Error(err)
				return "", echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
			}

			return userID, nil
		},
	}

	return middleware.RateLimiterWithConfig(config)
}

// QuestionnaireAdministratorAuthenticate アンケートの管理者かどうかの認証
func (m *Middleware) QuestionnaireAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
		}

		for _, adminID := range adminUserIDs {
			if userID == adminID {
				c.Set(questionnaireIDKey, questionnaireID)

				return next(c)
			}
		}
		isAdmin, err := m.CheckQuestionnaireAdmin(userID, questionnaireID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			return c.String(http.StatusForbidden, "You are not a administrator of this questionnaire.")
		}

		c.Set(questionnaireIDKey, questionnaireID)

		return next(c)
	}
}

// ResponseReadAuthenticate 回答閲覧権限があるかの認証
func (m *Middleware) ResponseReadAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		isRespondent, err := m.CheckRespondentByResponseID(userID, responseID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are a respondent: %w", err))
		}
		if isRespondent {
			return next(c)
		}

		responseReadPrivilegeInfo, err := m.GetResponseReadPrivilegeInfoByResponseID(userID, responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid responseID: %d", responseID))
		} else if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response read privilege info: %w", err))
		}

		haveReadPrivilege, err := checkResponseReadPrivilege(responseReadPrivilegeInfo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check response read privilege: %w", err))
		}
		if !haveReadPrivilege {
			return c.String(http.StatusForbidden, "You do not have permission to view this response.")
		}

		return next(c)
	}
}

// RespondentAuthenticate 回答者かどうかの認証
func (m *Middleware) RespondentAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		isRespondent, err := m.CheckRespondentByResponseID(userID, responseID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are a respondent: %w", err))
		}
		if !isRespondent {
			return c.String(http.StatusForbidden, "You are not a respondent of this response.")
		}

		c.Set(responseIDKey, responseID)

		return next(c)
	}
}

// QuestionAdministratorAuthenticate アンケートの管理者かどうかの認証
func (m *Middleware) QuestionAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionID := c.Param("questionID")
		questionID, err := strconv.Atoi(strQuestionID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionID:%s(error: %w)", strQuestionID, err))
		}

		for _, adminID := range adminUserIDs {
			if userID == adminID {
				c.Set(questionIDKey, questionID)

				return next(c)
			}
		}
		isAdmin, err := m.CheckQuestionAdmin(userID, questionID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			return c.String(http.StatusForbidden, "You are not a administrator of this questionnaire.")
		}

		c.Set(questionIDKey, questionID)

		return next(c)
	}
}

// ResultAuthenticate アンケートの回答を確認できるかの認証
func (m *Middleware) ResultAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
		}

		responseReadPrivilegeInfo, err := m.GetResponseReadPrivilegeInfoByQuestionnaireID(userID, questionnaireID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid responseID: %d", questionnaireID))
		} else if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response read privilege info: %w", err))
		}

		haveReadPrivilege, err := checkResponseReadPrivilege(responseReadPrivilegeInfo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check response read privilege: %w", err))
		}
		if !haveReadPrivilege {
			return c.String(http.StatusForbidden, "You do not have permission to view this response.")
		}

		return next(c)
	}
}

func checkResponseReadPrivilege(responseReadPrivilegeInfo *model.ResponseReadPrivilegeInfo) (bool, error) {
	switch responseReadPrivilegeInfo.ResSharedTo {
	case "administrators":
		return responseReadPrivilegeInfo.IsAdministrator, nil
	case "respondents":
		return responseReadPrivilegeInfo.IsAdministrator || responseReadPrivilegeInfo.IsRespondent, nil
	case "public":
		return true, nil
	}

	return false, errors.New("invalid resSharedTo")
}

func getValidator(c echo.Context) (*validator.Validate, error) {
	rowValidate := c.Get(validatorKay)
	validate, ok := rowValidate.(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to get validator")
	}

	return validate, nil
}

func getUserID(c echo.Context) (string, error) {
	rowUserID := c.Get(userIDKey)
	userID, ok := rowUserID.(string)
	if !ok {
		return "", errors.New("invalid context userID")
	}

	return userID, nil
}

func getQuestionnaireID(c echo.Context) (int, error) {
	rowQuestionnaireID := c.Get(questionnaireIDKey)
	questionnaireID, ok := rowQuestionnaireID.(int)
	if !ok {
		return 0, errors.New("invalid context userID")
	}

	return questionnaireID, nil
}

func getResponseID(c echo.Context) (int, error) {
	rowResponseID := c.Get(responseIDKey)
	questionnaireID, ok := rowResponseID.(int)
	if !ok {
		return 0, errors.New("invalid context userID")
	}

	return questionnaireID, nil
}

func getQuestionID(c echo.Context) (int, error) {
	rowQuestionID := c.Get(questionIDKey)
	questionID, ok := rowQuestionID.(int)
	if !ok {
		return 0, errors.New("invalid context userID")
	}

	return questionID, nil
}
