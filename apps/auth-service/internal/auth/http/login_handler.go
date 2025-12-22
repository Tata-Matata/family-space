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

func (handler *LoginHandler) ServeHTTP(responseWriter http.ResponseWriter, httpReq *http.Request) {
	if httpReq.Method != http.MethodPost {
		http.Error(responseWriter, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody loginRequest
	if err := json.NewDecoder(httpReq.Body).Decode(&reqBody); err != nil {
		http.Error(responseWriter, "invalid request", http.StatusBadRequest)
		return
	}

	reqBody.Email = strings.TrimSpace(reqBody.Email)
	if reqBody.Email == "" || reqBody.Password == "" {
		http.Error(responseWriter, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := handler.loginSvc.Login(httpReq.Context(), reqBody.Email, reqBody.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			http.Error(responseWriter, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(responseWriter, "internal error", http.StatusInternalServerError)
		return
	}

	respBody := loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(handler.tokenTTL.Seconds()),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(respBody)
}
