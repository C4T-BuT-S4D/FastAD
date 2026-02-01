package receiver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

type Metrics struct {
	FlagsSubmitted       *prometheus.CounterVec
	FlagsProcessed       *prometheus.CounterVec
	FlagVerdictsDetailed *prometheus.CounterVec
	AttacksByAttacker    *prometheus.CounterVec
	AttacksByVictim      *prometheus.CounterVec
	PointsByService      *prometheus.CounterVec
	ProcessingTime       *prometheus.HistogramVec
	SubmissionBatchSize  prometheus.Histogram
}

func NewMetrics(installation string) *Metrics {
	constLabels := prometheus.Labels{}
	if installation != "" {
		constLabels["installation"] = installation
	}

	return &Metrics{
		FlagsSubmitted: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_flags_submitted_total",
				Help:        "Total number of flags submitted",
				ConstLabels: constLabels,
			},
			[]string{"team_id", "team_name"},
		),
		FlagsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_flags_processed_total",
				Help:        "Total number of flags processed by verdict",
				ConstLabels: constLabels,
			},
			[]string{"verdict", "service_id", "service_name"},
		),
		FlagVerdictsDetailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_flag_verdicts_detailed_total",
				Help:        "Detailed flag verdicts with attacker, victim, service, and verdict",
				ConstLabels: constLabels,
			},
			[]string{"attacker_id", "attacker_name", "victim_id", "victim_name", "service_id", "service_name", "verdict"},
		),
		AttacksByAttacker: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_attacks_by_attacker_total",
				Help:        "Total successful attacks by attacker team",
				ConstLabels: constLabels,
			},
			[]string{"attacker_id", "attacker_name", "service_id", "service_name"},
		),
		AttacksByVictim: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_attacks_by_victim_total",
				Help:        "Total attacks received by victim team",
				ConstLabels: constLabels,
			},
			[]string{"victim_id", "victim_name", "service_id", "service_name"},
		),
		PointsByService: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_points_transferred_total",
				Help:        "Total points transferred by service",
				ConstLabels: constLabels,
			},
			[]string{"service_id", "service_name"},
		),
		ProcessingTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "fastad_flag_processing_seconds",
				Help:        "Time to process flag submission request",
				ConstLabels: constLabels,
				Buckets:     prometheus.ExponentialBucketsRange(0.0001, 5, 20),
			},
			[]string{"verdict_summary"},
		),
		SubmissionBatchSize: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:        "fastad_flag_submission_batch_size",
				Help:        "Number of flags per submission request",
				ConstLabels: constLabels,
				Buckets:     prometheus.ExponentialBuckets(1, 2, 8),
			},
		),
	}
}

func (m *Metrics) ObserveSubmission(teamID, teamName string, flagCount int) {
	m.FlagsSubmitted.WithLabelValues(teamID, teamName).Add(float64(flagCount))
	m.SubmissionBatchSize.Observe(float64(flagCount))
}

func (m *Metrics) ObserveVerdict(verdict receiverpb.FlagResponse_Verdict, serviceID, serviceName string) {
	m.FlagsProcessed.WithLabelValues(verdict.String(), serviceID, serviceName).Inc()
}

func (m *Metrics) ObserveVerdictDetailed(attackerID, attackerName, victimID, victimName, serviceID, serviceName string, verdict receiverpb.FlagResponse_Verdict) {
	m.FlagVerdictsDetailed.WithLabelValues(attackerID, attackerName, victimID, victimName, serviceID, serviceName, verdict.String()).Inc()
}

func (m *Metrics) ObserveAttack(attackerID, attackerName, victimID, victimName, serviceID, serviceName string, points float64) {
	m.AttacksByAttacker.WithLabelValues(attackerID, attackerName, serviceID, serviceName).Inc()
	m.AttacksByVictim.WithLabelValues(victimID, victimName, serviceID, serviceName).Inc()
	m.PointsByService.WithLabelValues(serviceID, serviceName).Add(points)
}

func (m *Metrics) ObserveProcessingTime(seconds float64, hasAccepted bool) {
	summary := "rejected"
	if hasAccepted {
		summary = "accepted"
	}
	m.ProcessingTime.WithLabelValues(summary).Observe(seconds)
}
