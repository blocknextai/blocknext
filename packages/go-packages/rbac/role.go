package rbac

// Role represents a named RBAC role together with its relative authority score.
type Role struct {
	Name  string
	Score int
}
