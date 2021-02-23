package router

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type resultResponseBody struct {
	questionnaireID int
	body            []responseBody
	userTraqid      string
	submittedAt     null.Time
	modifiedAt      time.Time
}

// func TestGetResults(t *testing.T) {
// 	testList := []struct {
// 		description     string
// 		questionnaireID int
// 		result          resultResponseBody
// 		expectCode      int
// 	}{
// 		{
// 			description:     "valid",
// 			questionnaireID: -1,
// 			result: resultResponseBody{
// 				questionnaireID: -1,
// 				submittedAt:     null.TimeFrom(time.Now()),
// 				userTraqid:      string(userOne),
// 				body: []responseBody{
// 					{
// 						QuestionID:     -1,
// 						QuestionType:   "Text",
// 						Body:           null.StringFrom("回答"),
// 						OptionResponse: []string{},
// 					},
// 				},
// 			},
// 			expectCode: http.StatusOK,
// 		},
// 		{
// 			description:     "unauthorized",
// 			questionnaireID: -1,
// 			result: resultResponseBody{
// 				questionnaireID: -1,
// 				submittedAt:     null.TimeFrom(time.Now()),
// 				userTraqid:      string(userTwo),
// 				body: []responseBody{
// 					{
// 						QuestionID:     -1,
// 						QuestionType:   "Text",
// 						ody:           null.StringFrom("回答"),
// 						optionResponse: []string{},
// 					},
// 				},
// 			},
// 			expectCode: http.StatusUnauthorized,
// 		},
// 		{
// 			description:     "not admin",
// 			questionnaireID: -1,
// 			result: resultResponseBody{
// 				questionnaireID: -1,
// 				submittedAt:     null.TimeFrom(time.Now()),
// 				userTraqid:      string(userTwo),
// 				body: []responseBody{
// 					{
// 						questionID:     -1,
// 						questionType:   "Text",
// 						body:           null.StringFrom("回答"),
// 						optionResponse: []string{},
// 					},
// 				},
// 			},
// 			expectCode: http.StatusUnauthorized,
// 		},
// 		{
// 			description:     "not respondents",
// 			questionnaireID: -1,
// 			result: resultResponseBody{
// 				questionnaireID: -1,
// 				submittedAt:     null.TimeFrom(time.Now()),
// 				userTraqid:      string(userTwo),
// 				body: []responseBody{
// 					{
// 						questionID:     -1,
// 						questionType:   "Text",
// 						body:           null.StringFrom("回答"),
// 						optionResponse: []string{},
// 					},
// 				},
// 			},
// 			expectCode: http.StatusUnauthorized,
// 		},
// 	}
// 	fmt.Println(testList)
// }
