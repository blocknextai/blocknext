package getroles

type RoleResponse struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type GetRolesResponse = []RoleResponse
