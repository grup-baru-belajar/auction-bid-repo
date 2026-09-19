package models

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	ID       int64
	Name     string
	Username string
	Password string
	Role     Role
}