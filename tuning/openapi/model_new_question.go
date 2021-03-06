/*
 * anke-to API
 *
 * anke-to API
 *
 * API version: 1.0.0-oas3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// NewQuestion struct for NewQuestion
type NewQuestion struct {
	QuestionnaireID int32 `json:"questionnaireID"`
	// アンケートの何ページ目の質問か
	PageNum int32 `json:"page_num"`
	// アンケートの質問のうち、何問目か
	QuestionNum int32 `json:"question_num"`
	// どのタイプの質問か (\"Text\", \"TextArea\", \"Number\", \"MultipleChoice\", \"Checkbox\", \"Dropdown\", \"LinearScale\", \"Date\", \"Time\")
	QuestionType string `json:"question_type"`
	Body         string `json:"body"`
	// 回答必須かどうか
	IsRequired      bool     `json:"is_required"`
	Options         []string `json:"options"`
	ScaleLabelRight string   `json:"scale_label_right"`
	ScaleLabelLeft  string   `json:"scale_label_left"`
	ScaleMin        int32    `json:"scale_min"`
	ScaleMax        int32    `json:"scale_max"`
	RegexPattern    string   `json:"regex_pattern"`
	MinBound        string   `json:"min_bound"`
	MaxBound        string   `json:"max_bound"`
}
