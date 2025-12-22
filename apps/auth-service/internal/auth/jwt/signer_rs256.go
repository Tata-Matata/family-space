package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type User = domain.User
type Membership = domain.Membership

type RS256Signer struct {
	privateKey *rsa.PrivateKey
	issuer     string
	audience   string
	ttl        time.Duration
}

func NewRS256Signer(
	privateKey *rsa.PrivateKey,
	issuer string,
	audience string,
	ttl time.Duration,
) *RS256Signer {
	return &RS256Signer{
		privateKey: privateKey,
		issuer:     issuer,
		audience:   audience,
		ttl:        ttl,
	}
}

func (s *RS256Signer) GenerateSignedAccessToken(
	user User,
	membership Membership,
) (string, error) {

	now := time.Now()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  []string{s.audience},
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
		FamilyID: membership.FamilyID,
		Role:     membership.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}
