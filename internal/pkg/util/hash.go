package util

import (
	"github.com/stakater/Reloader/internal/pkg/crypto"
)

// GenerateSHAFromString generates a SHA from the provided string
func GenerateSHAFromString(data string) string {
	return crypto.GenerateSHA(data)
}
