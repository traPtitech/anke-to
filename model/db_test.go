package model

import (
	"os"
	"testing"
)

const (
	userOne   = "mazrean"
	userTwo   = "ryoha"
	userThree = "YumizSui"
)

var (
	administratorImpl = new(Administrator)
	questionnaireImpl = new(Questionnaire)
	questionImpl      = new(Question)
	respondentImpl    = new(Respondent)
	responseImpl      = new(Response)
	scaleLabelImpl    = new(ScaleLabel)
	validationImpl    = new(Validation)
)

//TestMain テストのmain
func TestMain(m *testing.M) {
	db, err := EstablishConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = Migrate()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}
