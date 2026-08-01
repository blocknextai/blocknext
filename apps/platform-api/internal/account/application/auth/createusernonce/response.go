package createusernonce

type CreateUserNonceResponse struct {
	Nonce        string `json:"nonce"`
	URL          string `json:"url"`
	LoginMessage string `json:"loginMessage,omitempty"`
}
