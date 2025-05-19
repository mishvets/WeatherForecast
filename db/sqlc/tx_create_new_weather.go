package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type CreateNewWeatherTxParams struct {
	CreateWeatherParams
	ID int64
}

func (store *SQLStore) CreateNewWeatherTx(ctx context.Context, arg CreateNewWeatherTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		// If requested subscription not exist anymore, the weather shouldn't be created
		isExist, err := q.IsSubscriptionExist(ctx, arg.ID)
		if (!isExist && err == nil) || errors.Is(err, sql.ErrNoRows) {
			log.Printf("The requested subscription(%v) not exist anymore", arg.ID)
			return sql.ErrNoRows
		}
		if err != nil {
			return err
		}

		_, err = q.CreateWeather(ctx, arg.CreateWeatherParams)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
