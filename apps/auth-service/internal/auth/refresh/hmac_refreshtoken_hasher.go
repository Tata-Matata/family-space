package refresh

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

// implements RefreshTokenHasher using HMAC with SHA-256.
type HMACRefreshTokenHasher struct {
	secret []byte
}

func NewHMACRefreshTokenHasher(secret []byte) *HMACRefreshTokenHasher {
	return &HMACRefreshTokenHasher{secret: secret}
}

func (h *HMACRefreshTokenHasher) Hash(token string) (string, error) {
	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(token))
	sum := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum), nil
}

func (h *HMACRefreshTokenHasher) Compare(hash string, token string) error {
	expected, err := h.Hash(token)
	if err != nil {
		return err
	}

	if !hmac.Equal([]byte(hash), []byte(expected)) {
		return ErrInvalidRefreshToken
	}

	return nil
}
