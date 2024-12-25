package metrics

import (
	"fmt"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.uber.org/zap"
)

func TemporalHandler(addr string) (client.MetricsHandler, error) {
	scope, err := newPrometheusScope(prometheus.Configuration{
		ListenAddress: addr,
		TimerType:     "histogram",
	})
	if err != nil {
		return nil, fmt.Errorf("creating prometheus scope: %w", err)
	}
	return sdktally.NewMetricsHandler(scope), nil
}

func newPrometheusScope(c prometheus.Configuration) (tally.Scope, error) {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				zap.L().Error("error in prometheus reporter", zap.Error(err))
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating prometheus reporter: %w", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	return scope, nil
}
