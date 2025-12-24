package refresh

import (
	"crypto/rand"
	"encoding/base64"
)

type SecureRefreshTokenGenerator struct{}

func (g *SecureRefreshTokenGenerator) Generate() (string, error) {
	b := make([]byte, 32) // 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
