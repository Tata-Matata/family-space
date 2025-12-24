package password

import "testing"

func TestBcryptHasher_HashAndCompare(test *testing.T) {
	hasher := &BcryptHasher{}

	password := "my-very-secure-password"

	hash, err := hasher.Hash(password)
	if err != nil {
		test.Fatalf("hash failed: %v", err)
	}

	if hash == password {
		test.Fatalf("hash must not equal plaintext password")
	}

	if err := hasher.Compare(hash, password); err != nil {
		test.Fatalf("expected password to match")
	}
}

func TestBcryptHasher_CompareFails(test *testing.T) {
	hasher := &BcryptHasher{}

	password := "correct-password"
	wrong := "wrong-password"

	hash, err := hasher.Hash(password)
	if err != nil {
		test.Fatalf("hash failed: %v", err)
	}

	if err := hasher.Compare(hash, wrong); err == nil {
		test.Fatalf("expected compare to fail for wrong password")
	}
}
