package metadata

type Service struct {
	issuerURL string
}

func NewService(
	issuerURL string,
) *Service {
	return &Service{
		issuerURL: issuerURL,
	}
}
