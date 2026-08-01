package updateuserpreferences

type ThemeResponse struct {
	Color string `json:"color"`
	Mode  string `json:"mode"`
}

type UpdateUserPreferencesResponse struct {
	Theme    ThemeResponse `json:"theme"`
	Language string        `json:"language"`
}
