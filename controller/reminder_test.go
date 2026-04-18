package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/anke-to/openapi"
)

// nthJob returns the nth job (0-indexed) from the btree in ascending order.
func nthJob(re *Reminder, n int) *Job {
	var result *Job
	i := 0
	re.tree.Ascend(func(item *Job) bool {
		if i == n {
			result = item
			return false
		}
		i++
		return true
	})
	return result
}

func TestPushReminder(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionnaireID int
		time            time.Time
	}
	type expect struct {
		num   int
		isErr bool
		err   error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "5 reminders",
			args: args{
				questionnaireID: 1,
				time:            time.Now().Add(8 * 24 * time.Hour),
			},
			expect: expect{
				num: 5,
			},
		},
		{
			description: "4 reminders",
			args: args{
				questionnaireID: 1,
				time:            time.Now().Add(25 * time.Hour),
			},
			expect: expect{
				num: 4,
			},
		},
		{
			description: "3 reminders",
			args: args{
				questionnaireID: 1,
				time:            time.Now().Add(2 * time.Hour),
			},
			expect: expect{
				num: 3,
			},
		},
		{
			description: "2 reminders",
			args: args{
				questionnaireID: 1,
				time:            time.Now().Add(50 * time.Minute),
			},
			expect: expect{
				num: 2,
			},
		},
		{
			description: "1 reminders",
			args: args{
				questionnaireID: 1,
				time:            time.Now().Add(10 * time.Minute),
			},
			expect: expect{
				num: 1,
			},
		},
		{
			description: "0 reminders",
			args: args{
				questionnaireID: 1,
				time:            time.Now().Add(time.Minute),
			},
			expect: expect{
				num: 0,
			},
		},
	}

	for _, testCase := range testCases {
		re := NewReminder()
		err := re.PushReminder(testCase.args.questionnaireID, &testCase.args.time)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}
		assertion.Equal(testCase.expect.num, re.tree.Len(), "reminder num")
	}
}

func TestDeleteReminder(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionnaireID int
	}
	type expect struct {
		num   int
		isErr bool
		err   error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "delete multiple reminders",
			args: args{
				questionnaireID: 1,
			},
			expect: expect{
				num: 3,
			},
		},
		{
			description: "delete one reminder",
			args: args{
				questionnaireID: 2,
			},
			expect: expect{
				num: 1,
			},
		},
		{
			description: "delete no reminder",
			args: args{
				questionnaireID: 100,
			},
			expect: expect{
				num: 0,
			},
		},
	}

	for _, testCase := range testCases {
		re := NewReminder()
		reminderLimit1 := time.Now().Add(2 * time.Hour)
		reminderLimit2 := time.Now().Add(10 * time.Minute)
		reminderLimit3 := time.Now().Add(25 * time.Hour)
		err := re.PushReminder(1, &reminderLimit1)
		require.NoError(t, err)
		err = re.PushReminder(2, &reminderLimit2)
		require.NoError(t, err)
		err = re.PushReminder(3, &reminderLimit3)
		require.NoError(t, err)
		jobsNum := re.tree.Len()
		err = re.DeleteReminder(testCase.args.questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}
		assertion.Equal(jobsNum-testCase.expect.num, re.tree.Len(), testCase.description, "reminder num")
	}
}
func TestCheckRemindStatus(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	type args struct {
		questionnaireID int
	}
	type expect struct {
		status bool
		isErr  bool
		err    error
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "questionnaire with multiple reminders",
			args: args{
				questionnaireID: 1,
			},
			expect: expect{
				status: true,
			},
		},
		{
			description: "questionnaire with  one reminder",
			args: args{
				questionnaireID: 2,
			},
			expect: expect{
				status: true,
			},
		},
		{
			description: "questionnaire with  no reminder",
			args: args{
				questionnaireID: 100,
			},
			expect: expect{
				status: false,
			},
		},
	}

	for _, testCase := range testCases {
		re := NewReminder()
		reminderLimit1 := time.Now().Add(2 * time.Hour)
		reminderLimit2 := time.Now().Add(10 * time.Minute)
		reminderLimit3 := time.Now().Add(25 * time.Hour)
		err := re.PushReminder(1, &reminderLimit1)
		require.NoError(t, err)
		err = re.PushReminder(2, &reminderLimit2)
		require.NoError(t, err)
		err = re.PushReminder(3, &reminderLimit3)
		require.NoError(t, err)
		status, err := re.CheckRemindStatus(testCase.args.questionnaireID)
		if !testCase.expect.isErr {
			assertion.NoError(err, testCase.description, "no error")
		} else if testCase.expect.err != nil {
			assertion.Equal(true, errors.Is(err, testCase.expect.err), testCase.description, "errorIs")
		} else {
			assertion.Error(err, testCase.description, "any error")
		}
		if err != nil {
			continue
		}
		assertion.Equal(testCase.expect.status, status, "reminder status")
	}
}

func TestReminderWorkerReschedulesEarlierJob(t *testing.T) {
	re := NewReminder()

	now := time.Now()
	laterExecutedAtCh := make(chan time.Time, 1)
	earlierExecutedAtCh := make(chan time.Time, 1)

	re.push(&Job{
		Timestamp:       now.Add(500 * time.Millisecond),
		QuestionnaireID: 1,
		Action: func() {
			laterExecutedAtCh <- time.Now()
		},
	})

	go re.ReminderWorker()

	time.Sleep(100 * time.Millisecond)

	earlierDue := now.Add(150 * time.Millisecond)
	re.push(&Job{
		Timestamp:       earlierDue,
		QuestionnaireID: 2,
		Action: func() {
			earlierExecutedAtCh <- time.Now()
		},
	})

	select {
	case executedAt := <-earlierExecutedAtCh:
		assert.WithinDuration(t, earlierDue, executedAt, 150*time.Millisecond)
	case <-time.After(1 * time.Second):
		t.Fatal("earlier reminder was not executed in time")
	}

	select {
	case laterExecutedAt := <-laterExecutedAtCh:
		assert.WithinDuration(t, now.Add(500*time.Millisecond), laterExecutedAt, 150*time.Millisecond)
	case <-time.After(1 * time.Second):
		t.Fatal("later reminder was not executed in time")
	}
}

func TestPush(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	re := NewReminder()

	type args struct {
		job *Job
	}
	type expect struct {
		position int
	}
	type test struct {
		description string
		args
		expect
	}

	testCases := []test{
		{
			description: "push to empty queue",
			args: args{
				job: &Job{
					Timestamp:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					QuestionnaireID: 1,
					Action:          func() {},
				},
			},
			expect: expect{
				position: 0,
			},
		},
		{
			description: "push to queue with one job",
			args: args{
				job: &Job{
					Timestamp:       time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					QuestionnaireID: 2,
					Action:          func() {},
				},
			},
			expect: expect{
				position: 1,
			},
		},
		{
			description: "push to queue and sort",
			args: args{
				job: &Job{
					Timestamp:       time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					QuestionnaireID: 3,
					Action:          func() {},
				},
			},
			expect: expect{
				position: 1,
			},
		},
	}

	for i, testCase := range testCases {
		re.push(testCase.args.job)
		assertion.Equal(i, re.tree.Len()-1, "queue length")
		got := nthJob(re, testCase.expect.position)
		assertion.Equal(testCase.args.job.Timestamp, got.Timestamp, testCase.description, "pushed position timestamp")
		assertion.Equal(testCase.args.job.QuestionnaireID, got.QuestionnaireID, testCase.description, "pushed position questionnaire id")
	}
}

func TestPop(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	re := NewReminder()

	jobs := []*Job{
		{
			Timestamp:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			QuestionnaireID: 1,
			Action:          func() {},
		},
		{
			Timestamp:       time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			QuestionnaireID: 1,
			Action:          func() {},
		},
		{
			Timestamp:       time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			QuestionnaireID: 1,
			Action:          func() {},
		},
	}

	type expect struct {
		num int
	}
	type test struct {
		description string
		expect
	}

	for _, job := range jobs {
		re.push(job)
	}

	testCases := []test{
		{
			description: "pop queue with 3 jobs",
			expect: expect{
				num: 2,
			},
		},
		{
			description: "pop queue with 2 jobs",
			expect: expect{
				num: 1,
			},
		},
		{
			description: "pop queue with one jobs",
			expect: expect{
				num: 0,
			},
		},
		{
			description: "pop queue with no jobs",
			expect: expect{
				num: 0,
			},
		},
	}

	for _, testCase := range testCases {
		re.popDue(time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC))
		assertion.Equal(testCase.expect.num, re.tree.Len(), testCase.description, "queue length")
		if testCase.expect.num != 0 {
			earliest, ok := re.tree.Min()
			assertion.True(ok)
			assertion.Equal(jobs[3-testCase.expect.num].Timestamp, earliest.Timestamp, testCase.description, "first content timestamp")
			assertion.Equal(jobs[3-testCase.expect.num].QuestionnaireID, earliest.QuestionnaireID, testCase.description, "first content questionnaire id")
		}
	}
}

func TestReminderActionUnpublished(t *testing.T) {
	responseDueDateTimePlus := time.Now().Add(24 * time.Hour)
	params := openapi.PostQuestionnaireJSONRequestBody{
		Admin:                    sampleAdmin,
		Description:              "リマインダーテスト用アンケート",
		IsDuplicateAnswerAllowed: false,
		IsAnonymous:              false,
		IsPublished:              false,
		Questions:                []openapi.NewQuestion{},
		ResponseDueDateTime:      &responseDueDateTimePlus,
		ResponseViewableBy:       "anyone",
		Target:                   sampleTarget,
		Title:                    "未公開リマインダーテスト",
	}

	e := echo.New()
	body, err := json.Marshal(params)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/questionnaires", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, rec)

	detail, err := q.PostQuestionnaire(ctx, params)
	require.NoError(t, err)

	err = reminderAction(detail.QuestionnaireId, "5分")
	assert.NoError(t, err, "reminderAction should return nil for unpublished questionnaire")
}
