package api

import (
	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/subscribe", handler.CreateSubscription).Methods("POST")

	return router
}
