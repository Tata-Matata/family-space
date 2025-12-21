package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims

	FamilyID string `json:"family_id"`
	Role     string `json:"role"`
}
