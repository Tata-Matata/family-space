package refresh

type RefreshTokenGenerator interface {
	Generate() (string, error)
}
