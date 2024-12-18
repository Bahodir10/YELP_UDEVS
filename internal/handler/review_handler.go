package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"YALP/internal/service"
	"YALP/pkg/response"

	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	reviewSvc service.ReviewService
}

func NewReviewHandler(r service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: r}
}

func (h *ReviewHandler) ListReviewsForBusiness(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	bid, _ := strconv.ParseInt(idStr, 10, 64)
	reviews, err := h.reviewSvc.ListReviewsForBusiness(bid)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	idStr := mux.Vars(r)["id"]
	bid, _ := strconv.ParseInt(idStr, 10, 64)

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	rv, err := h.reviewSvc.CreateReview(bid, userID, req.Rating, req.Comment)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, rv)
}
