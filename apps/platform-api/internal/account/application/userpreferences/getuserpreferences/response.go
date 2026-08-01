package getuserpreferences

type Theme struct {
	Color string `json:"color"`
	Mode  string `json:"mode"`
}

type GetUserPreferencesResponse struct {
	Theme    Theme  `json:"theme"`
	Language string `json:"language"`
}
