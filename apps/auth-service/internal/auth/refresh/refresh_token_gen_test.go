package refresh

import "testing"

func TestRefreshTokenGenerator_Unique(test *testing.T) {
	gen := &SecureRefreshTokenGenerator{}

	t1, _ := gen.Generate()
	t2, _ := gen.Generate()

	if t1 == t2 {
		test.Fatalf("tokens must be unique")
	}
}
