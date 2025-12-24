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

func (h *RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	accessToken, newRefreshToken, err :=
		h.refreshSvc.Refresh(r.Context(), req.RefreshToken)

	if err != nil {
		if errors.Is(err, errs.ErrInvalidRefreshToken) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := refreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.accessTTL.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
