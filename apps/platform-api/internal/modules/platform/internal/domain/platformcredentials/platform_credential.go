package platformcredentials

type PlatformCredential struct {
	ID   string
	Data map[string]any
}

func NewPlatformCredential(
	id string,
	data map[string]any,
) *PlatformCredential {
	return &PlatformCredential{
		ID:   id,
		Data: data,
	}
}
