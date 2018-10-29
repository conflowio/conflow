package util

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

var seedMathRandOnce sync.Once

// SeedMathRand initializes the random generator
func SeedMathRand() {
	seedMathRandOnce.Do(func() {
		if n, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64)); err == nil {
			rand.Seed(n.Int64())
		} else {
			rand.Seed(time.Now().UTC().UnixNano())
		}
	})
}

const hexBytes = "0123456789abcdef"
const hexLetters = "abcdef"

// RandHexString generates a random hexadecimal string
func RandHexString(n int, startsWithletter bool) string {
	SeedMathRand()
	b := make([]byte, n)
	for i := range b {
		if i == 0 && startsWithletter {
			b[i] = hexLetters[rand.Int63()%int64(len(hexLetters))]
		} else {
			b[i] = hexBytes[rand.Int63()%int64(len(hexBytes))]
		}
	}
	return string(b)
}
