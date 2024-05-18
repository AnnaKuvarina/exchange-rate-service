package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorCode int

const (
	RequestUnmarshallingError ErrorCode = iota
	StoreError
	AlreadyExistsError
	InvalidEmailError
)

type ErrorMessage string

const (
	InvalidRequestError       ErrorMessage = "invalid request"
	InvalidEmailFormatError   ErrorMessage = "invalid email"
	InternalError             ErrorMessage = "internal error"
	SubscriptionAlreadyExists ErrorMessage = "subscription with provided email already exists"
)

type ErrorResponse struct {
	Message ErrorMessage `json:"message"`
	Type    string       `json:"type"`
	Code    ErrorCode    `json:"code"`
}

func WriteErrorResponse(w http.ResponseWriter, errResp *ErrorResponse, httpStatus int) {
	jsonStr, err := json.Marshal(errResp)
	if err != nil {
		log.Error().Err(err).Msg("Failed marshaling JSON response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
	_, _ = w.Write(jsonStr)
}

func WriteResponse[ResponsePayload any](w http.ResponseWriter, response ResponsePayload, statusCode int) {
	var (
		toMarshal any = response
		err       error
	)
	respJSON := make([]byte, 0)
	if toMarshal != nil {
		respJSON, err = json.Marshal(toMarshal)
		if err != nil {
			log.Error().Err(err).Msg("Failed marshaling JSON response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(respJSON)
}

func NewErrorResponse(msg ErrorMessage, code ErrorCode) *ErrorResponse {
	return &ErrorResponse{
		Message: msg,
		Type:    "ERROR",
		Code:    code,
	}
}
