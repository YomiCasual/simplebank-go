package api

import (
	lib "simplebank/libs"

	"github.com/go-playground/validator/v10"
)


var validCurrency validator.Func = func (fieldLevel validator.FieldLevel ) bool {
		if currency, ok := fieldLevel.Field().Interface().(string); ok {
				return lib.IsSupportedCurrency(currency)
		}

		return false
}