package constants

type Role string

const (
	Reader  Role = "reader"
	Writer  Role = "writer"
	Admin   Role = "admin"
	Invalid Role = "invalid"
)

func RolePower(r Role) int {
	switch r {
	case Reader:
		return 1
	case Writer:
		return 2
	case Admin:
		return 3
	default:
		return 0
	}
}
