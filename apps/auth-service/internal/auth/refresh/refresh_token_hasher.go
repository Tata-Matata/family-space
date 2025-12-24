package refresh

// intentionally duplicate of Hasher interface from hasher.go
// Even if the method signatures look identical today,
// they represent different security domains.
type RefreshTokenHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}
