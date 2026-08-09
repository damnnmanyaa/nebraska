package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetInstanceStatsLatest(ctx echo.Context) error {
	result, err := h.db.GetInstanceStatsLatest()
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.NoContent(http.StatusNotFound)
		}
		l.Error().Err(err).Msg("getInstanceStatsLatest - getting latest instance stats")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, result)
}
