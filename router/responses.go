package router

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/traPtitech/anke-to/model"
)

// PostResponse POST /responses
func PostResponse(c echo.Context) error {

	req := model.Responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	limit, err := model.GetQuestionnaireLimit(c, req.ID)
	if err != nil {
		return err
	}

	// 回答期限を過ぎた回答は許可しない
	if limit != "NULL" && limit < time.Now().Format(time.RFC3339) {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	//パターンマッチ
	for _, body := range req.Body {
		validation, err := model.GetValidations(body.QuestionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		switch body.QuestionType {
		case "Number":
			if err := model.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.NumberValidError{}) {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Text":
			if err := model.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.TextMatchError{}) {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}

	responseID, err := model.InsertRespondent(c, req.ID, req.SubmittedAt)
	if err != nil {
		return err
	}

	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				if err := model.InsertResponse(c, responseID, body.QuestionID, option); err != nil {
					return err
				}
			}
		default:
			if err := model.InsertResponse(c, responseID, body.QuestionID, body.Body.ValueOrZero()); err != nil {
				return err
			}
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"responseID":      responseID,
		"questionnaireID": req.ID,
		"submitted_at":    req.SubmittedAt,
		"body":            req.Body,
	})
}

// GetMyResponses GET /users/me/responses
func GetMyResponses(c echo.Context) error {
	userID := model.GetUserID(c)
	myResponses, err := model.GetRespondentInfos(c, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myResponses)
}

// GetMyResponsesByID GET /users/me/responses/:questionnaireID
func GetMyResponsesByID(c echo.Context) error {
	userID := model.GetUserID(c)
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	myresponses, err := model.GetRespondentInfos(c, userID, questionnaireID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myresponses)
}

// GetResponsesByID GET /results/:questionnaireID
func GetResponsesByID(c echo.Context) error {
	sort := c.QueryParam("sort")
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// アンケートの回答を確認する権限が無ければエラーを返す
	if err := checkResponseConfirmable(c, questionnaireID); err != nil {
		return err
	}

	respondentDetails, err := model.GetRespondentDetails(c, questionnaireID, sort)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, respondentDetails)
}

// GetResponse GET /responses/:id
func GetResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	respondentDetail, err := model.GetRespondentDetail(c, responseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, respondentDetail)
}

// EditResponse PATCH /responses/:id
func EditResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	req := model.Responses{}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	limit, err := model.GetQuestionnaireLimit(c, req.ID)
	if err != nil {
		return err
	}

	// 回答期限を過ぎた回答は許可しない
	if limit != "NULL" && limit < time.Now().Format(time.RFC3339) {
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	//パターンマッチ
	for _, body := range req.Body {
		validation, err := model.GetValidations(body.QuestionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		switch body.QuestionType {
		case "Number":
			if err := model.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.NumberValidError{}) {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		case "Text":
			if err := model.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
				if errors.Is(err, &model.TextMatchError{}) {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}

	if err := model.UpdateRespondents(c, req.ID, responseID); err != nil {
		return err
	}

	//全消し&追加(レコード数爆発しそう)
	if err := model.DeleteResponse(c, responseID); err != nil {
		return err
	}

	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				if err := model.InsertResponse(c, responseID, body.QuestionID, option); err != nil {
					return err
				}
			}
		default:
			if err := model.InsertResponse(c, responseID, body.QuestionID, body.Body.ValueOrZero()); err != nil {
				return err
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

// DeleteResponse DELETE /responses/:id
func DeleteResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := model.DeleteRespondent(c, responseID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// アンケートの回答を確認できるか
func checkResponseConfirmable(c echo.Context, questionnaireID int) error {
	resSharedTo, err := model.GetResShared(c, questionnaireID)
	if err != nil {
		return err
	}

	switch resSharedTo {
	case "administrators":
		AmAdmin, err := model.CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	case "respondents":
		AmAdmin, err := model.CheckAdmin(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin {
			isRespondent, err := model.IsRespondent(c, questionnaireID)
			if err != nil {
				return err
			}
			if !isRespondent {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("only admins and respondents can see this responses"))
			}
		}
	}
	return nil
}
