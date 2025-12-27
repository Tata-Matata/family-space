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

// Service interface expected by handler
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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (handler *LoginHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody loginRequest
	if err := json.NewDecoder(request.Body).Decode(&reqBody); err != nil {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	reqBody.Email = strings.TrimSpace(reqBody.Email)
	if reqBody.Email == "" || reqBody.Password == "" {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := handler.loginSvc.Login(request.Context(), reqBody.Email, reqBody.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			http.Error(response, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(response, "internal error", http.StatusInternalServerError)
		return
	}

	respBody := loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(handler.tokenTTL.Seconds()),
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(respBody)
}
