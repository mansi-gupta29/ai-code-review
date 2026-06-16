package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/mansi-gupta29/ai-code-review/internal/ai"
	"github.com/mansi-gupta29/ai-code-review/internal/models"
	"github.com/mansi-gupta29/ai-code-review/internal/store"
)

type Handler struct {
	store *store.Store
	ai    *ai.Client
}

func New(store *store.Store, ai *ai.Client) *Handler {
	return &Handler{store: store, ai: ai}
}

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.store.ListReviews(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	review, err := h.store.GetReviewByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "review not found", http.StatusNotFound) // 404
			return
		}
		http.Error(w, "failed to fetch review", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req models.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	// call the AI
	reviewText, err := h.ai.ReviewCode(r.Context(), req.Language, req.Code)
	if err != nil {
		http.Error(w, "failed to generate review", http.StatusInternalServerError)
		return
	}

	// save to DB
	saved, err := h.store.SaveReview(r.Context(), models.Review{
		Language: req.Language,
		Code:     req.Code,
		Review:   reviewText,
	})
	if err != nil {
		http.Error(w, "failed to save review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(saved)
}
