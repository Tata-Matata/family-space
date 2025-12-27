package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

// Service interface expected by handler
type RegisterService interface {
	Register(ctx context.Context, email, password, familyName string) error
}

type RegisterHandler struct {
	service RegisterService
}

func NewRegisterHandler(service RegisterService) *RegisterHandler {
	return &RegisterHandler{service: service}
}

type registerRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FamilyName string `json:"family_name"`
}

func (handler *RegisterHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	err := handler.service.Register(request.Context(), req.Email, req.Password, req.FamilyName)
	if err != nil {
		handler.handleError(response, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (handler *RegisterHandler) handleError(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrUserAlreadyExists):
		http.Error(response, "user already exists", http.StatusConflict)
	default:
		http.Error(response, "internal error", http.StatusInternalServerError)
	}
}
