package models

type UserRoles struct {
	Memberships map[string][]any `json:"memberships,omitempty"`
}
