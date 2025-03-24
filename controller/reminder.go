package controller

import (
	"context"
	"log"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
	// "golang.org/x/sync/semaphore"
)

type Job struct {
	Timestamp       time.Time
	QuestionnaireID int
	Action          func()
}

type JobQueue struct {
	jobs []*Job
	mu   sync.Mutex
}

var (
	// sem                   = semaphore.NewWeighted(1)
	Jq                    = &JobQueue{}
	Wg                    = &sync.WaitGroup{}
	reminderTimingMinutes = []int{5, 30, 60, 1440, 10080}
	reminderTimingStrings = []string{"5分", "30分", "1時間", "1日", "1週間"}
)

func (jq *JobQueue) Push(job *Job) {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	jq.jobs = append(jq.jobs, job)
	sort.Slice(jq.jobs, func(i, j int) bool {
		return jq.jobs[i].Timestamp.Before(jq.jobs[j].Timestamp)
	})
}

func (jq *JobQueue) Pop() *Job {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	if len(jq.jobs) == 0 {
		return nil
	}
	job := jq.jobs[0]
	jq.jobs = jq.jobs[1:]
	return job
}

func (jq *JobQueue) PushReminder(questionnaireID int, limit *time.Time) error {

	for i, timing := range reminderTimingMinutes {
		remindTimeStamp := limit.Add(-time.Duration(timing) * time.Minute)
		if remindTimeStamp.Before(time.Now()) {
			Jq.Push(&Job{
				Timestamp:       remindTimeStamp,
				QuestionnaireID: questionnaireID,
				Action: func() {
					err := reminderAction(questionnaireID, reminderTimingStrings[i])
					if err != nil {
						log.Printf("Failed to execute reminderAction for questionnaireID %d: %v", questionnaireID, err)
					}
				},
			})
		}
	}

	return nil
}

func (jq *JobQueue) DeleteReminder(questionnaireID int) error {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	if len(jq.jobs) == 1 && jq.jobs[0].QuestionnaireID == questionnaireID {
		jq.jobs = []*Job{}
	}
	for i, job := range jq.jobs {
		if job.QuestionnaireID == questionnaireID {
			jq.jobs = append(jq.jobs[:i], jq.jobs[i+1:]...)
		}
	}

	return nil
}

func (jq *JobQueue) CheckRemindStatus(questionnaireID int) (bool, error) {
	jq.mu.Lock()
	defer jq.mu.Unlock()
	for _, job := range jq.jobs {
		if job.QuestionnaireID == questionnaireID {
			return true, nil
		}
	}
	return false, nil
}

func reminderAction(questionnaireID int, leftTimeText string) error {
	ctx := context.Background()
	q := model.Questionnaire{}
	questionnaire, _, _, _, administrators, _, _, respondants, err := q.GetQuestionnaireInfo(ctx, questionnaireID)
	if err != nil {
		return err
	}

	var reminderTargets []string
	for _, target := range questionnaire.Targets {
		if target.IsCanceled {
			continue
		}
		if slices.Contains(respondants, target.UserTraqid) {
			continue
		}
		reminderTargets = append(reminderTargets, target.UserTraqid)
	}

	reminderMessage := createReminderMessage(questionnaireID, questionnaire.Title, questionnaire.Description, administrators, questionnaire.ResTimeLimit.Time, reminderTargets, leftTimeText)
	wh := traq.NewWebhook()
	err = wh.PostMessage(reminderMessage)
	if err != nil {
		return err
	}

	return nil
}

func ReminderWorker() {
	for {
		job := Jq.Pop()
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

func ReminderInit() {
	questionnaires, err := model.NewQuestionnaire().GetQuestionnairesInfoForReminder(context.Background())
	if err != nil {
		panic(err)
	}
	for _, questionnaire := range questionnaires {
		err := Jq.PushReminder(questionnaire.ID, &questionnaire.ResTimeLimit.Time)
		if err != nil {
			panic(err)
		}
	}
}
