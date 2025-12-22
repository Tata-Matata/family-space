package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	authjwt "github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/jwt"
)

const AUDIENCE = "family-space-api"
const ISSUER = "family-space-auth"

func generateTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return key
}

func TestRS256Signer_SignAndVerify(t *testing.T) {
	privateKey := generateTestKey(t)
	publicKey := &privateKey.PublicKey

	signer := authjwt.NewRS256Signer(
		privateKey,
		ISSUER,
		AUDIENCE,
		15*time.Minute,
	)

	user := authjwt.User{
		ID: "user-123",
	}
	membership := authjwt.Membership{
		FamilyID: "family-456",
		Role:     "admin",
	}

	tokenString, err := signer.GenerateSignedAccessToken(user, membership)
	if err != nil {
		t.Fatalf("failed to generate signed token: %v", err)
	}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&authjwt.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		jwt.WithAudience(AUDIENCE),
		jwt.WithIssuer(ISSUER),
	)

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("token is not valid")
	}

	claims, ok := parsedToken.Claims.(*authjwt.Claims)
	if !ok {
		t.Fatal("claims type mismatch")
	}

	if claims.Subject != "user-123" {
		t.Errorf("unexpected subject: %s", claims.Subject)
	}

	if claims.FamilyID != "family-456" {
		t.Errorf("unexpected family_id: %s", claims.FamilyID)
	}

	if claims.Role != "admin" {
		t.Errorf("unexpected role: %s", claims.Role)
	}
}

func TestRS256Signer_ExpiredToken(t *testing.T) {
	privateKey := generateTestKey(t)
	publicKey := &privateKey.PublicKey

	signer := authjwt.NewRS256Signer(
		privateKey,
		ISSUER,
		AUDIENCE,
		-1*time.Minute, // already expired
	)

	user := authjwt.User{ID: "user-123"}
	membership := authjwt.Membership{
		FamilyID: "family-456",
		Role:     "admin",
	}

	tokenString, err := signer.GenerateSignedAccessToken(user, membership)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = jwt.ParseWithClaims(
		tokenString,
		&authjwt.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		jwt.WithAudience(AUDIENCE),
		jwt.WithIssuer(ISSUER),
	)

	if err == nil {
		t.Fatal("expected token to be expired")
	}
}

func TestRS256Signer_WrongAudience(t *testing.T) {
	privateKey := generateTestKey(t)
	publicKey := &privateKey.PublicKey

	signer := authjwt.NewRS256Signer(
		privateKey,
		ISSUER,
		AUDIENCE,
		15*time.Minute,
	)

	tokenString, _ := signer.GenerateSignedAccessToken(
		authjwt.User{ID: "user-123"},
		authjwt.Membership{FamilyID: "f", Role: "admin"},
	)

	_, err := jwt.ParseWithClaims(
		tokenString,
		&authjwt.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		jwt.WithAudience("some-other-api"),
		jwt.WithIssuer("family-space-auth"),
	)

	if err == nil {
		t.Fatal("expected audience validation to fail")
	}
}
