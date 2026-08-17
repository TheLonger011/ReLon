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

type PostHandler struct {
	service *service.PostService
}

type createPostRequest struct {
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	CommunityID *uuid.UUID `json:"community_id,omitempty"`
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req createPostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.service.CreatePost(ctx, authorID, req.Title, req.Content, req.CommunityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(idParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.service.GetPostByID(ctx, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	post, err := h.service.GetPosts(ctx, limitInt, offsetInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "search query cannot be empty", http.StatusBadRequest)
		return
	}

	search, err := h.service.SearchPosts(ctx, query, limitInt, offsetInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(search); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := h.service.DeletePost(ctx, postID, authorID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorID, err := middleware.GetUserID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req createPostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePost(ctx, postID, authorID, req.Title, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}
