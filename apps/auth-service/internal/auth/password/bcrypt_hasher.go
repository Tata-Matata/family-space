package password

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher with the given cost.
// If cost is 0, it uses bcrypt.DefaultCost. Mostly used for password hashing.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		h.cost,
	)
	return string(bytes), err
}

func (h *BcryptHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
