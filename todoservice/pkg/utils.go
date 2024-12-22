package pkg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/google/uuid"
	"net/mail"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal server error")
)

func NewNullableStringFromJSONRawMessage(json json.RawMessage) sql.NullString {
	nullString := sql.NullString{}
	if json != nil && string(json) != "null" {
		nullString.String = string(json)
		nullString.Valid = true
	}
	return nullString
}

func JSONRawMessageFromNullableString(sqlString sql.NullString) json.RawMessage {
	if sqlString.Valid {
		return json.RawMessage(sqlString.String)
	}
	return nil
}
func NewValidNullableString(text string) sql.NullString {
	if text == "" || text == "null" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: text,
		Valid:  true,
	}
}

func ValidateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	return nil
}

func Contains(slice []constants.Role, item constants.Role) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}

func DetermineUserRole(org string) (constants.Role, error) {
	var orgList []struct {
		Login string `json:"login"`
	}
	err := json.Unmarshal([]byte(org), &orgList)
	if err != nil {
		return constants.Invalid, err
	}

	var roles []constants.Role

	for _, el := range orgList {
		if el.Login == constants.AdminOrganization {
			return constants.Admin, nil
		} else if el.Login == constants.WriterOrganization {
			roles = append(roles, constants.Writer)
		} else if el.Login == constants.ReaderOrganization {
			roles = append(roles, constants.Reader)
		}
	}

	if Contains(roles, constants.Admin) {
		return constants.Admin, nil
	} else if Contains(roles, constants.Writer) {
		return constants.Writer, nil
	} else if Contains(roles, constants.Reader) {
		return constants.Reader, nil
	}

	return constants.Invalid, errors.New("unknown role")
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidName(name string) bool {
	if len(name) <= 3 {
		return false
	}
	return true
}

func StringToRole(role string) constants.Role {
	switch role {
	case "admin":
		return constants.Admin
	case "writer":
		return constants.Writer
	case "reader":
		return constants.Reader
	default:
		return constants.Invalid
	}
}

func NullIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
