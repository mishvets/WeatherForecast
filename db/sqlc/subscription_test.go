package db

import (
	"context"
	"database/sql"
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

func TestConfirmSubscription(t *testing.T) {
	subscription := createRandomSubscription(t)

	arg := ConfirmSubscriptionParams{
		Token:     subscription.Token,
		Confirmed: true,
	}

	confirmedSubscription, err := testQueries.ConfirmSubscription(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, confirmedSubscription)

	require.Equal(t, subscription.ID, confirmedSubscription.ID)
	require.Equal(t, subscription.Email, confirmedSubscription.Email)
	require.Equal(t, subscription.City, confirmedSubscription.City)
	require.Equal(t, subscription.Frequency, confirmedSubscription.Frequency)

	require.True(t, confirmedSubscription.Confirmed)
	require.Equal(t, subscription.Token, confirmedSubscription.Token)
	require.Equal(t, subscription.CreatedAt, confirmedSubscription.CreatedAt)
}

func TestDeleteSubscription(t *testing.T) {
	subs := createRandomSubscription(t)

	subsBeforeDel, err := testQueries.GetSubscription(context.Background(), subs.Email)
	require.NoError(t, err)
	require.Equal(t, subs.Token, subsBeforeDel.Token)

	deletedToken, err := testQueries.DeleteSubscription(context.Background(), subs.Token)
	require.NoError(t, err)
	require.Equal(t, subs.Token, deletedToken)

	subscription2, err := testQueries.GetSubscription(context.Background(), subs.Email)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, subscription2)
}

func TestGetEmailsForUpdate(t *testing.T) {
	const expectedLen = 5
	var city = util.RandomCity()
	var freq = randomFrequency()

	for range expectedLen {
		arg := CreateSubscriptionParams{
			Email:     util.RandomEmail(),
			City:      city,
			Frequency: freq,
		}

		subs, err := testQueries.CreateSubscription(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, subs)

		confirmArg := ConfirmSubscriptionParams{Token: subs.Token, Confirmed: true}
		_, err = testQueries.ConfirmSubscription(context.Background(), confirmArg)
		require.NoError(t, err)
	}

	arg := GetEmailsForUpdateParams{
		Frequency: freq,
		City:      city,
	}
	emails, err := testQueries.GetEmailsForUpdate(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, expectedLen, len(emails))
}

func TestIsSubscriptionExist(t *testing.T) {
	subs := createRandomSubscription(t)

	exists, err := testQueries.IsSubscriptionExist(context.Background(), subs.ID)
	require.NoError(t, err)
	require.True(t, exists)

	_, err = testQueries.DeleteSubscription(context.Background(), subs.Token)
	require.NoError(t, err)

	exists, err = testQueries.IsSubscriptionExist(context.Background(), subs.ID)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestGetCitiesForUpdate(t *testing.T) {
	const expectedLen = 1
	city1 := util.RandomCity()
	city2 := util.RandomCity()
	freq := randomFrequency()

	for i := range 10 {
		arg := CreateSubscriptionParams{
			Email:     util.RandomEmail(),
			Frequency: freq,
		}
		if i%2 == 0 {
			arg.City = city1
		} else {
			arg.City = city2
		}
		subs, err := testQueries.CreateSubscription(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, subs)

		confirmArg := ConfirmSubscriptionParams{Token: subs.Token, Confirmed: true}
		_, err = testQueries.ConfirmSubscription(context.Background(), confirmArg)
		require.NoError(t, err)
	}

	cities, err := testQueries.GetCitiesForUpdate(context.Background(), freq)
	require.NoError(t, err)

	cntCity1 := 0
	cntCity2 := 0
	for _, city := range cities {
		if city == city1 {
			cntCity1++
		} else if city == city2 {
			cntCity2++
		}
	}
	require.Equal(t, expectedLen, cntCity1)
	require.Equal(t, expectedLen, cntCity2)
}
