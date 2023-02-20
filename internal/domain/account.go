package domain

// Account business layer account definition
type Account struct {
	ID         int
	Iban       string
	UserID     int
	CurrencyID int
	Blocked    bool
	Amount     float64
}

// Orderings type map[string]string for filters
type Orderings map[string]string
