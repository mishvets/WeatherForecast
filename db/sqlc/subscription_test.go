package db

import (
	"context"
	"testing"

	"github.com/mishvets/WeatherForecast/util"
	"github.com/stretchr/testify/require"
)

func createRandomSubscription(t *testing.T) Subscription {
	arg := CreateSubscriptionParams{
		Email:     util.RandomEmail(),
		City:      util.RandomCity(),
		Frequency: randomFrequency(),
	}

	subscription, err := testQueries.CreateSubscription(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, subscription)

	require.Equal(t, arg.Email, subscription.Email)
	require.Equal(t, arg.City, subscription.City)
	require.Equal(t, arg.Frequency, subscription.Frequency)

	require.False(t, subscription.Confirmed)
	require.NotZero(t, subscription.ID)
	require.NotZero(t, subscription.CreatedAt)
	require.NotZero(t, subscription.Token)

	return subscription
}

func randomFrequency() FrequencyEnum {
	avalFreqEnum := AllFrequencyEnumValues()
	i := util.RandomInt(0, int64(len(avalFreqEnum)-1))
	return avalFreqEnum[i]
}

func TestCreateSubscription(t *testing.T) {
	createRandomSubscription(t)
}
