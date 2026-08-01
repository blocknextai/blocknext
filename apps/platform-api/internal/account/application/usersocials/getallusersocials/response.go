package getallusersocials

type UserSocialResponse struct {
	Platform  string `json:"platform"`
	URL       string `json:"url"`
	SortOrder int    `json:"sortOrder"`
}

type GetAllUserSocialsResponse = []UserSocialResponse
