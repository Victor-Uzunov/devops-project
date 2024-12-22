package time

import (
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type Time struct{}

func (t Time) Now() time.Time {
	return time.Now()
}

func TimeToString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(constants.DateFormat)
	return &formatted
}
