package router

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
	"golang.org/x/sync/semaphore"
)

type Job struct {
	Timestamp       time.Time
	QuestionnaireID string
	Action          func()
}

type JobQueue struct {
	jobs []*Job
	mu   sync.Mutex
}

var (
	sem                    = semaphore.NewWeighted(1)
	Q                      = &JobQueue{}
	Wg                     = &sync.WaitGroup{}
	reiminderTimingMinutes = []int{5, 30, 60, 1440, 10080}
)

func (q *JobQueue) Push(j *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, j)
	sort.Slice(q.jobs, func(i, j int) bool {
		return q.jobs[i].Timestamp.Before(q.jobs[j].Timestamp)
	})
}

func (q *JobQueue) Pop() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.jobs) == 0 {
		return nil
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job
}

func (q *JobQueue) PushReminder(questionnaireID int) error {
	ctx := context.Background()
	questionnaire := model.Questionnaire{}
	limit, err := questionnaire.GetQuestionnaireLimit(ctx, questionnaireID)
	if err != nil {
		return err
	}
	if !limit.Valid {
		return nil
	}
	for _, timing := range reiminderTimingMinutes {
		remindTimeStamp := limit.Time.Add(-time.Duration(timing) * time.Minute)
		if remindTimeStamp.After(time.Now()) {
			Q.Push(&Job{
				Timestamp:       remindTimeStamp,
				QuestionnaireID: string(questionnaireID),
				Action: func() {
					reminderAction(questionnaireID, time.Until(limit.Time).String())
				},
			})

		}
	}
	return nil
}

func (q *JobQueue) DeleteReminder(questionnaireID int) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, job := range q.jobs {
		if job.QuestionnaireID == string(questionnaireID) {
			q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
		}
	}
	return nil
}

func reminderAction(questionnaireID int, lestTimeString string) error {
	ctx := context.Background()
	questionnaire := model.Questionnaire{}
	questionnaires, administorators, _, respondents, err := questionnaire.GetQuestionnaireInfo(ctx, questionnaireID)
	if err != nil {
		return err
	}

	var remindeTargets []string
	for _, target := range questionnaires.Targets {
		if !target.IsCanceled {
			if !slices.Contains(respondents, target.UserTraqid) {
				remindeTargets = append(remindeTargets, target.UserTraqid)
			}
		}
	}

	reminderMessage := createReminderMessage(questionnaireID, questionnaires.Title, questionnaires.Description, administorators, questionnaires.ResTimeLimit.Time, remindeTargets, lestTimeString)
	wh := traq.NewWebhook()
	err = wh.PostMessage(reminderMessage)
	if err != nil {
		return err
	}

	return nil
}

func ReminderWorker() {
	for {
		job := Q.Pop()
		if job == nil {
			time.Sleep(1 * time.Minute)
			continue
		}
		if time.Until(job.Timestamp) > 0 {
			time.Sleep(time.Until(job.Timestamp))
		}
		Wg.Add(1)
		go func() {
			defer Wg.Done()
			job.Action()
		}()
	}
}
