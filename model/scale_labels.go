package model

import ()

type ScaleLabels struct {
	ID              int    `json:"questionID" db:"question_id"`
	ScaleLabelRight string `json:"scale_label_right" db:"scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"  db:"scale_label_left"`
	ScaleMin        int    `json:"scale_min" db:"scale_min"`
	ScaleMax        int    `json:"scale_max" db:"scale_max"`
}
