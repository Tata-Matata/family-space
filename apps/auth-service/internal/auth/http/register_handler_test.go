package http_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authhttp "github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/http"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

type fakeRegistrationService struct {
	err error
}

func (f *fakeRegistrationService) Register(
	ctx context.Context,
	email, password, familyName string,
) error {
	return f.err
}

func TestRegisterHandler_Success(test *testing.T) {
	svc := &fakeRegistrationService{}
	handler := authhttp.NewRegisterHandler(svc)

	body := []byte(`{"email":"a@b.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusCreated {
		test.Fatalf("expected %d, got %d", http.StatusCreated, handlerResponse.Code)
	}
}

func TestRegisterHandler_InvalidJSON(test *testing.T) {
	handler := authhttp.NewRegisterHandler(&fakeRegistrationService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewReader([]byte("{invalid")),
	)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

func TestRegisterHandler_MissingFields(test *testing.T) {
	handler := authhttp.NewRegisterHandler(&fakeRegistrationService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewReader([]byte(`{"email":""}`)),
	)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

func TestRegisterHandler_UserExists(t *testing.T) {
	svc := &fakeRegistrationService{
		err: errs.ErrUserAlreadyExists,
	}
	handler := authhttp.NewRegisterHandler(svc)

	body := []byte(`{"email":"a@b.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestRegisterHandler_InternalError(test *testing.T) {
	svc := &fakeRegistrationService{
		err: errors.New("db down"),
	}
	handler := authhttp.NewRegisterHandler(svc)

	body := []byte(`{"email":"a@b.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusInternalServerError {
		test.Fatalf("expected %d, got %d", http.StatusInternalServerError, handlerResponse.Code)
	}
}

func TestRegisterHandler_MethodNotAllowed(test *testing.T) {
	handler := authhttp.NewRegisterHandler(&fakeRegistrationService{})

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	handlerResponse := httptest.NewRecorder()

	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusMethodNotAllowed {
		test.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, handlerResponse.Code)
	}
}
