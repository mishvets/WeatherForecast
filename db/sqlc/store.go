package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	SubscribeTx(ctx context.Context, arg SubscribeTxParams) (Subscription, error)
	ConfirmSubscriptionTx(ctx context.Context, arg ConfirmSubscriptionTxParams) (uuid.UUID, error)
	DeleteSubscriptionTx(ctx context.Context, arg DeleteSubscriptionTxParams) error
	CreateNewWeatherTx(ctx context.Context, arg CreateNewWeatherTxParams) error
	GetCitiesForUpdate(ctx context.Context, frequency FrequencyEnum) ([]string, error)
}

// SQLStore providers all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
