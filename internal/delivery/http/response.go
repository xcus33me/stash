package httpdelivery

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/xcus33me/stash/internal/domain/file"
	"github.com/xcus33me/stash/internal/domain/room"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, log *slog.Logger, err error) {
	status, msg := mapError(err)
	if status >= 500 {
		log.Error("request failed", "err", err)
	}
	writeJSON(w, status, errorResponse{Error: msg})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, room.ErrNotFound), errors.Is(err, file.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, room.ErrInvalidOwner):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, room.ErrLocked),
		errors.Is(err, file.ErrDeleted),
		errors.Is(err, file.ErrExpired):
		return http.StatusConflict, err.Error()
	case errors.Is(err, room.ErrTitleRequired),
		errors.Is(err, file.ErrNameRequired),
		errors.Is(err, file.ErrInvalidSize):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, file.ErrQuotaExceeded):
		return http.StatusRequestEntityTooLarge, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
