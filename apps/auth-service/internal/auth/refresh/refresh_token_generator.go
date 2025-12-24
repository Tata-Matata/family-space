package password

type RefreshTokenGenerator interface {
	Generate() (string, error)
}
