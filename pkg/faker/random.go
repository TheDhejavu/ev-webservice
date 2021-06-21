package faker

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefgijklmnopqrstuvwxyz"

// Randomint generates random number between Min and Max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUser generates random user
func RandomUser() string {
	return RandomString(6)
}

// RandomEmailgenerates random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
