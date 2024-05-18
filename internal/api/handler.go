package api

import (
	"encoding/json"
	"io"
	"net/http"

	"exchange-rate-service/internal/store/subscriptions"
	"github.com/pkg/errors"
)

type Handler struct {
	subscriptionsStore subscriptions.Store
}

func NewHandler(subscriptionsStore subscriptions.Store) *Handler {
	return &Handler{subscriptionsStore: subscriptionsStore}
}

func (h *Handler) CreateSubscription(resp http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		WriteErrorResponse(resp, NewErrorResponse(InvalidRequestError, RequestUnmarshallingError), http.StatusInternalServerError)
		return
	}
	ctx := req.Context()

	defer func() {
		_ = req.Body.Close()
	}()

	var payload CreateSubscriptionModel
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteErrorResponse(resp, NewErrorResponse(InvalidRequestError, RequestUnmarshallingError), http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		WriteErrorResponse(resp, NewErrorResponse(InvalidEmailFormatError, InvalidEmailError), http.StatusBadRequest)
		return
	}

	dbSubscription, err := h.subscriptionsStore.Get(ctx, payload.Email)
	if err != nil {
		if !errors.Is(err, subscriptions.ErrNotFound) {
			WriteErrorResponse(resp, NewErrorResponse(InternalError, StoreError), http.StatusInternalServerError)
			return
		}
	}

	if dbSubscription != nil {
		WriteErrorResponse(resp, NewErrorResponse(SubscriptionAlreadyExists, AlreadyExistsError), http.StatusConflict)
		return
	}

	err = h.subscriptionsStore.CreateSubscription(ctx, payload.Email)
	if err != nil {
		WriteErrorResponse(resp, NewErrorResponse(InternalError, StoreError), http.StatusInternalServerError)
		return
	}

	WriteResponse[any](resp, nil, http.StatusCreated)
}
