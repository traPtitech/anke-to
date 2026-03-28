package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func TestCleanupSoftDeletedQuestionnaires(t *testing.T) {
	now := time.Now()

	deletedQuestionnaire := Questionnaires{
		Title:        "cleanup target questionnaire",
		Description:  "cleanup target questionnaire",
		ResTimeLimit: null.Time{},
		ResSharedTo:  "administrators",
		IsPublished:  true,
		DeletedAt: gorm.DeletedAt{
			Time:  now,
			Valid: true,
		},
	}
	activeQuestionnaire := Questionnaires{
		Title:        "cleanup active questionnaire",
		Description:  "cleanup active questionnaire",
		ResTimeLimit: null.Time{},
		ResSharedTo:  "administrators",
		IsPublished:  true,
	}

	err := db.Create(&deletedQuestionnaire).Error
	if err != nil {
		t.Fatalf("failed to create deleted questionnaire: %v", err)
	}
	err = db.Create(&activeQuestionnaire).Error
	if err != nil {
		t.Fatalf("failed to create active questionnaire: %v", err)
	}

	deletedQuestion := Questions{
		QuestionnaireID: deletedQuestionnaire.ID,
		PageNum:         1,
		QuestionNum:     1,
		Type:            "Text",
		Body:            "deleted question",
		Description:     "deleted question",
	}
	activeQuestion := Questions{
		QuestionnaireID: activeQuestionnaire.ID,
		PageNum:         1,
		QuestionNum:     1,
		Type:            "Text",
		Body:            "active question",
		Description:     "active question",
	}

	err = db.Create(&deletedQuestion).Error
	if err != nil {
		t.Fatalf("failed to create deleted question: %v", err)
	}
	err = db.Create(&activeQuestion).Error
	if err != nil {
		t.Fatalf("failed to create active question: %v", err)
	}

	deletedRespondent := Respondents{
		QuestionnaireID: deletedQuestionnaire.ID,
		UserTraqid:      "deleted-user",
		SubmittedAt:     null.NewTime(now, true),
	}
	activeRespondent := Respondents{
		QuestionnaireID: activeQuestionnaire.ID,
		UserTraqid:      "active-user",
		SubmittedAt:     null.NewTime(now, true),
	}

	err = db.Create(&deletedRespondent).Error
	if err != nil {
		t.Fatalf("failed to create deleted respondent: %v", err)
	}
	err = db.Create(&activeRespondent).Error
	if err != nil {
		t.Fatalf("failed to create active respondent: %v", err)
	}

	for _, record := range []Targets{
		{QuestionnaireID: deletedQuestionnaire.ID, UserTraqid: "deleted-target"},
		{QuestionnaireID: activeQuestionnaire.ID, UserTraqid: "active-target"},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create target: %v", err)
		}
	}

	for _, record := range []TargetUsers{
		{QuestionnaireID: deletedQuestionnaire.ID, UserTraqid: "deleted-target-user"},
		{QuestionnaireID: activeQuestionnaire.ID, UserTraqid: "active-target-user"},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create target user: %v", err)
		}
	}

	for _, record := range []TargetGroups{
		{QuestionnaireID: deletedQuestionnaire.ID, GroupID: groupOne},
		{QuestionnaireID: activeQuestionnaire.ID, GroupID: groupTwo},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create target group: %v", err)
		}
	}

	for _, record := range []Administrators{
		{QuestionnaireID: deletedQuestionnaire.ID, UserTraqid: "deleted-admin"},
		{QuestionnaireID: activeQuestionnaire.ID, UserTraqid: "active-admin"},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create administrator: %v", err)
		}
	}

	for _, record := range []AdministratorUsers{
		{QuestionnaireID: deletedQuestionnaire.ID, UserTraqid: "deleted-admin-user"},
		{QuestionnaireID: activeQuestionnaire.ID, UserTraqid: "active-admin-user"},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create administrator user: %v", err)
		}
	}

	for _, record := range []AdministratorGroups{
		{QuestionnaireID: deletedQuestionnaire.ID, GroupID: groupOne},
		{QuestionnaireID: activeQuestionnaire.ID, GroupID: groupTwo},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create administrator group: %v", err)
		}
	}

	for _, record := range []Validations{
		{QuestionID: deletedQuestion.ID, RegexPattern: "^deleted$"},
		{QuestionID: activeQuestion.ID, RegexPattern: "^active$"},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create validation: %v", err)
		}
	}

	for _, record := range []ScaleLabels{
		{QuestionID: deletedQuestion.ID, ScaleLabelLeft: "low", ScaleLabelRight: "high", ScaleMin: 1, ScaleMax: 5},
		{QuestionID: activeQuestion.ID, ScaleLabelLeft: "low", ScaleLabelRight: "high", ScaleMin: 1, ScaleMax: 5},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create scale label: %v", err)
		}
	}

	for _, record := range []Options{
		{QuestionID: deletedQuestion.ID, OptionNum: 1, Body: "deleted option"},
		{QuestionID: activeQuestion.ID, OptionNum: 1, Body: "active option"},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create option: %v", err)
		}
	}

	for _, record := range []Responses{
		{ResponseID: deletedRespondent.ResponseID, QuestionID: deletedQuestion.ID, Body: null.NewString("deleted response", true)},
		{ResponseID: activeRespondent.ResponseID, QuestionID: activeQuestion.ID, Body: null.NewString("active response", true)},
	} {
		err = db.Create(&record).Error
		if err != nil {
			t.Fatalf("failed to create response: %v", err)
		}
	}

	err = cleanupSoftDeletedQuestionnaires(db)
	if err != nil {
		t.Fatalf("failed to cleanup soft-deleted questionnaires: %v", err)
	}

	assertRecordCount(t, "questionnaires", "id = ?", deletedQuestionnaire.ID, 1)
	assertRecordCount(t, "questionnaires", "id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "question", "id = ?", deletedQuestion.ID, 0)
	assertRecordCount(t, "question", "id = ?", activeQuestion.ID, 1)
	assertRecordCount(t, "respondents", "response_id = ?", deletedRespondent.ResponseID, 0)
	assertRecordCount(t, "respondents", "response_id = ?", activeRespondent.ResponseID, 1)
	assertRecordCount(t, "responses", "response_id = ?", deletedRespondent.ResponseID, 0)
	assertRecordCount(t, "responses", "response_id = ?", activeRespondent.ResponseID, 1)
	assertRecordCount(t, "targets", "questionnaire_id = ?", deletedQuestionnaire.ID, 0)
	assertRecordCount(t, "targets", "questionnaire_id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "target_users", "questionnaire_id = ?", deletedQuestionnaire.ID, 0)
	assertRecordCount(t, "target_users", "questionnaire_id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "target_groups", "questionnaire_id = ?", deletedQuestionnaire.ID, 0)
	assertRecordCount(t, "target_groups", "questionnaire_id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "administrators", "questionnaire_id = ?", deletedQuestionnaire.ID, 0)
	assertRecordCount(t, "administrators", "questionnaire_id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "administrator_users", "questionnaire_id = ?", deletedQuestionnaire.ID, 0)
	assertRecordCount(t, "administrator_users", "questionnaire_id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "administrator_groups", "questionnaire_id = ?", deletedQuestionnaire.ID, 0)
	assertRecordCount(t, "administrator_groups", "questionnaire_id = ?", activeQuestionnaire.ID, 1)
	assertRecordCount(t, "validations", "question_id = ?", deletedQuestion.ID, 0)
	assertRecordCount(t, "validations", "question_id = ?", activeQuestion.ID, 1)
	assertRecordCount(t, "scale_labels", "question_id = ?", deletedQuestion.ID, 0)
	assertRecordCount(t, "scale_labels", "question_id = ?", activeQuestion.ID, 1)
	assertRecordCount(t, "options", "question_id = ?", deletedQuestion.ID, 0)
	assertRecordCount(t, "options", "question_id = ?", activeQuestion.ID, 1)
}

func assertRecordCount(t *testing.T, table string, query string, arg interface{}, expected int64) {
	t.Helper()

	var count int64
	err := db.Table(table).Where(query, arg).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count %s: %v", table, err)
	}

	assert.Equal(t, expected, count, table)
}
