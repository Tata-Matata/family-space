package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authhttp "github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/http"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

type fakeRefreshService struct {
	accessToken  string
	refreshToken string
	err          error
}

func (f *fakeRefreshService) Refresh(
	ctx context.Context,
	rawRefreshToken string,
) (string, string, error) {
	return f.accessToken, f.refreshToken, f.err
}

func TestRefreshHandler_Success(test *testing.T) {
	fakeSvc := &fakeRefreshService{
		accessToken:  "new.jwt.token",
		refreshToken: "new-refresh-token",
	}

	handler := authhttp.NewRefreshHandler(fakeSvc, 15*time.Minute)

	body := []byte(`{"refresh_token":"old-refresh-token"}`)
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusOK {
		test.Fatalf("expected %d, got %d", http.StatusOK, handlerResponse.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(handlerResponse.Body).Decode(&resp); err != nil {
		test.Fatalf("invalid JSON response")
	}

	if resp["access_token"] != "new.jwt.token" {
		test.Fatalf("unexpected access_token")
	}

	if resp["refresh_token"] != "new-refresh-token" {
		test.Fatalf("unexpected refresh_token")
	}

	if resp["token_type"] != "Bearer" {
		test.Fatalf("unexpected token_type")
	}
}

// 401
func TestRefreshHandler_InvalidToken(test *testing.T) {
	fakeSvc := &fakeRefreshService{
		err: errs.ErrInvalidRefreshToken,
	}

	handler := authhttp.NewRefreshHandler(fakeSvc, 15*time.Minute)

	body := []byte(`{"refresh_token":"bad-token"}`)
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusUnauthorized {
		test.Fatalf("expected %d, got %d", http.StatusUnauthorized, handlerResponse.Code)
	}
}

// 400
func TestRefreshHandler_InvalidJSON(test *testing.T) {
	handler := authhttp.NewRefreshHandler(&fakeRefreshService{}, 15*time.Minute)

	req := httptest.NewRequest(
		http.MethodPost,
		"/refresh",
		bytes.NewReader([]byte("{invalid json")),
	)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

// 400
func TestRefreshHandler_MissingToken(test *testing.T) {
	handler := authhttp.NewRefreshHandler(&fakeRefreshService{}, 15*time.Minute)

	req := httptest.NewRequest(
		http.MethodPost,
		"/refresh",
		bytes.NewReader([]byte(`{}`)),
	)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}
