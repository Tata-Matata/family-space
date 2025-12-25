package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeLogoutService struct {
	err error
}

func (f *fakeLogoutService) Logout(ctx context.Context, rawRefreshToken string) error {
	return f.err
}

func TestLogoutHandler_Success(test *testing.T) {
	fakeSvc := &fakeLogoutService{err: nil}
	handler := NewLogoutHandler(fakeSvc)

	body := []byte(`{"refresh_token":"some-token"}`)
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewReader(body))
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusNoContent {
		test.Fatalf("expected %d, got %d", http.StatusNoContent, handlerResponse.Code)
	}

	if handlerResponse.Body.Len() != 0 {
		test.Fatalf("expected empty body")
	}
}

func TestLogoutHandler_InvalidJSON(test *testing.T) {
	handler := NewLogoutHandler(&fakeLogoutService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/logout",
		bytes.NewReader([]byte("{invalid json")),
	)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

func TestLogoutHandler_MissingToken(test *testing.T) {
	handler := NewLogoutHandler(&fakeLogoutService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/logout",
		bytes.NewReader([]byte(`{}`)),
	)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

func TestLogoutHandler_ServiceFailure(test *testing.T) {
	fakeSvc := &fakeLogoutService{
		err: errors.New("db unavailable"),
	}
	handler := NewLogoutHandler(fakeSvc)

	body := []byte(`{"refresh_token":"some-token"}`)
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewReader(body))
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusInternalServerError {
		test.Fatalf("expected %d, got %d", http.StatusInternalServerError, handlerResponse.Code)
	}
}

func TestLogoutHandler_MethodNotAllowed(test *testing.T) {
	handler := NewLogoutHandler(&fakeLogoutService{})

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusMethodNotAllowed {
		test.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, handlerResponse.Code)
	}
}
