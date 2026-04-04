package controller

import (
	"context"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/google/btree"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
)

type Job struct {
	Timestamp       time.Time
	QuestionnaireID int
	TimingIndex     int
	Action          func()
}

type Reminder struct {
	tree     *btree.BTreeG[*Job]
	index    map[int][]*Job
	mu       sync.Mutex
	Wg       sync.WaitGroup
	wakeUpCh chan struct{}
}

func jobLess(a, b *Job) bool {
	if !a.Timestamp.Equal(b.Timestamp) {
		return a.Timestamp.Before(b.Timestamp)
	}
	if a.QuestionnaireID != b.QuestionnaireID {
		return a.QuestionnaireID < b.QuestionnaireID
	}
	return a.TimingIndex < b.TimingIndex
}

func NewReminder() *Reminder {
	return &Reminder{
		tree:     btree.NewG(32, jobLess),
		index:    make(map[int][]*Job),
		mu:       sync.Mutex{},
		Wg:       sync.WaitGroup{},
		wakeUpCh: make(chan struct{}, 1),
	}
}

var (
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
				TimingIndex:     i,
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
	jobs, exists := re.index[questionnaireID]
	if exists {
		for _, job := range jobs {
			re.tree.Delete(job)
		}
		delete(re.index, questionnaireID)
	}
	re.mu.Unlock()

	if exists {
		re.notifyWorker()
	}

	return nil
}

func (re *Reminder) CheckRemindStatus(questionnaireID int) (bool, error) {
	re.mu.Lock()
	defer re.mu.Unlock()
	jobs, exists := re.index[questionnaireID]
	return exists && len(jobs) > 0, nil
}

func (re *Reminder) push(job *Job) {
	re.mu.Lock()
	re.tree.ReplaceOrInsert(job)
	re.index[job.QuestionnaireID] = append(re.index[job.QuestionnaireID], job)
	re.mu.Unlock()

	re.notifyWorker()
}

func (re *Reminder) peek() *Job {
	re.mu.Lock()
	defer re.mu.Unlock()
	earliest, ok := re.tree.Min()
	if !ok {
		return nil
	}
	return earliest
}

func (re *Reminder) popDue(now time.Time) *Job {
	re.mu.Lock()
	defer re.mu.Unlock()
	earliest, ok := re.tree.Min()
	if !ok {
		return nil
	}
	if earliest.Timestamp.After(now) {
		return nil
	}
	re.tree.DeleteMin()
	jobs := re.index[earliest.QuestionnaireID]
	for i, j := range jobs {
		if j == earliest {
			re.index[earliest.QuestionnaireID] = append(jobs[:i], jobs[i+1:]...)
			break
		}
	}
	if len(re.index[earliest.QuestionnaireID]) == 0 {
		delete(re.index, earliest.QuestionnaireID)
	}
	return earliest
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

	reminderMessages := createReminderMessage(questionnaireID, questionnaire.Title, questionnaire.Description, administrators, questionnaire.ResTimeLimit.Time, reminderTargets, leftTimeText)
	wh := traq.NewWebhook()
	for _, msg := range reminderMessages {
		err = wh.PostMessage(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
