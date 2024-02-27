package utils

import (
	"github.com/google/uuid"
)

func GenUUID() string { // GenUUID generates a new UUID and returns it as a string.
	return uuid.New().String()
}
