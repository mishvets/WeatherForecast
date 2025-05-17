package db

import "context"

type SubscribeTxParams struct {
	CreateSubscriptionParams
	AfterCreate func(subscription Subscription) error
}

type SubscribeTxResult struct {
	Subscription Subscription
}

func (store *SQLStore) SubscribeTx(ctx context.Context, arg SubscribeTxParams) (SubscribeTxResult, error) {
	var result SubscribeTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Subscription, err = q.CreateSubscription(ctx, arg.CreateSubscriptionParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Subscription)
	})

	return result, err
}
