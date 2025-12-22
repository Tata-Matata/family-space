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

// fakeLoginService is a mock implementation of LoginService for testing.
// we can configure it to return a specific token or error.
type fakeLoginService struct {
	token string
	err   error
}

// we skip real authentication and return predefined values
// so tests can focus on handling the token, erros and return status codes.
func (f *fakeLoginService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	return f.token, f.err
}

const EXPIRATION_SECONDS time.Duration = 900 // 15 minutes
const TOKEN = "test.jwt.token"
const TOKEN_TYPE = "Bearer"

var REQUEST_CREDS = map[string]string{
	"email":    "user@",
	"password": "password123",
}

func TestLoginHandler_Success(test *testing.T) {

	//PREPARE
	fakeSvc := &fakeLoginService{
		token: TOKEN,
		err:   nil,
	}
	handler := createLoginHandler(fakeSvc)

	// login service will not authenticate, just return the token
	// so these values are arbitrary
	req := createRequest(REQUEST_CREDS)
	handlerResponse := httptest.NewRecorder()

	// ACT
	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusOK {
		test.Fatalf("expected %d, got %d", http.StatusOK, handlerResponse.Code)
	}

	// ASSERT
	// parse response and check that it contains the expected token
	var responseMap map[string]interface{}
	if err := json.NewDecoder(handlerResponse.Body).Decode(&responseMap); err != nil {
		test.Fatalf("invalid JSON response %v. Error: %v", handlerResponse.Body, err)
	}

	if responseMap["access_token"] != TOKEN {
		test.Errorf("expected access_token %s, got %s", TOKEN, responseMap["access_token"])
	}

	if responseMap["token_type"] != TOKEN_TYPE {
		test.Errorf("expected token type %s, got %s", TOKEN_TYPE, responseMap["token_type"])
	}

	if responseMap["expires_in"] != float64(EXPIRATION_SECONDS) {
		test.Errorf("expected expires_in %d, got %v", EXPIRATION_SECONDS, responseMap["expires_in"])
	}
}

func createLoginHandler(fakeSvc authhttp.LoginService) authhttp.LoginHandler {
	return *authhttp.NewLoginHandler(
		fakeSvc,
		EXPIRATION_SECONDS*time.Second,
	)
}

func createRequest(reqBody map[string]string) *http.Request {

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewReader(body),
	)

	return req
}

func TestLoginHandler_InvalidCredentials(test *testing.T) {
	//PREPARE
	fakeSvc := &fakeLoginService{
		err: errs.ErrInvalidCredentials,
	}

	handler := createLoginHandler(fakeSvc)

	// can be arbitrary since authentication is not tested here
	// login service will return invalid credentials error
	req := createRequest(REQUEST_CREDS)

	handlerResponse := httptest.NewRecorder()

	// ACT
	handler.ServeHTTP(handlerResponse, req)

	// ASSERT
	if handlerResponse.Code != http.StatusUnauthorized {
		test.Fatalf("expected %d, got %d", http.StatusUnauthorized, handlerResponse.Code)
	}
}

func TestLoginHandler_InvalidJSON(test *testing.T) {
	handler := createLoginHandler(&fakeLoginService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewReader([]byte("{invalid json")), // malformed JSON
	)

	handlerResponse := httptest.NewRecorder()
	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected response code %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

func TestLoginHandler_MissingFields(test *testing.T) {
	handler := createLoginHandler(&fakeLoginService{})

	req := createRequest(map[string]string{"email": ""}) // empty email

	handlerResponse := httptest.NewRecorder()
	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusBadRequest {
		test.Fatalf("expected response code %d, got %d", http.StatusBadRequest, handlerResponse.Code)
	}
}

func TestLoginHandler_MethodNotAllowed(test *testing.T) {
	handler := createLoginHandler(&fakeLoginService{})
	body, _ := json.Marshal(REQUEST_CREDS)

	req := httptest.NewRequest(
		http.MethodGet, // wrong method
		"/login",
		bytes.NewReader(body),
	)

	handlerResponse := httptest.NewRecorder()
	handler.ServeHTTP(handlerResponse, req)

	if handlerResponse.Code != http.StatusMethodNotAllowed {
		test.Fatalf("expected responce code %d, got %d", http.StatusMethodNotAllowed, handlerResponse.Code)
	}
}
