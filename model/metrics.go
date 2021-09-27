package model

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
	gormPrometheus "gorm.io/plugin/prometheus"
)

type MetricsCollector struct {
	Prefix             string
	Interval           uint32
	questionnaireGauge *prometheus.GaugeVec
	questionGauge      *prometheus.GaugeVec
	respondentGauge    *prometheus.GaugeVec
	responseGauge      *prometheus.GaugeVec
	administratorGauge *prometheus.GaugeVec
}

func (mc *MetricsCollector) Metrics(p *gormPrometheus.Prometheus) []prometheus.Collector {
	if mc.Prefix == "" {
		mc.Prefix = "gorm_anke_to"
	}

	if mc.Interval == 0 {
		mc.Interval = p.RefreshInterval
	}

	if mc.questionnaireGauge == nil {
		mc.questionnaireGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: mc.Prefix,
			Subsystem: "questionnaire",
			Name:      "count",
			Help:      "Number of questionnaires",
		}, []string{"status"})
	}

	if mc.questionGauge == nil {
		mc.questionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: mc.Prefix,
			Subsystem: "question",
			Name:      "count",
			Help:      "Number of questions",
		}, []string{"status", "type", "required"})
	}

	if mc.respondentGauge == nil {
		mc.respondentGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: mc.Prefix,
			Subsystem: "respondent",
			Name:      "count",
			Help:      "Number of respondents",
		}, []string{"status"})
	}

	if mc.responseGauge == nil {
		mc.responseGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: mc.Prefix,
			Subsystem: "response",
			Name:      "count",
			Help:      "Number of responses",
		}, []string{"status"})
	}

	if mc.administratorGauge == nil {
		mc.administratorGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: mc.Prefix,
			Subsystem: "administrator",
			Name:      "count",
			Help:      "Number of administrators",
		}, []string{"traq_id"})
	}

	go func() {
		for range time.Tick(time.Duration(mc.Interval) * time.Second) {
			mc.collect(p)
		}
	}()

	mc.collect(p)

	return []prometheus.Collector{
		mc.questionnaireGauge,
		mc.questionGauge,
		mc.respondentGauge,
		mc.responseGauge,
		mc.administratorGauge,
	}
}

func (mc *MetricsCollector) collect(p *gormPrometheus.Prometheus) {
	ctx := context.Background()

	err := mc.collectQuestionnaireMetrics(ctx, p)
	if err != nil {
		p.DB.Logger.Error(ctx, "failed to collect questionnaire metrics: %v", err)
	}

	err = mc.collectQuestionMetrics(ctx, p)
	if err != nil {
		p.DB.Logger.Error(ctx, "failed to collect question metrics: %v", err)
	}

	err = mc.collectRespondentMetrics(ctx, p)
	if err != nil {
		p.DB.Logger.Error(ctx, "failed to collect respondent metrics: %v", err)
	}

	err = mc.collectResponseMetrics(ctx, p)
	if err != nil {
		p.DB.Logger.Error(ctx, "failed to collect response metrics: %v", err)
	}

	err = mc.collectAdministratorMetrics(ctx, p)
	if err != nil {
		p.DB.Logger.Error(ctx, "failed to collect administrator metrics: %v", err)
	}
}

func (mc *MetricsCollector) collectQuestionnaireMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	var questionnaireCounts []struct {
		IsDeleted bool  `gorm:"column:is_deleted"`
		Count     int64 `gorm:"column:count"`
	}

	err := p.DB.
		Session(&gorm.Session{
			NewDB:   true,
			Context: ctx,
		}).
		Unscoped().
		Model(&Questionnaires{}).
		Select("deleted_at IS NOT NULL AS is_deleted, count(*) as count").
		Group("is_deleted").
		Find(&questionnaireCounts).Error
	if err != nil {
		return fmt.Errorf("failed to get questionnaire count from db: %v", err)
	}

	mc.questionnaireGauge.Reset()
	for _, count := range questionnaireCounts {
		var label string
		if count.IsDeleted {
			label = "deleted"
		} else {
			label = "active"
		}

		mc.questionnaireGauge.
			WithLabelValues(label).
			Set(float64(count.Count))
	}

	return nil
}

func (mc *MetricsCollector) collectQuestionMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	var questionCounts []struct {
		IsDeleted bool   `gorm:"column:is_deleted"`
		Type      string `gorm:"column:type"`
		Required  bool   `gorm:"column:is_required"`
		Count     int64  `gorm:"column:count"`
	}

	err := p.DB.
		Session(&gorm.Session{
			NewDB:   true,
			Context: ctx,
		}).
		Unscoped().
		Model(&Questions{}).
		Select("deleted_at IS NOT NULL AS is_deleted, type, is_required, count(*) as count").
		Group("is_deleted, type, is_required").
		Find(&questionCounts).Error
	if err != nil {
		return fmt.Errorf("failed to get question count from db: %v", err)
	}

	mc.questionGauge.Reset()
	for _, count := range questionCounts {
		var label string
		if count.IsDeleted {
			label = "deleted"
		} else {
			label = "active"
		}

		var required string
		if count.Required {
			required = "required"
		} else {
			required = "optional"
		}

		mc.questionGauge.
			WithLabelValues(label, count.Type, required).
			Set(float64(count.Count))
	}

	return nil
}

func (mc *MetricsCollector) collectRespondentMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	var respondentCounts []struct {
		IsDeleted bool  `gorm:"column:is_deleted"`
		Count     int64 `gorm:"column:count"`
	}

	err := p.DB.
		Session(&gorm.Session{
			NewDB:   true,
			Context: ctx,
		}).
		Unscoped().
		Model(&Respondents{}).
		Select("deleted_at IS NOT NULL AS is_deleted, count(*) as count").
		Group("is_deleted").
		Find(&respondentCounts).Error
	if err != nil {
		return fmt.Errorf("failed to get respondent count from db: %v", err)
	}

	mc.respondentGauge.Reset()
	for _, count := range respondentCounts {
		var label string
		if count.IsDeleted {
			label = "deleted"
		} else {
			label = "active"
		}

		mc.respondentGauge.
			WithLabelValues(label).
			Set(float64(count.Count))
	}

	return nil
}

func (mc *MetricsCollector) collectResponseMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	var responseCounts []struct {
		IsDeleted bool  `gorm:"column:is_deleted"`
		Count     int64 `gorm:"column:count"`
	}

	err := p.DB.
		Session(&gorm.Session{
			NewDB:   true,
			Context: ctx,
		}).
		Unscoped().
		Model(&Responses{}).
		Select("deleted_at IS NOT NULL AS is_deleted, count(*) as count").
		Group("is_deleted").
		Find(&responseCounts).Error
	if err != nil {
		return fmt.Errorf("failed to get response count from db: %v", err)
	}

	mc.responseGauge.Reset()
	for _, count := range responseCounts {
		var label string
		if count.IsDeleted {
			label = "deleted"
		} else {
			label = "active"
		}

		mc.responseGauge.
			WithLabelValues(label).
			Set(float64(count.Count))
	}

	return nil
}

func (mc *MetricsCollector) collectAdministratorMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	var adminCounts []struct {
		UserTraqid string `gorm:"column:user_traqid"`
		Count      int64  `gorm:"column:count"`
	}

	err := p.DB.
		Session(&gorm.Session{
			NewDB:   true,
			Context: ctx,
		}).
		Unscoped().
		Model(&Administrator{}).
		Select("user_traqid, count(*) as count").
		Group("user_traqid").
		Find(&adminCounts).Error
	if err != nil {
		return fmt.Errorf("failed to get admin count from db: %v", err)
	}

	mc.administratorGauge.Reset()
	for _, count := range adminCounts {
		mc.administratorGauge.
			WithLabelValues(count.UserTraqid).
			Set(float64(count.Count))
	}

	return nil
}
