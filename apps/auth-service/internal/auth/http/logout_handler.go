package http

import (
	"context"
	"encoding/json"
	"net/http"
)

// Service interface expected by handler
type LogoutService interface {
	Logout(ctx context.Context, refreshToken string) error
}

type LogoutHandler struct {
	logoutSvc LogoutService
}

func NewLogoutHandler(logoutSvc LogoutService) *LogoutHandler {
	return &LogoutHandler{logoutSvc: logoutSvc}
}

func (handler *LogoutHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(response, "invalid request", http.StatusBadRequest)
		return
	}

	if err := handler.logoutSvc.Logout(request.Context(), req.RefreshToken); err != nil {
		// we DO NOT say logout failed, we say server failed
		http.Error(response, "internal error", http.StatusInternalServerError)
		return
	}

	// 204 no content to return
	response.WriteHeader(http.StatusNoContent)
}
