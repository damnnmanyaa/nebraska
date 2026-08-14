package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetInstanceStatsLatest(ctx echo.Context) error {
	result, err := h.db.GetInstanceStatsLatest()
	if err != nil {
		l.Error().Err(err).Msg("getInstanceStatsLatest - getting latest instance stats")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	// Empty slice is a valid "no data yet" response (not 404).
	return ctx.JSON(http.StatusOK, result)
}
