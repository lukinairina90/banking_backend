package messages

import "time"

// Card object representation of response connection with cards functionality.
type Card struct {
	Id             int       `json:"id"`
	AccountID      int       `json:"account_id"`
	CardNumber     string    `json:"card_number"`
	CardholderName string    `json:"cardholder_name"`
	ExpirationDate time.Time `json:"expiration_date"`
	CvvCode        string    `json:"cvv_code"`
}
