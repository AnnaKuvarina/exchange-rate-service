package subscriptions

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const tableName = "subscriptions"

type pgSubscription struct {
	ID    string `gorm:"primaryKey"`
	Email string `gorm:"not null;type:varchar(255)"`
}

func (pgSubscription) TableName() string {
	return tableName
}

type PostgresStore struct {
	db *gorm.DB
}

func NewPostgresStore(db *gorm.DB) Store {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) CreateSubscription(ctx context.Context, email string) error {
	id := uuid.New().String()
	pgNewSubscription := pgSubscription{
		ID:    id,
		Email: email,
	}

	err := s.db.
		WithContext(ctx).
		Create(&pgNewSubscription).Error
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to create subscription. email=%s", email)
		return err
	}

	log.Ctx(ctx).Debug().Msgf("email is sucessfully subscribed, email=%s, id=%s", email, id)
	return nil
}

func (s *PostgresStore) Get(ctx context.Context, email string) (*Subscription, error) {
	var pgSub *pgSubscription
	err := s.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&pgSub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to get subscription by email: %s", email)
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &Subscription{
		ID:    pgSub.ID,
		Email: pgSub.Email,
	}, nil
}
