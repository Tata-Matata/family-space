package jwt

type TokenSigner interface {
	GenerateSignedAccessToken(user User, membership Membership) (string, error)
}
