package router

// API api全体の構造体
type API struct {
	*Middleware
	*Questionnaire
	*Question
	*Response
	*Result
	*User
	*Oauth
}

// NewAPI APIのコンストラクタ
func NewAPI(middleware *Middleware, questionnaire *Questionnaire, question *Question, response *Response, result *Result, user *User, oauth *Oauth) *API {
	return &API{
		Middleware:    middleware,
		Questionnaire: questionnaire,
		Question:      question,
		Response:      response,
		Result:        result,
		User:          user,
		Oauth:         oauth,
	}
}
