package subscriptions

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("subscription not found")

type Subscription struct {
	ID    string
	Email string
}

type Store interface {
	CreateSubscription(ctx context.Context, email string) error
	Get(ctx context.Context, email string) (*Subscription, error)
}
