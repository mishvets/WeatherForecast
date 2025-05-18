package db

import "context"

type SubscribeTxParams struct {
	CreateSubscriptionParams
	AfterCreate func(subscription Subscription) error
}

func (store *SQLStore) SubscribeTx(ctx context.Context, arg SubscribeTxParams) (Subscription, error) {
	var result Subscription
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result, err = q.CreateSubscription(ctx, arg.CreateSubscriptionParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result)
	})

	return result, err
}
