package utils

var supportedCurrencies = map[string]bool{
	"NGN": true,
	"USD": true,
	"EUR": true,
}

func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrencies[currency]
	return ok
}
