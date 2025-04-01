package randomness

import (
	"crypto/rand"
	"encoding/binary"
)

// GenerateTestRandomValue generates a cryptographically secure random value for testing.
// It uses crypto/rand.Reader to ensure a good distribution of random values.
func GenerateTestRandomValue() uint64 {
	var randomValue uint64
	if err := binary.Read(rand.Reader, binary.BigEndian, &randomValue); err != nil {
		// In a test context, we can panic as this would indicate a serious issue
		// with the system's random number generator
		panic(err)
	}
	return randomValue
}
