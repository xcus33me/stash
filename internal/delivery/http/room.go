package httpdelivery

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	roomapp "github.com/xcus33me/stash/internal/application/room"
)

const ownerTokenHeader = "X-Owner-Token"

type RoomHandler struct {
	uc  *roomapp.Usecase
	log *slog.Logger
}

func NewRoomHandler(uc *roomapp.Usecase, log *slog.Logger) *RoomHandler {
	return &RoomHandler{uc: uc, log: log.With("component", "room_handler")}
}

type createRoomRequest struct {
	Title        string  `json:"title"`
	Desciption   *string `json:"description,omitempty"`
	MaxSizeBytes *int64  `json:"max_size_bytes,omitempty"`
	FileTTLSec   *int64  `json:"file_ttl_sec,omitempty"`
}

type createRoomResponse struct {
	RoomID     uuid.UUID `json:"room_id"`
	OwnerToken string    `json:"owner_token"`
}

type roomResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Locked      bool      `json:"locked"`
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	var ttl *time.Duration
	if req.FileTTLSec != nil {
		d := time.Duration(*req.FileTTLSec) * time.Second
		ttl = &d
	}

	out, err := h.uc.CreateRoom(r.Context(), roomapp.CreateRoomInput{
		Title:        req.Title,
		Description:  req.Desciption,
		MaxSizeBytes: req.MaxSizeBytes,
		FileTTL:      ttl,
	})
	if err != nil {
		writeError(w, h.log, err)
		return
	}

	writeJSON(w, http.StatusCreated, createRoomResponse{
		RoomID:     out.RoomID,
		OwnerToken: out.OwnerToken,
	})
}

func (h *RoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid room id"})
		return
	}

	out, err := h.uc.GetRoom(r.Context(), roomapp.GetRoomInput{RoomID: id})
	if err != nil {
		writeError(w, h.log, err)
		return
	}

	writeJSON(w, http.StatusOK, roomResponse{
		ID:          out.ID,
		Title:       out.Title,
		Description: out.Description,
		Locked:      out.Locked,
	})
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.withOwner(w, r, func(ctx context.Context, id uuid.UUID, token string) error {
		return h.uc.DeleteRoom(ctx, roomapp.DeleteRoomInput{RoomID: id, OwnerToken: token})
	})
}

func (h *RoomHandler) Lock(w http.ResponseWriter, r *http.Request) {
	h.withOwner(w, r, func(ctx context.Context, id uuid.UUID, token string) error {
		return h.uc.LockRoom(ctx, roomapp.LockUnlockRoomInput{RoomID: id, OwnerToken: token})
	})
}

func (h *RoomHandler) Unlock(w http.ResponseWriter, r *http.Request) {
	h.withOwner(w, r, func(ctx context.Context, id uuid.UUID, token string) error {
		return h.uc.UnlockRoom(ctx, roomapp.LockUnlockRoomInput{RoomID: id, OwnerToken: token})
	})
}

func (h *RoomHandler) withOwner(
	w http.ResponseWriter,
	r *http.Request,
	fn func(ctx context.Context, id uuid.UUID, token string) error,
) {
	id, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid room id"})
		return
	}

	token := r.Header.Get(ownerTokenHeader)
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing owner token"})
		return
	}
	if err := fn(r.Context(), id, token); err != nil {
		writeError(w, h.log, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
