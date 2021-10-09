package common

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	IDR = "IDR"
	RM = "RM"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, IDR, RM:
		return true
	}
	return false
}
