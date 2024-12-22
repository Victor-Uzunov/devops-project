package converters

import (
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
)

var tagTypeMap = map[string]constants.TagType{
	"work":      constants.TagWork,
	"personal":  constants.TagPersonal,
	"shopping":  constants.TagShopping,
	"health":    constants.TagHealth,
	"fitness":   constants.TagFitness,
	"finance":   constants.TagFinance,
	"important": constants.TagImportant,
	"urgent":    constants.TagUrgent,
}

func ToTagType(tagStr string) (constants.TagType, error) {
	tag, exists := tagTypeMap[tagStr]
	if !exists {
		return "", fmt.Errorf("invalid tag type %s", tagStr)
	}
	return tag, nil
}

var tagTypeReverseMap = map[constants.TagType]string{
	constants.TagWork:      "work",
	constants.TagPersonal:  "personal",
	constants.TagShopping:  "shopping",
	constants.TagHealth:    "health",
	constants.TagFitness:   "fitness",
	constants.TagFinance:   "finance",
	constants.TagImportant: "important",
	constants.TagUrgent:    "urgent",
}

func TagToString(tag constants.TagType) (string, error) {
	if str, ok := tagTypeReverseMap[tag]; ok {
		return str, nil
	}
	return "", fmt.Errorf("invalid tag type %s", tag)
}
