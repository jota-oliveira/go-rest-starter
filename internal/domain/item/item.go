package item

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidName = errors.New("name must be between 1 and 120 characters")

type Item struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) > 120 {
		return "", ErrInvalidName
	}
	return name, nil
}
