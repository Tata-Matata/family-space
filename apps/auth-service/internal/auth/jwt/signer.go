package jwt

type TokenSigner interface {
	SignAccessToken(user User, membership Membership) (string, error)
}
