package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flatcar/nebraska/backend/pkg/api"
	"github.com/flatcar/nebraska/backend/pkg/codegen"
)

func TestGetInstanceStatsLatest(t *testing.T) {
	// 1. Get real API with DB init
	a, err := api.NewForTest(api.OptionInitDB, api.OptionDisableUpdatesOnFailedRollout)
	require.NoError(t, err)
	defer a.Close()

	// 2. Seed data directly via DB handle
	db := a.DB()
	ts := time.Now().UTC()
	_, err = db.Exec(`INSERT INTO instance_stats (timestamp, channel_name, arch, version, instances) VALUES ($1, $2, $3, $4, $5)`, ts, "channel", "AMD64", "1.0.0", 5)
	require.NoError(t, err)

	// 3. Construct minimal handler
	h := &Handler{
		db: a,
	}

	// 4. Build echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/instances_stats/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 5. Call handler directly
	err = h.GetInstanceStatsLatest(c)
	assert.NoError(t, err)

	// 6. Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	var stats []codegen.InstanceStats
	err = json.Unmarshal(rec.Body.Bytes(), &stats)
	assert.NoError(t, err)

	require.NotEmpty(t, stats)

	found := false
	for _, s := range stats {
		if s.Version == "1.0.0" && s.ChannelName == "channel" && s.Arch == "AMD64" {
			assert.Equal(t, 5, s.Instances)
			found = true
			break
		}
	}
	assert.True(t, found, "seeded instance_stats row was not found in response")
}
