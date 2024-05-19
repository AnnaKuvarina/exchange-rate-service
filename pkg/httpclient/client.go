package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	// GET for get method
	GET = "GET"
)

type ResponseObject interface{}

type RestClient interface {
	Get(context.Context, string, map[string]string, ResponseObject) <-chan error
}

type Client struct {
	client *http.Client
}

func NewHTTPClient(c *http.Client) *Client {
	return &Client{client: c}
}

type RequestError struct {
	StatusCode int
	Body       string
}

func (m RequestError) Error() string {
	message, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("Request failed, error marshaling response to json. Code: %v", m.StatusCode)
	}
	return string(message)
}

func (c *Client) sendRequest(ctx context.Context, method, baseURL string,
	requestBody io.Reader, headers map[string]string) (string, error) {
	// Creating request
	req, reqErr := http.NewRequestWithContext(ctx, method, baseURL, requestBody)
	if reqErr != nil {
		log.Error().Err(reqErr).Msg("Failed to create request object")
		return "", reqErr
	}
	// Attaching headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	// Calling request to get response
	resp, respErr := c.client.Do(req)
	if respErr != nil {
		log.Error().Err(respErr).Msg("Failed to send request")
		return "", respErr
	}
	// Defer response close
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		err = RequestError{resp.StatusCode, string(body)}
		return "", err
	}
	// Read the response body
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Error().Err(bodyErr).Msg("Failed to read response body")
		return "", bodyErr
	}
	sb := string(body)
	return sb, nil
}

func (c *Client) doSend(
	ctx context.Context,
	method, baseURL string,
	requestBody io.Reader,
	headers map[string]string,
	ch chan<- error,
	responseObj ResponseObject,
) {
	responseString, err := c.sendRequest(ctx, method, baseURL, requestBody, headers)
	if err != nil {
		log.Error().Err(err).Msg("Request failed")
		ch <- err
	}
	if responseObj != nil {
		err = json.Unmarshal([]byte(responseString), &responseObj)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse response")
			ch <- err
		}
	}
	ch <- nil
}

func (c *Client) Get(ctx context.Context, baseURL string, headers map[string]string, responseObj ResponseObject) <-chan error {
	ch := make(chan error)
	go c.doSend(ctx, GET, baseURL, nil, headers, ch, responseObj)
	return ch
}
