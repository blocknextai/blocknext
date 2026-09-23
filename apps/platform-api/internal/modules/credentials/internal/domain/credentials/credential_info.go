package credentials

type CredentialInfo struct {
	Key        string
	SourceType SourceType
	Data       map[string]any
}
