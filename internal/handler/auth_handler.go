package handler

import (
	"encoding/json"
	"net/http"
	"YALP/internal/service"
	"YALP/pkg/response"
)

type AuthHandler struct {
	userSvc service.UserService
}

func NewAuthHandler(u service.UserService) *AuthHandler {
	return &AuthHandler{userSvc: u}
}

// Register handler: Registers a new user
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Register the user
	u, err := h.userSvc.Register(req.Email, req.Password, req.Name)
	if err != nil {
		response.JSONError(w, http.StatusConflict, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, u)
}

// Login handler: Authenticates a user and returns a JWT token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Attempt to log in
	token, err := h.userSvc.Login(req.Email, req.Password)
	if err != nil {
		response.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Return the token
	response.JSON(w, http.StatusOK, map[string]string{"token": token})
}