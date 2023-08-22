package utils

import (
	"github.com/google/uuid"
	"strings"
)

func GenerateUUID() string {
	id, _ := uuid.NewRandom()
	return strings.ReplaceAll(id.String(), "-", "")
}
