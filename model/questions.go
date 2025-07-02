//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "context"

// IQuestion QuestionのRepository
type IQuestion interface {
	InsertQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, title string, description string, isRequired bool) (int, error)
	UpdateQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, title string, description string, isRequired bool, questionID int) error
	DeleteQuestion(ctx context.Context, questionID int) error
	GetQuestions(ctx context.Context, questionnaireID int) ([]Questions, error)
	CheckQuestionAdmin(ctx context.Context, userID string, questionID int) (bool, error)
	CheckQuestionNum(ctx context.Context, questionnaireID, questionNum int) (bool, error)
}
