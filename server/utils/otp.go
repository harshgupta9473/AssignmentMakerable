package utils

import (
	"crypto/rand"
	"fmt"

	math "math/rand"

	
)



func GenerateOTPToken() (string, error) {
	length := math.Intn(16) + 10      // Random length between 10-25
	bytes := make([]byte, length/2+1) // Enough bytes for hex conversion

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes)[:length], nil // Convert to hex string and trim to 25 chars
}

