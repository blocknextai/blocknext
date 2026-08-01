package getprofile

type GetProfileResponse struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	IsVerified  bool     `json:"isVerified"`
}
