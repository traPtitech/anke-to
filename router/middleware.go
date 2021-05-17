package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router/session"
	"github.com/traPtitech/anke-to/traq"
)

// Middleware Middlewareの構造体
type Middleware struct {
	model.IAdministrator
	model.IRespondent
	model.IQuestion
	session.IStore
	traq.IUser
}

// NewMiddleware Middlewareのコンストラクタ
func NewMiddleware(administrator model.IAdministrator, respondent model.IRespondent, question model.IQuestion, session session.IStore, user traq.IUser) *Middleware {
	return &Middleware{
		IAdministrator: administrator,
		IRespondent:    respondent,
		IQuestion:      question,
		IStore:  session,
		IUser:          user,
	}
}

const (
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

/* 消せないアンケートの発生を防ぐための管理者
暫定的にハードコーディングで対応*/
var adminUserIDs = []string{"temma", "sappi_red", "ryoha", "mazrean", "YumizSui", "pure_white_404"}

func (m *Middleware) SessionMiddleware() echo.MiddlewareFunc {
	return m.IStore.GetMiddleware()
}

// UserAuthenticate traPのメンバーかの認証
func (m *Middleware) UserAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := m.IStore.GetSession(c)
		if errors.Is(err, session.ErrNoSession) {
			return echo.NewHTTPError(http.StatusUnauthorized, "no session")
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get session: %w", err))
		}

		userID, err := sess.GetUserID()
		if err != nil && !errors.Is(err, session.ErrNoValue) {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}
		if errors.Is(err, session.ErrNoValue) {
			token, err := sess.GetToken()
			if errors.Is(err, session.ErrNoValue) {
				return echo.NewHTTPError(http.StatusUnauthorized, "no token")
			}
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get token: %w", err))
			}

			userID, err = m.IUser.GetMyID(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
			}
		}

		c.Set(userIDKey, userID)

		return next(c)
	}
}

// QuestionnaireAdministratorAuthenticate アンケートの管理者かどうかの認証
func (m *Middleware) QuestionnaireAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			return c.String(http.StatusForbidden, "You are not a administrator of this questionnaire.")
		}

		c.Set(questionnaireIDKey, questionnaireID)

		return next(c)
	}
}

// RespondentAuthenticate 回答者かどうかの認証
func (m *Middleware) RespondentAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		isRespondent, err := m.CheckRespondentByResponseID(userID, responseID)
		if err != nil {
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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionID := c.Param("questionID")
		questionID, err := strconv.Atoi(strQuestionID)
		if err != nil {
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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			return c.String(http.StatusForbidden, "You are not a administrator of this questionnaire.")
		}

		c.Set(questionIDKey, questionID)

		return next(c)
	}
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
