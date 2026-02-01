package checkers

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type Metrics struct {
	ExecutionsTotal *prometheus.CounterVec
	ExecutionStatus *prometheus.CounterVec
}

func NewMetrics(installation string) *Metrics {
	constLabels := prometheus.Labels{}
	if installation != "" {
		constLabels["installation"] = installation
	}

	return &Metrics{
		ExecutionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_checker_executions_total",
				Help:        "Total number of checker executions",
				ConstLabels: constLabels,
			},
			[]string{"action", "service_id"},
		),
		ExecutionStatus: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "fastad_checker_execution_status_total",
				Help:        "Checker execution results by status",
				ConstLabels: constLabels,
			},
			[]string{"action", "status", "service_id"},
		),
	}
}

func (m *Metrics) ObserveExecution(action checkerpb.Action, status checkerpb.Status, serviceID int) {
	serviceIDStr := strconv.Itoa(serviceID)
	actionStr := action.String()
	statusStr := status.String()

	m.ExecutionsTotal.WithLabelValues(actionStr, serviceIDStr).Inc()
	m.ExecutionStatus.WithLabelValues(actionStr, statusStr, serviceIDStr).Inc()
}
