package exchange_rate

import (
	"context"
	"fmt"
	"time"

	"exchange-rate-service/pkg/httpclient"
	"exchange-rate-service/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	currencyPath = "/bank/currency"
)

type MonobankExchangeRateService struct {
	client     httpclient.RestClient
	baseURL    string
	ratesCache map[string]*CurrencyRate
	cacheDate  string
}

func NewMonobankExchangeRateService(client httpclient.RestClient, baseURL string) ExchangeRateService {
	return &MonobankExchangeRateService{
		client:  client,
		baseURL: baseURL,
	}
}

func (s *MonobankExchangeRateService) GetRate(ctx context.Context, currencyCodeA, currencyCodeB int) (*CurrencyRate, error) {
	if s.ratesCache == nil || !s.isTodayCache() {
		s.cacheDate = utils.ConvertUnixTimestampToDate(time.Now().Unix())
		err := s.processRates(ctx)
		if err != nil {
			return nil, err
		}
	}

	rate, ok := s.ratesCache[getCurrencyRateKey(currencyCodeA, currencyCodeB)]
	if !ok {
		return nil, fmt.Errorf("enable to find rates for provided curency codes")
	}

	return rate, nil
}

func (s *MonobankExchangeRateService) processRates(ctx context.Context) error {
	reqURL := fmt.Sprintf("%s%s", s.baseURL, currencyPath)

	var rates []*CurrencyRate
	err := s.doGet(ctx, reqURL, &rates, "failed to get currency rates from monobank")
	if err != nil {
		log.Error().Err(err).Msg("failed to send request to monobank")
		return err
	}

	ratesCache := make(map[string]*CurrencyRate, len(rates)*2)
	for _, r := range rates {
		key1 := getCurrencyRateKey(r.CurrencyCodeA, r.CurrencyCodeB)
		key2 := getCurrencyRateKey(r.CurrencyCodeB, r.CurrencyCodeA)
		ratesCache[key1] = r
		ratesCache[key2] = r
	}

	s.ratesCache = ratesCache
	return nil
}

func (s *MonobankExchangeRateService) doGet(ctx context.Context, requestURL string, res interface{}, errMsg string) error {
	err := <-s.client.Get(ctx, requestURL, nil, res)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	return nil
}

func (s *MonobankExchangeRateService) isTodayCache() bool {
	if s.cacheDate == "" {
		return false
	}
	return utils.ConvertUnixTimestampToDate(time.Now().Unix()) == s.cacheDate
}

func getCurrencyRateKey(currencyCodeA, currencyCodeB int) string {
	return fmt.Sprintf("%d-%d", currencyCodeA, currencyCodeB)
}
