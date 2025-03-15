package controller

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
		assertion.Equal(testCase.expect.num, len(re.jobs), "reminder num")
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
		re.PushReminder(1, &reminderLimit1)
		re.PushReminder(2, &reminderLimit2)
		re.PushReminder(3, &reminderLimit3)
		jobsNum := len(re.jobs)
		err := re.DeleteReminder(testCase.args.questionnaireID)
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
		assertion.Equal(jobsNum-testCase.expect.num, len(re.jobs), testCase.description, "reminder num")
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
		re.PushReminder(1, &reminderLimit1)
		re.PushReminder(2, &reminderLimit2)
		re.PushReminder(3, &reminderLimit3)
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
		assertion.Equal(i, len(re.jobs)-1, "queue length")
		assertion.Equal(testCase.args.job.Timestamp, re.jobs[testCase.expect.position].Timestamp, testCase.description, "pushed position timestamp")
		assertion.Equal(testCase.args.job.QuestionnaireID, re.jobs[testCase.expect.position].QuestionnaireID, testCase.description, "pushed position questionnaire id")
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
		re.pop()
		assertion.Equal(testCase.expect.num, len(re.jobs), testCase.description, "queue length")
		if testCase.expect.num != 0 {
			assertion.Equal(jobs[3-testCase.expect.num].Timestamp, re.jobs[0].Timestamp, testCase.description, "first content timestamp")
			assertion.Equal(jobs[3-testCase.expect.num].QuestionnaireID, re.jobs[0].QuestionnaireID, testCase.description, "first content questionnaire id")
		}
	}
}
