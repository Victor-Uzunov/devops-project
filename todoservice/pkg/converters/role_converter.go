package converters

import (
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
)

var roleMap = map[string]constants.Role{
	"admin":  constants.Admin,
	"reader": constants.Reader,
	"writer": constants.Writer,
}

func ToRole(role string) (constants.Role, error) {
	if role, ok := roleMap[role]; ok {
		return role, nil
	}
	return "", fmt.Errorf("role %s not found", role)
}

var roleRevereMap = map[constants.Role]string{
	constants.Admin:  "admin",
	constants.Reader: "reader",
	constants.Writer: "writer",
}

func RoleToString(role constants.Role) (string, error) {
	if str, ok := roleRevereMap[role]; ok {
		return str, nil
	}
	return "", fmt.Errorf("role %s not found", role)
}
