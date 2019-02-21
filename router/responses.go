package router

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"git.trapti.tech/SysAd/anke-to/model"
)

func PostResponse(c echo.Context) error {

	req := model.Responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	responseID, err := model.InsertRespondents(c, req)
	if err != nil {
		return nil
	}

	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				if err := model.InsertResponse(c, responseID, req, body, option); err != nil {
					return err
				}
			}
		default:
			if err := model.InsertResponse(c, responseID, req, body, body.Response); err != nil {
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

func GetMyResponses(c echo.Context) error {
	responsesinfo, err := model.GetMyResponses(c)
	if err != nil {
		return err
	}

	myresponses, err := model.GetResponsesInfo(c, responsesinfo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myresponses)
}

func GetMyResponsesByID(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	responsesinfo, err := model.GetMyResponsesByID(c, questionnaireID)
	if err != nil {
		return err
	}

	myresponses, err := model.GetResponsesInfo(c, responsesinfo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, myresponses)
}

func GetResponsesByID(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

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
		RespondedAt, err := model.RespondedAt(c, questionnaireID)
		if err != nil {
			return err
		}
		if !AmAdmin && RespondedAt == "NULL" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
	}

	sortQuery := c.QueryParam("sort")
	responsesinfo, sortNum, err := model.GetSortedResponses(c, questionnaireID, sortQuery)
	if err != nil {
		return err
	}

	type ResponsesInfo struct {
		ResponseID  int                  `json:"responseID"`
		UserID      string               `json:"traqID"`
		SubmittedAt string               `json:"submitted_at"`
		ModifiedAt  string               `json:"modified_at"`
		Body        []model.ResponseBody `json:"response_body"`
	}
	responses := []ResponsesInfo{}

	questionTypeList, err := model.GetQuestionsType(c, questionnaireID)
	if err != nil {
		return err
	}

	for _, response := range responsesinfo {
		bodyList := []model.ResponseBody{}
		for _, questionType := range questionTypeList {
			body, err := model.GetResponseBody(c, response.ResponseID, questionType.ID, questionType.Type)
			if err != nil {
				return err
			}
			bodyList = append(bodyList, body)
		}
		responses = append(responses,
			ResponsesInfo{
				ResponseID:  response.ResponseID,
				UserID:      response.UserID,
				SubmittedAt: model.NullTimeToString(response.SubmittedAt),
				ModifiedAt:  response.ModifiedAt.Format(time.RFC3339),
				Body:        bodyList,
			})
	}

	// 昇順
	if sortNum > 0 {
		sort.Slice(responses, func(i, j int) bool {
			return responses[i].Body[sortNum-1].Response < responses[j].Body[sortNum-1].Response
		})
	}
	// 降順
	if sortNum < 0 {
		sort.Slice(responses, func(i, j int) bool {
			return responses[i].Body[-sortNum-1].Response > responses[j].Body[-sortNum-1].Response
		})
	}

	return c.JSON(http.StatusOK, responses)
}

func GetResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	respondentInfo, err := model.GetRespondentByID(c, responseID)
	if err != nil {
		return nil
	}

	responses := struct {
		QuestionnaireID int                  `json:"questionnaireID"`
		SubmittedAt     string               `json:"submitted_at"`
		ModifiedAt      string               `json:"modified_at"`
		Body            []model.ResponseBody `json:"body"`
	}{
		respondentInfo.QuestionnaireID,
		model.NullTimeToString(respondentInfo.SubmittedAt),
		model.NullTimeToString(respondentInfo.ModifiedAt),
		[]model.ResponseBody{},
	}

	questionTypeList, err := model.GetQuestionsType(c, responses.QuestionnaireID)
	if err != nil {
		return err
	}

	bodyList := []model.ResponseBody{}
	for _, questionType := range questionTypeList {
		body, err := model.GetResponseBody(c, responseID, questionType.ID, questionType.Type)
		if err != nil {
			return err
		}
		bodyList = append(bodyList, body)
	}
	responses.Body = bodyList
	return c.JSON(http.StatusOK, responses)
}

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

	if err := model.UpdateRespondents(c, req.ID, responseID, req.SubmittedAt); err != nil {
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
				if err := model.InsertResponse(c, responseID, req, body, option); err != nil {
					return err
				}
			}
		default:
			if err := model.InsertResponse(c, responseID, req, body, body.Response); err != nil {
				return err
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

func DeleteResponse(c echo.Context) error {
	responseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := model.DeleteMyResponse(c, responseID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
