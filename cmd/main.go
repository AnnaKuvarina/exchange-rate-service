package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exchange-rate-service/internal/api"
	"exchange-rate-service/internal/store/subscriptions"
	"exchange-rate-service/pkg/store/pg"
	"github.com/rs/zerolog/log"
)

func main() {
	config := LoadConfig()
	ctx, cancel := context.WithCancel(log.Logger.WithContext(context.Background()))
	defer cancel()
	// Connect to a DB and create subscriptions store
	db, err := pg.Connect(ctx, config.Database.ReadURL, config.Database.WriteURL,
		config.Database.ConnRetries, config.Database.RetryWaitDuration)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("unable to connect to database")
	}

	subscriptionsStore := subscriptions.NewPostgresStore(db)
	log.Ctx(ctx).Info().Msg("successfully connected to subscriptions store")

	errChan := make(chan error)
	// Run HTTP API handler
	go func() {
		handler := api.NewHandler(subscriptionsStore)
		server := &http.Server{
			Addr:              fmt.Sprintf(":%d", config.HTTPPort),
			Handler:           api.NewRouter(handler),
			ReadHeaderTimeout: 2 * time.Second,
		}

		log.Ctx(ctx).Info().Msgf("starting API handler on port %d", config.HTTPPort)
		errChan <- server.ListenAndServe()
	}()

	// Listen on signals for graceful shutdowns
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		log.Ctx(ctx).Warn().Str("sig", sig.String()).Msg("got termination signal, exiting")
	case err = <-errChan:
		log.Ctx(ctx).Error().Err(err).Msg("received error from error channel")
	}
}
