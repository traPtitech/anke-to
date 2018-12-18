package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func InsertResponse(c echo.Context, responseID int, req responses, body responseBody, data string) error {
	if req.SubmittedAt.Valid {
		if _, err := db.Exec(
			`INSERT INTO responses
				(questionnaire_id, question_id, response_id, body, user_traqid, submitted_at)
				VALUES (?, ?, ?, ?, ?, ?)`,
			req.ID, body.QuestionID, responseID, data, getUserID(c), req.SubmittedAt); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	} else {
		if _, err := db.Exec(
			`INSERT INTO responses
				(questionnaire_id, question_id, response_id, body, user_traqid)
				VALUES (?, ?, ?, ?, ?)`,
			req.ID, body.QuestionID, responseID, data, getUserID(c)); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	return nil
}

func postResponse(c echo.Context) error {

	req := responses{}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// responseIDを全部取って最大値+1をIDとする(もっといい方法がありそう)
	responseIDlist := []int{}

	if err := db.Select(&responseIDlist, "SELECT response_id from responses"); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	maxResponseID := 0
	for _, responseID := range responseIDlist {
		if maxResponseID < responseID {
			maxResponseID = responseID
		}
	}

	responseID := maxResponseID + 1

	for _, body := range req.Body {
		switch body.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, option := range body.OptionResponse {
				if err := InsertResponse(c, responseID, req, body, option); err != nil {
					return err
				}
			}
		default:
			if err := InsertResponse(c, responseID, req, body, body.Response); err != nil {
				return err
			}
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"questionnaireID": req.ID,
		"submitted_at":    timeConvert(req.SubmittedAt),
		"body":            req.Body,
	})
}

func getResponse(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func editResponse(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func deleteResponse(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
