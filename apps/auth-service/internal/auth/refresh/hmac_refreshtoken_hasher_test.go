package password

import "testing"

func TestRefreshTokenHasher(test *testing.T) {
	hasher := NewHMACRefreshTokenHasher([]byte("secret"))

	token := "refresh-token"
	hash, _ := hasher.Hash(token)

	if err := hasher.Compare(hash, token); err != nil {
		test.Fatalf("expected token to match")
	}

	if err := hasher.Compare(hash, "wrong"); err == nil {
		test.Fatalf("expected mismatch")
	}
}
