package exchange_rate

type CurrencyRate struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	Date          int64   `json:"date"`
	RateBuy       float32 `json:"rateBuy"`
	RateSell      float32 `json:"rateSell"`
}
