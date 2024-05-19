package api

type CreateSubscriptionModel struct {
	Email string `json:"email"`
}

type RateResponse struct {
	RateBuy  float32 `json:"rateBuy"`
	RateSell float32 `json:"rateSell"`
}
