package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	c "github.com/SemmiDev/chi-bank/util/currency"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func Integer(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func Str(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func Owner() string {
	return Str(6)
}

// RandomMoney generates a random amount of money
func Money() int64 {
	return Integer(0, 1000)
}

// RandomCurrency generates a random currency code
func Currency() string {
	currencies := []string{c.USD, c.EUR, c.CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email
func Email() string {
	return fmt.Sprintf("%s@email.com", Str(6))
}
