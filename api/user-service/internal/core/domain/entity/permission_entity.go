package entity

type PermissionEntity struct {
	ID       int64  `json:"id"`
	Resource string `json:"name"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}
