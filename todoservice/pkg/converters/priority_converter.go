package converters

import (
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
)

var priorityLevelMap = map[string]constants.PriorityLevel{
	"high":   constants.PriorityHigh,
	"medium": constants.PriorityMedium,
	"low":    constants.PriorityLow,
}

func ToPriorityLevel(priorityStr string) (constants.PriorityLevel, error) {
	if priority, ok := priorityLevelMap[priorityStr]; ok {
		return priority, nil
	}
	return "", fmt.Errorf("invalid priority level %s", priorityStr)
}

var priorityLevelReverseMap = map[constants.PriorityLevel]string{
	constants.PriorityLow:    "low",
	constants.PriorityMedium: "medium",
	constants.PriorityHigh:   "high",
}

func PriorityLevelToString(priority constants.PriorityLevel) (string, error) {
	if priorityStr, ok := priorityLevelReverseMap[priority]; ok {
		return priorityStr, nil
	}
	return "", fmt.Errorf("invalid priority level %s", priority)
}
