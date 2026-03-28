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

type Reminder struct {
	jobs     []*Job
	mu       sync.Mutex
	Wg       sync.WaitGroup
	wakeUpCh chan struct{}
}

func NewReminder() *Reminder {
	return &Reminder{
		jobs:     []*Job{},
		mu:       sync.Mutex{},
		Wg:       sync.WaitGroup{},
		wakeUpCh: make(chan struct{}, 1),
	}
}

var (
	// sem                   = semaphore.NewWeighted(1)
	reminderTimingMinutes = []int{5, 30, 60, 1440, 10080}
	reminderTimingStrings = []string{"5分", "30分", "1時間", "1日", "1週間"}
)

func (re *Reminder) ReminderInit() {
	questionnaires, err := model.NewQuestionnaire().GetQuestionnairesInfoForReminder(context.Background())
	if err != nil {
		panic(err)
	}
	for _, questionnaire := range questionnaires {
		err := re.PushReminder(questionnaire.ID, &questionnaire.ResTimeLimit.Time)
		if err != nil {
			panic(err)
		}
	}
}

func (re *Reminder) ReminderWorker() {
	for {
		job := re.peek()
		if job == nil {
			<-re.wakeUpCh
			continue
		}

		wait := time.Until(job.Timestamp)
		if wait > 0 {
			timer := time.NewTimer(wait)
			select {
			case <-timer.C:
			case <-re.wakeUpCh:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				continue
			}
		}

		job = re.popDue(time.Now())
		if job == nil {
			continue
		}

		re.Wg.Add(1)
		go func() {
			defer re.Wg.Done()
			job.Action()
		}()
	}
}

func (re *Reminder) PushReminder(questionnaireID int, limit *time.Time) error {

	for i := range reminderTimingMinutes {
		timing := reminderTimingMinutes[i]
		timingStrings := reminderTimingStrings[i]
		remindTimeStamp := limit.Add(-time.Duration(timing) * time.Minute)
		if remindTimeStamp.After(time.Now()) {
			re.push(&Job{
				Timestamp:       remindTimeStamp,
				QuestionnaireID: questionnaireID,
				Action: func() {
					err := reminderAction(questionnaireID, timingStrings)
					if err != nil {
						log.Printf("Failed to execute reminderAction for questionnaireID %d: %v", questionnaireID, err)
					}
				},
			})
		}
	}

	return nil
}

func (re *Reminder) DeleteReminder(questionnaireID int) error {
	re.mu.Lock()

	newJobs := []*Job{}
	removed := false
	for _, job := range re.jobs {
		if job.QuestionnaireID != questionnaireID {
			newJobs = append(newJobs, job)
			continue
		}
		removed = true
	}
	re.jobs = newJobs
	re.mu.Unlock()

	if removed {
		re.notifyWorker()
	}

	return nil
}

func (re *Reminder) CheckRemindStatus(questionnaireID int) (bool, error) {
	re.mu.Lock()
	defer re.mu.Unlock()
	for _, job := range re.jobs {
		if job.QuestionnaireID == questionnaireID {
			return true, nil
		}
	}
	return false, nil
}

func (re *Reminder) push(job *Job) {
	re.mu.Lock()
	re.jobs = append(re.jobs, job)
	sort.Slice(re.jobs, func(i, j int) bool {
		return re.jobs[i].Timestamp.Before(re.jobs[j].Timestamp)
	})
	re.mu.Unlock()

	re.notifyWorker()
}

func (re *Reminder) peek() *Job {
	re.mu.Lock()
	defer re.mu.Unlock()
	if len(re.jobs) == 0 {
		return nil
	}
	return re.jobs[0]
}

func (re *Reminder) popDue(now time.Time) *Job {
	re.mu.Lock()
	defer re.mu.Unlock()
	if len(re.jobs) == 0 {
		return nil
	}
	if re.jobs[0].Timestamp.After(now) {
		return nil
	}
	job := re.jobs[0]
	re.jobs = re.jobs[1:]
	return job
}

func (re *Reminder) notifyWorker() {
	select {
	case re.wakeUpCh <- struct{}{}:
	default:
	}
}

func reminderAction(questionnaireID int, leftTimeText string) error {
	ctx := context.Background()
	questionnaire, _, _, _, administrators, _, _, respondants, err := model.NewQuestionnaire().GetQuestionnaireInfo(ctx, questionnaireID)
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
