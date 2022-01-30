package router

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/traPtitech/anke-to/router/session/mock_session"
	"github.com/traPtitech/anke-to/traq/mock_traq"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/model/mock_model"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type myResponse struct {
	Title           string    `json:"questionnaire_title"`
	ResTimeLimit    null.Time `json:"res_time_limit"`
	ResponseID      int       `json:"responseID"`
	QuestionnaireID int       `json:"questionnaireID"`
	ModifiedAt      time.Time `json:"modified_at"`
	SubmittedAt     null.Time `json:"submitted_at"`
	DeletedAt       null.Time `json:"deleted_at"`
}

type targettedQuestionnaire struct {
	QuestionnaireID int       `json:"questionnaireID"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ResTimeLimit    null.Time `json:"res_time_limit"`
	DeletedAt       null.Time `json:"deleted_at"`
	ResSharedTo     string    `json:"res_shared_to"`
	CreatedAt       time.Time `json:"created_at"`
	ModifiedAt      time.Time `json:"modified_at"`
	RespondedAt     null.Time `json:"responded_at"`
	HasResponse     bool      `json:"has_response"`
}

func TestGetTargettedQuestionnairesBytraQIDValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		request     *UserQueryparam
		isErr       bool
	}{
		{
			description: "一般的なQueryParameterなのでエラーなし",
			request: &UserQueryparam{
				Sort:     "created_at",
				Answered: "answered",
			},
		},
		{
			description: "Sortが-created_atでもエラーなし",
			request: &UserQueryparam{
				Sort:     "-created_at",
				Answered: "answered",
			},
		},
		{
			description: "Sortがtitleでもエラーなし",
			request: &UserQueryparam{
				Sort:     "title",
				Answered: "answered",
			},
		},
		{
			description: "Sortが-titleでもエラーなし",
			request: &UserQueryparam{
				Sort:     "-title",
				Answered: "answered",
			},
		},
		{
			description: "Sortがmodified_atでもエラーなし",
			request: &UserQueryparam{
				Sort:     "modified_at",
				Answered: "answered",
			},
		},
		{
			description: "Sortが-modified_atでもエラーなし",
			request: &UserQueryparam{
				Sort:     "-modified_at",
				Answered: "answered",
			},
		},
		{
			description: "Answeredがunansweredでもエラーなし",
			request: &UserQueryparam{
				Sort:     "created_at",
				Answered: "unanswered",
			},
		},
		{
			description: "Sortが空文字でもエラーなし",
			request: &UserQueryparam{
				Sort:     "",
				Answered: "answered",
			},
		},
		{
			description: "Answeredが空文字でもエラーなし",
			request: &UserQueryparam{
				Sort:     "created_at",
				Answered: "",
			},
		},
		{
			description: "Sortが指定された文字列ではないためエラー",
			request: &UserQueryparam{
				Sort:     "sort",
				Answered: "answered",
			},
			isErr: true,
		},
		{
			description: "Answeredが指定された文字列ではないためエラー",
			request: &UserQueryparam{
				Sort:     "created_at",
				Answered: "answer",
			},
			isErr: true,
		},
	}

	for _, test := range tests {
		validate := validator.New()

		t.Run(test.description, func(t *testing.T) {
			err := validate.Struct(test.request)
			if test.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUsersMe(t *testing.T) {

	type meResponseBody struct {
		TraqID string `json:"traqID"`
	}

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockStore := mock_session.NewMockIStore(ctrl)
	mockUser := mock_traq.NewMockIUser(ctrl)
	u := NewUser(
		mockRespondent,
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
	)
	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
		mockStore,
		mockUser,
	)

	type request struct {
		user users
	}
	type expect struct {
		isErr    bool
		code     int
		response meResponseBody
	}

	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user: userOne,
			},
			expect: expect{
				isErr: false,
				code:  http.StatusOK,
				response: meResponseBody{
					string(userOne),
				},
			},
		},
	}

	e := echo.New()
	e.GET("api/users/me", u.GetUsersMe, m.SetUserIDMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {
		rec := createRecorder(e, testCase.request.user, methodGet, makePath("/users/me"), typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		responseByte, jsonErr := json.Marshal(testCase.expect.response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

func TestGetMyResponses(t *testing.T) {

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	responseID1 := 1
	questionnaireID1 := 1
	responseID2 := 2
	questionnaireID2 := 2
	responseID3 := 3
	questionnaireID3 := 3
	myResponses := []myResponse{
		{
			ResponseID:      responseID1,
			QuestionnaireID: questionnaireID1,
			Title:           "質問1",
			ResTimeLimit:    null.NewTime(nowTime, false),
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
		},
		{
			ResponseID:      responseID2,
			QuestionnaireID: questionnaireID2,
			Title:           "質問2",
			ResTimeLimit:    null.NewTime(nowTime, false),
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
		},
		{
			ResponseID:      responseID3,
			QuestionnaireID: questionnaireID3,
			Title:           "質問3",
			ResTimeLimit:    null.NewTime(nowTime, false),
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
		},
	}
	respondentInfos := []model.RespondentInfo{
		{
			Title:        "質問1",
			ResTimeLimit: null.NewTime(nowTime, false),
			Respondents: model.Respondents{
				ResponseID:      responseID1,
				QuestionnaireID: questionnaireID1,
				SubmittedAt:     null.TimeFrom(nowTime),
				ModifiedAt:      nowTime,
			},
		},
		{
			Title:        "質問2",
			ResTimeLimit: null.NewTime(nowTime, false),
			Respondents: model.Respondents{
				ResponseID:      responseID2,
				QuestionnaireID: questionnaireID2,
				SubmittedAt:     null.TimeFrom(nowTime),
				ModifiedAt:      nowTime,
			},
		},
		{
			Title:        "質問3",
			ResTimeLimit: null.NewTime(nowTime, false),
			Respondents: model.Respondents{
				ResponseID:      responseID3,
				QuestionnaireID: questionnaireID3,
				SubmittedAt:     null.TimeFrom(nowTime),
				ModifiedAt:      nowTime,
			},
		},
	}

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockStore := mock_session.NewMockIStore(ctrl)
	mockUser := mock_traq.NewMockIUser(ctrl)
	u := NewUser(
		mockRespondent,
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
	)

	m := NewMiddleware(mockAdministrator, mockRespondent, mockQuestion, mockQuestionnaire, mockStore, mockUser)

	// Respondent
	// GetRespondentInfos
	// success
	mockRespondent.EXPECT().
		GetRespondentInfos(gomock.Any(), string(userOne)).
		Return(respondentInfos, nil).AnyTimes()
	// empty
	mockRespondent.EXPECT().
		GetRespondentInfos(gomock.Any(), "empty").
		Return([]model.RespondentInfo{}, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		GetRespondentInfos(gomock.Any(), "StatusInternalServerError").
		Return(nil, errMock).AnyTimes()

	type request struct {
		user users
	}
	type expect struct {
		isErr    bool
		code     int
		response []myResponse
	}

	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user: userOne,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: myResponses,
			},
		},
		{
			description: "empty",
			request: request{
				user: "empty",
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: []myResponse{},
			},
		},
		{
			description: "StatusInternalServerError",
			request: request{
				user: "StatusInternalServerError",
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
	}

	e := echo.New()
	e.GET("api/users/me/responses", u.GetMyResponses, m.SetUserIDMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {
		rec := createRecorder(e, testCase.request.user, methodGet, makePath("/users/me/responses"), typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		responseByte, jsonErr := json.Marshal(testCase.expect.response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

func TestGetMyResponsesByID(t *testing.T) {

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	responseID1 := 1
	responseID2 := 2
	questionnaireIDSuccess := 1
	questionnaireIDNotFound := -1
	myResponses := []myResponse{
		{
			ResponseID:      responseID1,
			QuestionnaireID: questionnaireIDSuccess,
			Title:           "質問1",
			ResTimeLimit:    null.NewTime(nowTime, false),
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
		},
		{
			ResponseID:      responseID2,
			QuestionnaireID: questionnaireIDSuccess,
			Title:           "質問2",
			ResTimeLimit:    null.NewTime(nowTime, false),
			SubmittedAt:     null.TimeFrom(nowTime),
			ModifiedAt:      nowTime,
		},
	}
	respondentInfos := []model.RespondentInfo{
		{
			Title:        "質問1",
			ResTimeLimit: null.NewTime(nowTime, false),
			Respondents: model.Respondents{
				ResponseID:      responseID1,
				QuestionnaireID: questionnaireIDSuccess,
				SubmittedAt:     null.TimeFrom(nowTime),
				ModifiedAt:      nowTime,
			},
		},
		{
			Title:        "質問2",
			ResTimeLimit: null.NewTime(nowTime, false),
			Respondents: model.Respondents{
				ResponseID:      responseID2,
				QuestionnaireID: questionnaireIDSuccess,
				SubmittedAt:     null.TimeFrom(nowTime),
				ModifiedAt:      nowTime,
			},
		},
	}

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockStore := mock_session.NewMockIStore(ctrl)
	mockUser := mock_traq.NewMockIUser(ctrl)
	m := NewMiddleware(mockAdministrator, mockRespondent, mockQuestion, mockQuestionnaire, mockStore, mockUser)
	u := NewUser(
		mockRespondent,
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
	)

	// Respondent
	// GetRespondentInfos
	// success
	mockRespondent.EXPECT().
		GetRespondentInfos(gomock.Any(), string(userOne), questionnaireIDSuccess).
		Return(respondentInfos, nil).AnyTimes()
	// questionnaireIDNotFound
	mockRespondent.EXPECT().
		GetRespondentInfos(gomock.Any(), string(userOne), questionnaireIDNotFound).
		Return([]model.RespondentInfo{}, nil).AnyTimes()
	// failure
	mockRespondent.EXPECT().
		GetRespondentInfos(gomock.Any(), "StatusInternalServerError", questionnaireIDSuccess).
		Return(nil, errMock).AnyTimes()

	type request struct {
		user            users
		questionnaireID int
		isBadParam      bool
	}
	type expect struct {
		isErr    bool
		code     int
		response []myResponse
	}

	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user:            userOne,
				questionnaireID: questionnaireIDSuccess,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: myResponses,
			},
		},
		{
			description: "questionnaireID does not exist",
			request: request{
				user:            userOne,
				questionnaireID: questionnaireIDNotFound,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: []myResponse{},
			},
		},
		{
			description: "StatusInternalServerError",
			request: request{
				user:            "StatusInternalServerError",
				questionnaireID: questionnaireIDSuccess,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
		{
			description: "badParam",
			request: request{
				user:            userOne,
				questionnaireID: questionnaireIDSuccess,
				isBadParam:      true,
			},
			expect: expect{
				isErr: true,
				code:  http.StatusBadRequest,
			},
		},
	}

	e := echo.New()
	e.GET("api/users/me/responses/:questionnaireID", u.GetMyResponsesByID, m.SetUserIDMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {
		reqPath := fmt.Sprint(rootPath, "/users/me/responses/", testCase.request.questionnaireID)
		if testCase.request.isBadParam {
			reqPath = fmt.Sprint(rootPath, "/users/me/responses/", "badParam")
		}
		rec := createRecorder(e, testCase.request.user, methodGet, reqPath, typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		responseByte, jsonErr := json.Marshal(testCase.expect.response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

func TestGetTargetedQuestionnaire(t *testing.T) {

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	questionnaireID1 := 1
	questionnaireID2 := 2
	targettedQuestionnaires := []model.TargettedQuestionnaire{
		{
			Questionnaires: model.Questionnaires{
				ID:           questionnaireID1,
				Title:        "questionnaireID1",
				Description:  "questionnaireID1",
				ResTimeLimit: null.TimeFrom(nowTime),
				DeletedAt: gorm.DeletedAt{
					Time:  nowTime,
					Valid: false,
				},
				ResSharedTo: "public",
				CreatedAt:   nowTime,
				ModifiedAt:  nowTime,
			},
			RespondedAt: null.NewTime(nowTime, false),
			HasResponse: false,
		},
		{
			Questionnaires: model.Questionnaires{
				ID:           questionnaireID2,
				Title:        "questionnaireID2",
				Description:  "questionnaireID2",
				ResTimeLimit: null.TimeFrom(nowTime),
				DeletedAt: gorm.DeletedAt{
					Time:  nowTime,
					Valid: false,
				},
				ResSharedTo: "public",
				CreatedAt:   nowTime,
				ModifiedAt:  nowTime,
			},
			RespondedAt: null.NewTime(nowTime, false),
			HasResponse: false,
		},
	}

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockStore := mock_session.NewMockIStore(ctrl)
	mockUser := mock_traq.NewMockIUser(ctrl)
	u := NewUser(
		mockRespondent,
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
	)
	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
		mockStore,
		mockUser,
	)

	// Questionnaire
	// GetTargettedQuestionnaires
	// success
	mockQuestionnaire.EXPECT().
		GetTargettedQuestionnaires(gomock.Any(), string(userOne), "", gomock.Any()).
		Return(targettedQuestionnaires, nil).AnyTimes()
	// empty
	mockQuestionnaire.EXPECT().
		GetTargettedQuestionnaires(gomock.Any(), "empty", "", gomock.Any()).
		Return([]model.TargettedQuestionnaire{}, nil).AnyTimes()
	// failure
	mockQuestionnaire.EXPECT().
		GetTargettedQuestionnaires(gomock.Any(), "StatusInternalServerError", "", gomock.Any()).
		Return(nil, errMock).AnyTimes()

	type request struct {
		user users
	}
	type expect struct {
		isErr    bool
		code     int
		response []model.TargettedQuestionnaire
	}

	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user: userOne,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: targettedQuestionnaires,
			},
		},
		{
			description: "empty",
			request: request{
				user: "empty",
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: []model.TargettedQuestionnaire{},
			},
		},
		{
			description: "StatusInternalServerError",
			request: request{
				user: "StatusInternalServerError",
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
	}

	e := echo.New()
	e.GET("api/users/me/targeted", u.GetTargetedQuestionnaire, m.SetUserIDMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {
		rec := createRecorder(e, testCase.request.user, methodGet, makePath("/users/me/targeted"), typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		responseByte, jsonErr := json.Marshal(testCase.expect.response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

func TestGetTargettedQuestionnairesBytraQID(t *testing.T) {

	t.Parallel()
	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nowTime := time.Now()

	questionnaireID1 := 1
	questionnaireID2 := 2
	targettedQuestionnaires := []model.TargettedQuestionnaire{
		{
			Questionnaires: model.Questionnaires{
				ID:           questionnaireID1,
				Title:        "questionnaireID1",
				Description:  "questionnaireID1",
				ResTimeLimit: null.TimeFrom(nowTime),
				DeletedAt: gorm.DeletedAt{
					Time:  nowTime,
					Valid: false,
				},
				ResSharedTo: "public",
				CreatedAt:   nowTime,
				ModifiedAt:  nowTime,
			},
			RespondedAt: null.NewTime(nowTime, false),
			HasResponse: false,
		},
		{
			Questionnaires: model.Questionnaires{
				ID:           questionnaireID2,
				Title:        "questionnaireID2",
				Description:  "questionnaireID2",
				ResTimeLimit: null.TimeFrom(nowTime),
				DeletedAt: gorm.DeletedAt{
					Time:  nowTime,
					Valid: false,
				},
				ResSharedTo: "public",
				CreatedAt:   nowTime,
				ModifiedAt:  nowTime,
			},
			RespondedAt: null.NewTime(nowTime, false),
			HasResponse: false,
		},
	}

	mockRespondent := mock_model.NewMockIRespondent(ctrl)
	mockQuestionnaire := mock_model.NewMockIQuestionnaire(ctrl)
	mockTarget := mock_model.NewMockITarget(ctrl)
	mockAdministrator := mock_model.NewMockIAdministrator(ctrl)

	mockQuestion := mock_model.NewMockIQuestion(ctrl)
	mockStore := mock_session.NewMockIStore(ctrl)
	mockUser := mock_traq.NewMockIUser(ctrl)
	u := NewUser(
		mockRespondent,
		mockQuestionnaire,
		mockTarget,
		mockAdministrator,
	)
	m := NewMiddleware(
		mockAdministrator,
		mockRespondent,
		mockQuestion,
		mockQuestionnaire,
		mockStore,
		mockUser,
	)

	// Questionnaire
	// GetTargettedQuestionnaires
	// success
	mockQuestionnaire.EXPECT().
		GetTargettedQuestionnaires(gomock.Any(), string(userOne), "", gomock.Any()).
		Return(targettedQuestionnaires, nil).AnyTimes()
	// empty
	mockQuestionnaire.EXPECT().
		GetTargettedQuestionnaires(gomock.Any(), "empty", "", gomock.Any()).
		Return([]model.TargettedQuestionnaire{}, nil).AnyTimes()
	// failure
	mockQuestionnaire.EXPECT().
		GetTargettedQuestionnaires(gomock.Any(), "StatusInternalServerError", "", gomock.Any()).
		Return(nil, errMock).AnyTimes()

	type request struct {
		user       users
		targetUser users
	}
	type expect struct {
		isErr    bool
		code     int
		response []model.TargettedQuestionnaire
	}

	type test struct {
		description string
		request
		expect
	}
	testCases := []test{
		{
			description: "success",
			request: request{
				user:       userOne,
				targetUser: userOne,
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: targettedQuestionnaires,
			},
		},
		{
			description: "empty",
			request: request{
				user:       userOne,
				targetUser: "empty",
			},
			expect: expect{
				isErr:    false,
				code:     http.StatusOK,
				response: []model.TargettedQuestionnaire{},
			},
		},
		{
			description: "StatusInternalServerError",
			request: request{
				user:       userOne,
				targetUser: "StatusInternalServerError",
			},
			expect: expect{
				isErr: true,
				code:  http.StatusInternalServerError,
			},
		},
	}

	e := echo.New()
	e.GET("api/users/:traQID/targeted", u.GetTargettedQuestionnairesBytraQID, m.SetUserIDMiddleware, m.SetValidatorMiddleware, m.TraPMemberAuthenticate)

	for _, testCase := range testCases {
		rec := createRecorder(e, testCase.request.user, methodGet, fmt.Sprint(rootPath, "/users/", testCase.request.targetUser, "/targeted"), typeNone, "")

		assertion.Equal(testCase.expect.code, rec.Code, testCase.description, "status code")
		if rec.Code < 200 || rec.Code >= 300 {
			continue
		}

		responseByte, jsonErr := json.Marshal(testCase.expect.response)
		require.NoError(t, jsonErr)
		responseStr := string(responseByte) + "\n"
		assertion.Equal(responseStr, rec.Body.String(), testCase.description, "responseBody")
	}
}

// func TestGetUsersMe(t *testing.T) {
// 	testList := []struct {
// 		description string
// 		result      meResponseBody
// 		expectCode  int
// 	}{}
// 	fmt.Println(testList)
// }

// func TestGetMyResponses(t *testing.T) {
// 	testList := []struct {
// 		description string
// 		result      respondentInfos
// 		expectCode  int
// 	}{}
// 	fmt.Println(testList)
// }

// func TestGetMyResponsesByID(t *testing.T) {
// 	testList := []struct {
// 		description     string
// 		questionnaireID int
// 		result          respondentInfos
// 		expectCode      int
// 	}{}
// 	fmt.Println(testList)
// }

// func TestGetTargetedQuestionnaire(t *testing.T) {
// 	testList := []struct {
// 		description string
// 		result      targettedQuestionnaire
// 		expectCode  int
// 	}{}
// 	fmt.Println(testList)
// }

// func TestGetMyQuestionnaire(t *testing.T) {
// 	testList := []struct {
// 		description string
// 		result      targettedQuestionnaire
// 		expectCode  int
// 	}{}
// 	fmt.Println(testList)
// }
// func TestGetTargettedQuestionnairesBytraQID(t *testing.T) {
// 	testList := []struct {
// 		description string
// 		result      targettedQuestionnaire
// 		expectCode  int
// 	}{}
// 	fmt.Println(testList)
// }
