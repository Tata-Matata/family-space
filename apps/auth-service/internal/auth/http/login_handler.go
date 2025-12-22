package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

type LoginService interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type LoginHandler struct {
	loginSvc LoginService
	tokenTTL time.Duration
}

func NewLoginHandler(
	loginSvc LoginService,
	tokenTTL time.Duration,
) *LoginHandler {
	return &LoginHandler{
		loginSvc: loginSvc,
		tokenTTL: tokenTTL,
	}
}

func (handler *LoginHandler) ServeHTTP(responseWriter http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(responseWriter, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(responseWriter, "invalid request", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		http.Error(responseWriter, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := handler.loginSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			http.Error(responseWriter, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(responseWriter, "internal error", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(handler.tokenTTL.Seconds()),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(resp)
}
