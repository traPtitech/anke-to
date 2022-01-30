package router

import (
	"errors"
	"fmt"
	"github.com/traPtitech/anke-to/router/session"
	"github.com/traPtitech/anke-to/traq"
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
	session.IStore
	traq.IUser
}

// NewMiddleware Middlewareのコンストラクタ
func NewMiddleware(IAdministrator model.IAdministrator, IRespondent model.IRespondent, IQuestion model.IQuestion, IQuestionnaire model.IQuestionnaire, IStore session.IStore, IUser traq.IUser) *Middleware {
	return &Middleware{
		IAdministrator: IAdministrator,
		IRespondent: IRespondent,
		IQuestion: IQuestion,
		IQuestionnaire: IQuestionnaire,
		IStore: IStore,
		IUser: IUser,
	}
}



const (
	validatorKey       = "validator"
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

func (m *Middleware) SetValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		c.Set(validatorKey, validate)

		return next(c)
	}
}

/* 消せないアンケートの発生を防ぐための管理者
暫定的にハードコーディングで対応*/
var adminUserIDs = []string{"temma", "sappi_red", "ryoha", "mazrean", "xxarupakaxx", "asari"}

func (m *Middleware) SessionMiddleware() echo.MiddlewareFunc {
	return m.IStore.GetMiddleware()
}

// SetUserIDMiddleware SessionからUserIDを取得
func (m *Middleware) SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := m.IStore.GetSession(c)
		if errors.Is(err, session.ErrNoSession) {
			return echo.NewHTTPError(http.StatusUnauthorized, "no session")
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userID, err := sess.GetUserID()
		if err != nil && !errors.Is(err, session.ErrNoValue) {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID :%w", err))
		}
		if errors.Is(err, session.ErrNoValue) {
			token, err := sess.GetToken()
			if errors.Is(err, session.ErrNoValue) {
				return echo.NewHTTPError(http.StatusUnauthorized, "no token")
			}
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get token :%w", err))
			}

			userID, err = m.IUser.GetMyID(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get UserID:%w", err))
			}
		}
		c.Set(userIDKey, userID)

		return next(c)
	}
}

// TraPMemberAuthenticate traP部員かの認証
func (m *Middleware) TraPMemberAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Errorf("failed to get userID: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		// トークンを持たないユーザはアクセスできない
		if userID == "-" {
			c.Logger().Info("not logged in")
			return echo.NewHTTPError(http.StatusUnauthorized, "You are not logged in")
		}

		return next(c)
	}
}

// TrapRateLimitMiddlewareFunc traP IDベースのリクエスト制限
func (m *Middleware) TrapRateLimitMiddlewareFunc() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(5),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			userID, err := getUserID(c)
			if err != nil {
				c.Logger().Errorf("failed to get userID: %+v", err)
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
			c.Logger().Errorf("failed to get userID: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
			c.Logger().Infof("failed to convert questionnaireID to int: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
		}

		for _, adminID := range adminUserIDs {
			if userID == adminID {
				c.Set(questionnaireIDKey, questionnaireID)

				return next(c)
			}
		}
		isAdmin, err := m.CheckQuestionnaireAdmin(c.Request().Context(), userID, questionnaireID)
		if err != nil {
			c.Logger().Errorf("failed to check questionnaire admin: %+v", err)
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
			c.Logger().Errorf("failed to get userID: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			c.Logger().Info("failed to convert responseID to int: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		// 回答者ならOK
		respondent, err := m.GetRespondent(c.Request().Context(), responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("response not found: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("response not found:%d", responseID))
		}
		if err != nil {
			c.Logger().Errorf("failed to check if you are a respondent: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are a respondent: %w", err))
		}
		if respondent == nil {
			c.Logger().Error("respondent is nil")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if respondent.UserTraqid == userID {
			return next(c)
		}

		// 回答者以外は一時保存の回答は閲覧できない
		if !respondent.SubmittedAt.Valid {
			c.Logger().Info("not submitted")

			// Note: 一時保存の回答の存在もわかってはいけないので、Respondentが見つからない時と全く同じエラーを返す
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("response not found:%d", responseID))
		}

		// アンケートごとの回答閲覧権限チェック
		responseReadPrivilegeInfo, err := m.GetResponseReadPrivilegeInfoByResponseID(c.Request().Context(), userID, responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("response not found: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid responseID: %d", responseID))
		} else if err != nil {
			c.Logger().Errorf("failed to get response read privilege info: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response read privilege info: %w", err))
		}

		haveReadPrivilege, err := checkResponseReadPrivilege(responseReadPrivilegeInfo)
		if err != nil {
			c.Logger().Errorf("failed to check response read privilege: %+v", err)
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
			c.Logger().Errorf("failed to get userID: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			c.Logger().Infof("failed to convert responseID to int: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		respondent, err := m.GetRespondent(c.Request().Context(), responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("response not found: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("response not found:%d", responseID))
		}
		if err != nil {
			c.Logger().Errorf("failed to check if you are a respondent: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are a respondent: %w", err))
		}
		if respondent == nil {
			c.Logger().Error("respondent is nil")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if respondent.UserTraqid != userID {
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
			c.Logger().Errorf("failed to get userID: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionID := c.Param("questionID")
		questionID, err := strconv.Atoi(strQuestionID)
		if err != nil {
			c.Logger().Infof("failed to convert questionID to int: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionID:%s(error: %w)", strQuestionID, err))
		}

		for _, adminID := range adminUserIDs {
			if userID == adminID {
				c.Set(questionIDKey, questionID)

				return next(c)
			}
		}
		isAdmin, err := m.CheckQuestionAdmin(c.Request().Context(), userID, questionID)
		if err != nil {
			c.Logger().Errorf("failed to check if you are a question administrator: %+v", err)
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
			c.Logger().Errorf("failed to get userID: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
			c.Logger().Infof("failed to convert questionnaireID to int: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
		}

		responseReadPrivilegeInfo, err := m.GetResponseReadPrivilegeInfoByQuestionnaireID(c.Request().Context(), userID, questionnaireID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("response not found: %+v", err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid responseID: %d", questionnaireID))
		} else if err != nil {
			c.Logger().Errorf("failed to get responseReadPrivilegeInfo: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response read privilege info: %w", err))
		}

		haveReadPrivilege, err := checkResponseReadPrivilege(responseReadPrivilegeInfo)
		if err != nil {
			c.Logger().Errorf("failed to check response read privilege: %+v", err)
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
	rowValidate := c.Get(validatorKey)
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
		return 0, errors.New("invalid context questionnaireID")
	}

	return questionnaireID, nil
}

func getResponseID(c echo.Context) (int, error) {
	rowResponseID := c.Get(responseIDKey)
	responseID, ok := rowResponseID.(int)
	if !ok {
		return 0, errors.New("invalid context responseID")
	}

	return responseID, nil
}

func getQuestionID(c echo.Context) (int, error) {
	rowQuestionID := c.Get(questionIDKey)
	questionID, ok := rowQuestionID.(int)
	if !ok {
		return 0, errors.New("invalid context questionID")
	}

	return questionID, nil
}
