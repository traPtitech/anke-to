package handler

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
func NewMiddleware() *Middleware {
	return &Middleware{}
}

const (
	validatorKey       = "validator"
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

/*
	消せないアンケートの発生を防ぐための管理者

暫定的にハードコーディングで対応
*/
var adminUserIDs = []string{"ryoha", "xxarupakaxx", "kaitoyama", "cp20", "itzmeowww"}

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
func (*Middleware) TrapRateLimitMiddlewareFunc() echo.MiddlewareFunc {
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

// QuestionnaireReadAuthenticate アンケートの閲覧権限があるかの認証
func (m *Middleware) QuestionnaireReadAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
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

		// 管理者ならOK
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
		if isAdmin {
			c.Set(questionnaireIDKey, questionnaireID)
			return next(c)
		}

		// 公開されたらOK
		questionnaire, _, _, _, _, _, err := m.GetQuestionnaireInfo(c.Request().Context(), questionnaireID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Infof("questionnaire not found: %+v", err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("questionnaire not found:%d", questionnaireID))
		}
		if err != nil {
			c.Logger().Errorf("failed to get questionnaire read privilege info: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get questionnaire read privilege info: %w", err))
		}
		if !questionnaire.IsPublished {
			return c.String(http.StatusForbidden, "The questionnaire is not published.")
		}

		c.Set(questionnaireIDKey, questionnaireID)

		return next(c)
	}
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

// getValidator Validatorを設定する
func getValidator(c echo.Context) (*validator.Validate, error) {
	rowValidate := c.Get(validatorKey)
	validate, ok := rowValidate.(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to get validator")
	}

	return validate, nil
}

// getUserID ユーザーIDを取得する
func getUserID(c echo.Context) (string, error) {
	rowUserID := c.Get(userIDKey)
	userID, ok := rowUserID.(string)
	if !ok {
		return "", errors.New("invalid context userID")
	}

	return userID, nil
}
