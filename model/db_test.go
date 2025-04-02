package model

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

const (
	userOne   = "mazrean"
	userTwo   = "ryoha"
	userThree = "YumizSui"
)

var (
	groupOne   = uuid.MustParse("e0ac793f-8c4b-441f-a6d4-a3719e01e99b")
	groupTwo   = uuid.MustParse("d39fb117-8a61-437b-860d-6216bdbec197")
	groupThree = uuid.MustParse("025794c3-2a1c-4647-947c-36e2a150fc47")
)

var (
	administratorImpl      = new(Administrator)
	administratorUserImpl  = new(AdministratorUser)
	administratorGroupImpl = new(AdministratorGroup)
	questionnaireImpl      = new(Questionnaire)
	questionImpl           = new(Question)
	respondentImpl         = new(Respondent)
	responseImpl           = new(Response)
	optionImpl             = new(Option)
	scaleLabelImpl         = new(ScaleLabel)
	validationImpl         = new(Validation)
	targetImpl             = new(Target)
	targetUserImpl         = new(TargetUser)
	targetGroupImpl        = new(TargetGroup)
)

// TestMain テストのmain
func TestMain(m *testing.M) {
	err := EstablishConnection(true)
	if err != nil {
		panic(err)
	}

	_, err = Migrate()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}
