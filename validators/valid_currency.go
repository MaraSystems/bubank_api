package validators

import (
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/go-playground/validator/v10"
)

var ValidateCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return utils.IsSupportedCurrency(currency)
}
