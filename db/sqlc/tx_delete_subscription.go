package db

import (
	"context"

	"github.com/google/uuid"
)

type DeleteSubscriptionTxParams struct {
	Token uuid.UUID
	City string
	Email string
	AfterDelete func(city string, email string) error
}

func (store *SQLStore) DeleteSubscriptionTx(ctx context.Context, arg DeleteSubscriptionTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		_, err := q.DeleteSubscription(ctx, arg.Token)
		if err != nil {
			return err
		}

		return arg.AfterDelete(arg.City, arg.Email)
	})

	return err
}
