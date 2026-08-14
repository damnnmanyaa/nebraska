package metrics

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"

	"github.com/flatcar/nebraska/backend/pkg/api"
)

const defaultTestDbURL = "postgres://postgres:nebraska@127.0.0.1:5432/nebraska_tests?sslmode=disable&connect_timeout=10"

func TestMain(m *testing.M) {
	if os.Getenv("NEBRASKA_SKIP_TESTS") != "" {
		return
	}

	if _, ok := os.LookupEnv("NEBRASKA_DB_URL"); !ok {
		log.Printf("NEBRASKA_DB_URL not set, setting to default %q\n", defaultTestDbURL)
		_ = os.Setenv("NEBRASKA_DB_URL", defaultTestDbURL)
	}

	os.Exit(m.Run())
}

func ensureMetricsRegistered(t *testing.T) {
	t.Helper()
	for _, c := range []prometheus.Collector{
		appInstancePerChannelGaugeMetric,
		failedUpdatesGaugeMetric,
		openConnections,
		inUseConnections,
		idleConnections,
		instanceStatsLatestGaugeMetric,
	} {
		if err := prometheus.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				require.NoError(t, err)
			}
		}
	}
}

func gaugeValue(t *testing.T, g prometheus.Gauge) float64 {
	t.Helper()
	var m dto.Metric
	require.NoError(t, g.Write(&m))
	require.NotNil(t, m.Gauge)
	return m.GetGauge().GetValue()
}

func TestCalculateMetricsInstanceStatsLatest(t *testing.T) {
	ensureMetricsRegistered(t)

	a, err := api.NewForTest(api.OptionInitDB, api.OptionDisableUpdatesOnFailedRollout)
	require.NoError(t, err)
	defer a.Close()

	ts := time.Now().UTC()
	_, err = a.DB().Exec(
		`INSERT INTO instance_stats (timestamp, channel_name, arch, version, instances) VALUES ($1, $2, $3, $4, $5)`,
		ts, "stable", "AMD64", "4500.1.0", 128,
	)
	require.NoError(t, err)

	_, err = a.DB().Exec(
		`INSERT INTO instance_stats (timestamp, channel_name, arch, version, instances) VALUES ($1, $2, $3, $4, $5)`,
		ts, "beta", "ARM64", "4500.0.0", 12,
	)
	require.NoError(t, err)

	require.NoError(t, calculateMetrics(a))

	require.Equal(t, 128.0, gaugeValue(t, instanceStatsLatestGaugeMetric.WithLabelValues("stable", "AMD64", "4500.1.0")))
	require.Equal(t, 12.0, gaugeValue(t, instanceStatsLatestGaugeMetric.WithLabelValues("beta", "ARM64", "4500.0.0")))

	// Assert the Prometheus text exposition shape that /metrics serves.
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(rec, req)
	body := rec.Body.String()
	require.Contains(t, body, `nebraska_instance_stats_latest{arch="AMD64",channel_name="stable",version="4500.1.0"} 128`)
	require.Contains(t, body, `nebraska_instance_stats_latest{arch="ARM64",channel_name="beta",version="4500.0.0"} 12`)
	t.Logf("sample /metrics lines:\n%s", filterMetricLines(body, "nebraska_instance_stats_latest"))
}

func TestCalculateMetricsInstanceStatsLatestEmptyTable(t *testing.T) {
	ensureMetricsRegistered(t)

	a, err := api.NewForTest(api.OptionInitDB, api.OptionDisableUpdatesOnFailedRollout)
	require.NoError(t, err)
	defer a.Close()

	_, err = a.DB().Exec(`DELETE FROM instance_stats`)
	require.NoError(t, err)

	// Fresh install: empty instance_stats must not fail the ticker.
	require.NoError(t, calculateMetrics(a))
}

func filterMetricLines(body, name string) string {
	var out []string
	for _, line := range strings.Split(body, "\n") {
		if strings.Contains(line, name) {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
