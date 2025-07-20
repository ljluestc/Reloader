package testutil
package testutil

import "math/rand"

// GenerateRandomSequence generates a random sequence for test resources
func GenerateRandomSequence(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
import (
	"math/rand"
	"time"
)

func init() {
	// Initialize random seed for tests
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomSequence generates a random sequence of characters of the given length
func GenerateRandomSequence(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
