package tuning

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/traPtitech/anke-to/tuning/openapi"

	"golang.org/x/sync/errgroup"
)

func Inititial() {
	config := openapi.NewConfiguration()
	config.BasePath = "http://localhost:1323/api"
	client := openapi.NewAPIClient(config)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	newQuestionnaire := openapi.NewQuestionnaire{
		Title:          "後期プロジェクト/班所属募集",
		Description:    "後期プロジェクト/班所属のアンケートを実施します。\n所属希望の方は回答をお願いします。\n継続希望の方も回答は必須です。必ず回答するようにしてください。",
		ResTimeLimit:   time.Now().AddDate(0, 0, 7),
		ResSharedTo:    "public",
		Targets:        []string{"mds_boy"},
		Administrators: []string{"mds_boy"},
	}
	newQuestions := []openapi.NewQuestion{
		{
			PageNum:      1,
			QuestionNum:  0,
			QuestionType: "MultipleChoice",
			Body:         "複数のプロジェクトに所属を希望しますか(原則最大2つ迄)",
			IsRequired:   false,
			Options:      []string{"希望する", "希望しない"},
		},
		{
			PageNum:      1,
			QuestionNum:  1,
			QuestionType: "MultipleChoice",
			Body:         "第一希望のプロジェクトを選択してください",
			IsRequired:   false,
			Options:      []string{"[NEW] traPortation", "[NEW] Jump Jump Jump", "Presto Ray", "Clay Plate’s Story", "gameCreateTool", "JapariPark", "Hack and Slash", "神様のいない世界", "Arts（仮）", "Neo Showcase [新規メンバーの募集は行いません]", "traPortal v2 [新規メンバーの募集は行いません]", "anke-to v2"},
		},
		{
			PageNum:      1,
			QuestionNum:  2,
			QuestionType: "Text",
			Body:         "上記希望プロジェクトでの希望役職",
			IsRequired:   false,
		},
		{
			PageNum:      1,
			QuestionNum:  3,
			QuestionType: "MultipleChoice",
			Body:         "新規所属希望か継続所属希望かを回答してください",
			IsRequired:   false,
			Options:      []string{"新規所属希望", "継続所属希望"},
		},
		{
			PageNum:      1,
			QuestionNum:  4,
			QuestionType: "MultipleChoice",
			Body:         "第二希望のプロジェクトを選択してください",
			IsRequired:   false,
			Options:      []string{"[NEW] traPortation", "[NEW] Jump Jump Jump", "Presto Ray", "Clay Plate’s Story", "gameCreateTool", "JapariPark", "Hack and Slash", "神様のいない世界", "Arts（仮）", "Neo Showcase [新規メンバーの募集は行いません]", "traPortal v2 [新規メンバーの募集は行いません]", "anke-to v2"},
		},
		{
			PageNum:      1,
			QuestionNum:  5,
			QuestionType: "Text",
			Body:         "上記希望プロジェクトでの希望役職",
			IsRequired:   false,
		},
		{
			PageNum:      1,
			QuestionNum:  6,
			QuestionType: "MultipleChoice",
			Body:         "新規所属希望か継続所属希望かを回答してください",
			IsRequired:   false,
			Options:      []string{"新規所属希望", "継続所属希望"},
		},
		{
			PageNum:      1,
			QuestionNum:  7,
			QuestionType: "MultipleChoice",
			Body:         "第三希望のプロジェクトを選択してください",
			IsRequired:   false,
			Options:      []string{"[NEW] traPortation", "[NEW] Jump Jump Jump", "Presto Ray", "Clay Plate’s Story", "gameCreateTool", "JapariPark", "Hack and Slash", "神様のいない世界", "Arts（仮）", "Neo Showcase [新規メンバーの募集は行いません]", "traPortal v2 [新規メンバーの募集は行いません]", "anke-to v2"},
		},
		{
			PageNum:      1,
			QuestionNum:  8,
			QuestionType: "Text",
			Body:         "上記希望プロジェクトでの希望役職",
			IsRequired:   false,
		},
		{
			PageNum:      1,
			QuestionNum:  9,
			QuestionType: "MultipleChoice",
			Body:         "新規所属希望か継続所属希望かを回答してください",
			IsRequired:   false,
			Options:      []string{"新規所属希望", "継続所属希望"},
		},
		{
			PageNum:      1,
			QuestionNum:  10,
			QuestionType: "Checkbox",
			Body:         "所属を希望する班を選択してください(複数回答可)",
			IsRequired:   false,
			Options:      []string{"アルゴリズム班", "CTF班", "ゲーム班", "グラフィック班", "サウンド班", "SysAd班"},
		},
		{
			PageNum:      1,
			QuestionNum:  11,
			QuestionType: "TextArea",
			Body:         "自由記述欄",
			IsRequired:   false,
		},
	}
	newResponse := openapi.NewResponse{
		Body: []openapi.ResponseBody{
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"希望する"},
			},
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"[NEW] traPortation"},
			},
			{
				QuestionType: "Text",
				Response:     "にゃんこ",
			},
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"新規所属希望"},
			},
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"[NEW] traPortation"},
			},
			{
				QuestionType: "Text",
				Response:     "にゃんこ",
			},
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"新規所属希望"},
			},
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"[NEW] traPortation"},
			},
			{
				QuestionType: "Text",
				Response:     "にゃんこ",
			},
			{
				QuestionType:   "MultipleChoice",
				OptionResponse: []string{"新規所属希望"},
			},
			{
				QuestionType:   "Checkbox",
				OptionResponse: []string{"アルゴリズム班"},
			},
			{
				QuestionType: "TextArea",
				Response:     "がんばるぞい!",
			},
		},
	}

	routineNum := 3
	eg, ctx := errgroup.WithContext(ctx)
	reqFuncChan := make(chan func() error, 1)

	wg := sync.WaitGroup{}
	for i := 0; i < routineNum; i++ {
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
		L:
			for {
				select {
				case <-ctx.Done():
					return nil
				case reqFunc, ok := <-reqFuncChan:
					if !ok {
						break L
					}
					err := reqFunc()
					if err != nil {
						return fmt.Errorf("failed to send request: %w", err)
					}
				}
			}

			return nil
		})
	}
	wg.Wait()

	for i := 0; i < 750; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			questionnaireIDChan := make(chan int32, 1)
			reqFuncChan <- func(questinnaireIDChan chan int32) func() error {
				return func() error {
					questionnaireRes, _, err := client.QuestionnaireApi.PostQuestionnaire(ctx, newQuestionnaire)
					if err != nil {
						return fmt.Errorf("failed to make a questionnaire: %w", err)
					}

					questionnaireIDChan <- questionnaireRes.QuestionnaireID

					return nil
				}
			}(questionnaireIDChan)
			questionnaireID := <-questionnaireIDChan

			sm := sync.Map{}
			questionChan := make(chan struct{}, 1)
			reqFuncChan <- func(questionnaireID int32, sm *sync.Map, questionChan chan struct{}) func() error {
				return func() error {
					defer func() {
						questionChan <- struct{}{}
					}()
					for i, question := range newQuestions {
						question.QuestionnaireID = questionnaireID
						questionRes, _, err := client.QuestionApi.PostQuestion(ctx, question)
						if err != nil {
							return fmt.Errorf("failed to make question: %w", err)
						}

						sm.Store(i, questionRes.QuestionID)
					}

					return nil
				}
			}(questionnaireID, &sm, questionChan)
			<-questionChan

			responseNum := 10
			for i := 0; i < responseNum; i++ {
				reqFuncChan <- func(questionnaireID int32, sm *sync.Map) func() error {
					return func() error {
						response := newResponse
						response.QuestionnaireID = questionnaireID
						for j := range response.Body {
							iQuestionID, ok := sm.Load(j)
							if !ok {
								return errors.New("No questionID")
							}

							response.Body[j].QuestionID, ok = iQuestionID.(int32)
							if !ok {
								return errors.New("invalid questionID")
							}
						}

						_, _, err := client.ResponseApi.PostResponse(ctx, response)
						if err != nil {
							return fmt.Errorf("failed to make response: %w", err)
						}

						return nil
					}
				}(questionnaireID, &sm)
			}

			return nil
		})
	}

	go func() {
		wg.Wait()
		close(reqFuncChan)
	}()

	err := eg.Wait()
	if err != nil {
		panic(err)
	}
}

func Bench() {
	config := openapi.NewConfiguration()
	config.BasePath = "http://localhost:1323/api"
	client := openapi.NewAPIClient(config)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	questionnaireID := int32(1)
	questionDetails, _, err := client.QuestionnaireApi.GetQuestions(ctx, questionnaireID)
	if err != nil {
		panic(fmt.Errorf("failed to get questions: %w", err))
	}

	newResponse := openapi.NewResponse{
		QuestionnaireID: 1,
		Body: []openapi.ResponseBody{
			{
				QuestionID:   questionDetails[0].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "希望する",
			},
			{
				QuestionID:   questionDetails[1].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "[NEW] traPortation",
			},
			{
				QuestionID:   questionDetails[2].QuestionID,
				QuestionType: "Text",
				Response:     "にゃんこ",
			},
			{
				QuestionID:   questionDetails[3].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "新規所属希望",
			},
			{
				QuestionID:   questionDetails[4].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "[NEW] traPortation",
			},
			{
				QuestionID:   questionDetails[5].QuestionID,
				QuestionType: "Text",
				Response:     "にゃんこ",
			},
			{
				QuestionID:   questionDetails[6].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "新規所属希望",
			},
			{
				QuestionID:   questionDetails[7].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "[NEW] traPortation",
			},
			{
				QuestionID:   questionDetails[8].QuestionID,
				QuestionType: "Text",
				Response:     "にゃんこ",
			},
			{
				QuestionID:   questionDetails[9].QuestionID,
				QuestionType: "MultipleChoice",
				Response:     "新規所属希望",
			},
			{
				QuestionID:     questionDetails[10].QuestionID,
				QuestionType:   "Checkbox",
				OptionResponse: []string{"アルゴリズム班"},
			},
			{
				QuestionID:   questionDetails[11].QuestionID,
				QuestionType: "TextArea",
				Response:     "がんばるぞい!",
			},
		},
	}

	routineNum := 3
	ctx, cancel = context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < routineNum; i++ {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					_, _, err := client.ResponseApi.PostResponse(ctx, newResponse)
					if err != nil {
						if errors.Is(err, context.DeadlineExceeded) {
							return nil
						}

						return fmt.Errorf("failed to make a response: %w", err)
					}
				}
			}
		})
	}

	err = eg.Wait()
	if err != nil {
		panic(err)
	}
}
