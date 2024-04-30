package router

import (
	"context"
	"log"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
	"golang.org/x/sync/semaphore"
	"gopkg.in/guregu/null.v4"
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
	sem                   = semaphore.NewWeighted(1)
	Q                     = &JobQueue{}
	Wg                    = &sync.WaitGroup{}
	reminderTimingMinutes = []int{5, 30, 60, 1440, 10080}
	reminderTimingStrings = []string{"5分", "30分", "1時間", "1日", "1週間"}
)

// jobQueueにjobを追加する
func (q *JobQueue) Push(j *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, j)
	sort.Slice(q.jobs, func(i, j int) bool {
		return q.jobs[i].Timestamp.Before(q.jobs[j].Timestamp)
	})
}

// jobQueueからjobを取り出す
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

// jobQueueにquestionnaireIDに対応するアンケートのリマインダーを追加する
func (q *JobQueue) PushReminder(questionnaireID int, limit null.Time) error {
	if !limit.Valid {
		return nil
	}
	log.Printf("[DEBUG] PushReminder: questionnaireID=%d, limit=%s\n", questionnaireID, limit.Time.String())
	for i, timing := range reminderTimingMinutes {
		remindTimeStamp := limit.Time.Add(-time.Duration(timing) * time.Minute)
		if remindTimeStamp.After(time.Now()) {
			Q.Push(&Job{
				Timestamp:       remindTimeStamp,
				QuestionnaireID: string(questionnaireID),
				Action: func() {
					reminderAction(questionnaireID, reminderTimingStrings[i])
				},
			})

		}
	}
	return nil
}

// jobQueueからquestionnaireIDに対応するアンケートのリマインダーを削除する
func (q *JobQueue) DeleteReminder(questionnaireID int) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, job := range q.jobs {
		if len(q.jobs) == 1 {
			q.jobs = []*Job{}
		} else {
			if job.QuestionnaireID == string(questionnaireID) {
				q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
			}
		}
	}
	return nil
}

// リマインダーのメッセージを送信する
func reminderAction(questionnaireID int, lestTimeString string) error {
	ctx := context.Background()
	questionnaire := model.Questionnaire{}
	questionnaires, _, administorators, respondents, err := questionnaire.GetQuestionnaireInfo(ctx, questionnaireID)
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
	log.Printf("[DEBUG] questionnaires.Targets=%v", questionnaires.Targets)
	log.Printf("[DEBUG] reminderAction: questionnaireID=%d, title=%s, description=%s, administorators=%v, resTimeLimit=%s, remindeTargets=%v, lestTimeString=%s\n", questionnaireID, questionnaires.Title, questionnaires.Description, administorators, questionnaires.ResTimeLimit.Time.String(), remindeTargets, lestTimeString)

	reminderMessage := createReminderMessage(questionnaireID, questionnaires.Title, questionnaires.Description, administorators, questionnaires.ResTimeLimit.Time, remindeTargets, lestTimeString)
	wh := traq.NewWebhook()
	err = wh.PostMessage(reminderMessage)
	if err != nil {
		return err
	}

	return nil
}

// リマインダーのメッセージを作成する
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
