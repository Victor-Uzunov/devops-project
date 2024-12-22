package converters

import (
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
)

var visibilityMap = map[string]constants.Visibility{
	"public":  constants.VisibilityPublic,
	"private": constants.VisibilityPrivate,
	"shared":  constants.VisibilityShared,
}

func ToVisibility(visibilityStr string) (constants.Visibility, error) {
	visibility, exists := visibilityMap[visibilityStr]
	if !exists {
		return "", fmt.Errorf("invalid visibility str: %s", visibilityStr)
	}
	return visibility, nil
}

var visibilityReverseMap = map[constants.Visibility]string{
	constants.VisibilityPublic:  "public",
	constants.VisibilityPrivate: "private",
	constants.VisibilityShared:  "shared",
}

func VisibilityToString(visibility constants.Visibility) (string, error) {
	if str, ok := visibilityReverseMap[visibility]; ok {
		return str, nil
	}
	return "", fmt.Errorf("invalid visibility str: %s", visibility)
}
