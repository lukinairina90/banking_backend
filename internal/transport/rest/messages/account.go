package messages

// Account object representation of response connection with accounts functionality.
type Account struct {
	ID         int     `json:"id"`
	Iban       string  `json:"iban"`
	UserID     int     `json:"user_id"`
	CurrencyID int     `json:"currency_id"`
	Blocked    bool    `json:"blocked"`
	Amount     float64 `json:"amount"`
}

// CreateAccountRequestBody object representation of response.
type CreateAccountRequestBody struct {
	CurrencyID int `json:"currency_id" binding:"required,gte=1,lte=3"`
}

// DepositAccountRequestBody object representation of response.
type DepositAccountRequestBody struct {
	Amount float64 `json:"amount" binding:"required,gte=1"`
}

// TransferAccountRequestBody object representation of response.
type TransferAccountRequestBody struct {
	Iban   string  `json:"iban" binding:"required,gte=29"`
	Amount float64 `json:"amount" binding:"required,gte=1"`
}
