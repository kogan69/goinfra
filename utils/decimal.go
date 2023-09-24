package utils

import "github.com/shopspring/decimal"

func RoundMoney(d decimal.Decimal) decimal.Decimal {
	if d.Abs().LessThan(decimal.NewFromFloat(0.01)) {
		return d.Round(3)
	}

	return d.Round(2)
}
