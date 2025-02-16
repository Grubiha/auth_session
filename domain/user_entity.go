package domain

type User struct {
	Id    string
	Name  string
	Phone string
	Role  string
}

const (
	UserRoleUser    = "user"
	UserRoleManager = "manager"
	UserRoleAdmin   = "admin"
)

var ValidUserRoles = map[string]bool{
	UserRoleUser:    true,
	UserRoleManager: true,
	UserRoleAdmin:   true,
}

var UserRolesLevel = map[string]int{
	UserRoleUser:    0,
	UserRoleManager: 1,
	UserRoleAdmin:   2,
}
