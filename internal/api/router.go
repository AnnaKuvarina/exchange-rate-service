package api

import (
	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/subscribe", handler.CreateSubscription).Methods("POST")
	router.HandleFunc("/rate", handler.GetRate).Methods("GET")

	return router
}
