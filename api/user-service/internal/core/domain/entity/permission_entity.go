package entity

type PermissionEntity struct {
	ID       int64  `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}
