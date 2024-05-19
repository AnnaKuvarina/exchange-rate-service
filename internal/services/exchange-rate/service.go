package exchange_rate

import "context"

type ExchangeRateService interface {
	GetRate(ctx context.Context, currencyCodeA, currencyCodeB int) (*CurrencyRate, error)
}
