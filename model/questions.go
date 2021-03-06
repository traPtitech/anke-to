//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

// IQuestion QuestionのRepository
type IQuestion interface {
	InsertQuestion(questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool) (int, error)
	UpdateQuestion(questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool, questionID int) error
	DeleteQuestion(questionID int) error
	GetQuestions(questionnaireID int) ([]Questions, error)
	CheckQuestionAdmin(userID string, questionID int) (bool, error)
}
