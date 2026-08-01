package token

type Service interface {
	Generate() (string, error)
}
