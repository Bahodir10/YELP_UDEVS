package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"YALP/internal/service"
	"YALP/pkg/response"

	"github.com/gorilla/mux"
)

type BusinessHandler struct {
	businessSvc service.BusinessService
}

func NewBusinessHandler(b service.BusinessService) *BusinessHandler {
	return &BusinessHandler{businessSvc: b}
}

func (h *BusinessHandler) ListBusinesses(w http.ResponseWriter, r *http.Request) {
	businesses, err := h.businessSvc.ListAll()
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, businesses)
}

func (h *BusinessHandler) SearchBusinesses(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		response.JSONError(w, http.StatusBadRequest, "query param 'q' is required")
		return
	}
	businesses, err := h.businessSvc.Search(keyword)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, businesses)
}

func (h *BusinessHandler) GetBusiness(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	b, err := h.businessSvc.GetBusiness(id)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if b == nil {
		response.JSONError(w, http.StatusNotFound, "not found")
		return
	}
	response.JSON(w, http.StatusOK, b)
}

func (h *BusinessHandler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string `json:"name"`
		Category   string `json:"category"`
		Description string `json:"description"`
		Address    string `json:"address"`
		ContactInfo string `json:"contact_info"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userID := r.Context().Value("user_id").(int64)
	b, err := h.businessSvc.CreateBusiness(req.Name, req.Category, req.Description, req.Address, req.ContactInfo, userID)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, b)
}

func (h *BusinessHandler) ClaimBusiness(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	idStr := mux.Vars(r)["id"]
	bid, _ := strconv.ParseInt(idStr, 10, 64)

	b, err := h.businessSvc.ClaimBusiness(bid, userID)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, b)
}
