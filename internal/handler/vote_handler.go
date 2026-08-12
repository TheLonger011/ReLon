package handler

import (
	"encoding/json"
	"errors"
	"github.com/TheLonger011/ReLon/internal/middleware"
	"github.com/TheLonger011/ReLon/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type VoteHandler struct {
	service *service.VoteService
}

type voteRequest struct {
	VoteType int `json:"vote_type"`
}

func NewVoteHandler(service *service.VoteService) *VoteHandler {
	return &VoteHandler{service: service}
}

func (h *VoteHandler) Vote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(idParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserID(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req voteRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Vote(ctx, userID, postID, req.VoteType)
	if err != nil {
		if errors.Is(err, service.ErrInvalidVoteType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
