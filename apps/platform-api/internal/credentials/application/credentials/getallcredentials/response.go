package getallcredentials

import (
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
)

type GetAllCredentialsResponse struct {
	Items      []credentialsApplicationCredentials.CredentialResponse `json:"items"`
	TotalCount int64                                                  `json:"totalCount"`
}
