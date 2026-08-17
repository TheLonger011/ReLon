package handler

import (
	"github.com/TheLonger011/ReLon/internal/middleware"
	"github.com/TheLonger011/ReLon/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type CommunityMemberHandler struct {
	service *service.CommunityMemberService
}

func NewCommunityMemberHandler(service *service.CommunityMemberService) *CommunityMemberHandler {
	return &CommunityMemberHandler{service: service}
}

func (h *CommunityMemberHandler) LeaveCommunity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := chi.URLParam(r, "id")
	communityID, err := uuid.Parse(paramID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.service.LeaveCommunity(ctx, communityID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
