package generator

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	cardNumberLength = 16
	cvvCodeLength    = 3
)

// Generator struct.
type Generator struct {
	countryCode string
	bankCode    string

	digitRunes []rune
}

// NewGenerator constructor for Generator.
func NewGenerator(countryCode string, bankCode string) *Generator {
	return &Generator{
		countryCode: countryCode,
		bankCode:    bankCode,
		digitRunes:  []rune("1234567890"),
	}
}

// GenerateRandomIban generates a random iban and returns it.
func (g Generator) GenerateRandomIban() string {
	return fmt.Sprintf("%s%s%s00000%s", g.countryCode, g.randStringRunes(2), g.bankCode, g.randStringRunes(14))
}

// GenerateRandomCardNumber generates a random card number and returns it.
func (g Generator) GenerateRandomCardNumber() string {
	return g.randStringRunes(cardNumberLength)
}

// GenerateRandomCvv generates a random cvv code and returns it.
func (g Generator) GenerateRandomCvv() string {
	return g.randStringRunes(cvvCodeLength)
}

// randStringRunes generates the given number of random bytes and returns a string.
func (g Generator) randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = g.digitRunes[rand.Intn(len(g.digitRunes))]
	}
	return string(b)
}
