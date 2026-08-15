package handler

import (
	"encoding/json"
	"github.com/TheLonger011/ReLon/internal/middleware"
	"github.com/TheLonger011/ReLon/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	service *service.CommentService
}

type commentRequest struct {
	Content string `json:"content"`
}

func NewCommentHandler(service *service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := chi.URLParam(r, "id")
	postID, err := uuid.Parse(paramID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req commentRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := h.service.CreateComment(ctx, authorID, postID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(comment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paramID := chi.URLParam(r, "id")
	postID, err := uuid.Parse(paramID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	comment, err := h.service.GetCommentsByPostID(ctx, postID, limitInt, offsetInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	commentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteComment(ctx, commentID, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req commentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateComment(ctx, commentID, authorID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)

}
