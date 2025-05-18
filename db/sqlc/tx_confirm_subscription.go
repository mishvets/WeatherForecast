package db

import (
	"context"

	"github.com/google/uuid"
)

type ConfirmSubscriptionTxParams struct {
	ConfirmSubscriptionParams
	AfterConfirm func(subscription Subscription) error
}

func (store *SQLStore) ConfirmSubscriptionTx(ctx context.Context, arg ConfirmSubscriptionTxParams) (uuid.UUID, error) {
	var result Subscription
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		subscription, err := q.ConfirmSubscription(ctx, arg.ConfirmSubscriptionParams)
		if err != nil {
			return err
		}

		return arg.AfterConfirm(subscription)
	})

	return result.Token, err
}
