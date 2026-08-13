package entity

type RoleEntity struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Permissions []PermissionEntity `json:"permissions"`
}

type PermissionEntity struct {
	ID       int64  `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}
