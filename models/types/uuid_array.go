package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

// Scan parses a value into a UUIDArray. It supports the following types:
// * []byte
// * string
//
// The value must be a string representation of a UUID array in the following format:
// {uuid1, uuid2, ...}
//
// It returns an error if the value is of an unsupported type, or if any of the UUIDs in the array are invalid.
func (a *UUIDArray) Scan(value interface{}) error {
	// Get the string representation of the value
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("failed to parse UUIDArray: unsupported data type")
	}

	// Remove the quotes and spaces from the string
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	// Split the string into individual UUIDs
	parts := strings.Split(str, ",")
	*a = make(UUIDArray, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(strings.Trim(part, `"`)) // Remove quotes and space if present
		if part == "" {
			continue
		}

		// Parse the UUID
		u, err := uuid.Parse(part)
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %v", err)
		}

		// Append the UUID to the array
		*a = append(*a, u)
	}

	return nil
}

// Value returns a driver.Value for a UUIDArray. It formats the UUIDs as a PostgreSQL array literal.
//
// If the array is empty, it returns "{}".
//
// The driver.Value returned by Value is safe to pass to the database.
func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		// Empty array, return an empty PostgreSQL array literal
		return "{}", nil
	}

	// Format the UUIDs as a PostgreSQL array literal
	postgreFormat := make([]string, 0, len(a))

	for _, v := range a {
		// Format each UUID as a string in the PostgreSQL array literal
		postgreFormat = append(postgreFormat, fmt.Sprintf(`"%s"`, v.String()))
	}

	value := strings.Join(postgreFormat, ",")
	// Return the formatted string as a driver.Value
	return fmt.Sprintf("{%s}", value), nil
}

// GormDataType returns the data type for a UUIDArray as a string.
//
// It is used by Gorm to determine the type of the field in the database.
func (UUIDArray) GormDataType() string {
	return "uuid[]"
}