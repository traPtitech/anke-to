package router

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type meResponseBody struct {
	traqID string
}

type respondentInfos struct {
	responnseID     int
	questionnaireID int
	title           string
	resTimeLimit    null.Time
	submittedAt     null.Time
	modifiedAt      time.Time
}

type targettedQuestionnaire struct {
	questionnaireID int
	title           string
	description     string
	resTimeLimit    null.Time
	createdAt       time.Time
	resSharedTo     string
	modifiedAt      time.Time
	respondedAt     null.Time
	hasResponse     bool
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
