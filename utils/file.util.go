package utils

import (
	"github.com/google/uuid"
)

func AddUUIDToString(text string) string {
	return uuid.New().String() + "-" + text
}
