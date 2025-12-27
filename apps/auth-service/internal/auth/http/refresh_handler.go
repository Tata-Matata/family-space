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
type RefreshService interface {
	Refresh(
		ctx context.Context,
		rawRefreshToken string,
	) (accessToken string, refreshToken string, err error)
}

type RefreshHandler struct {
	refreshSvc RefreshService
	accessTTL  time.Duration
}

func NewRefreshHandler(
	refreshSvc RefreshService,
	accessTTL time.Duration,
) *RefreshHandler {
	return &RefreshHandler{
		refreshSvc: refreshSvc,
		accessTTL:  accessTTL,
	}
}

func (handler *RefreshHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	accessToken, newRefreshToken, err :=
		handler.refreshSvc.Refresh(request.Context(), req.RefreshToken)

	if err != nil {
		if errors.Is(err, errs.ErrInvalidRefreshToken) {
			http.Error(response, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(response, "internal error", http.StatusInternalServerError)
		return
	}

	resp := refreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(handler.accessTTL.Seconds()),
	}

	response.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(response).Encode(resp)
}
