package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	// "github.com/lib/pq"
	mockdb "github.com/mishvets/WeatherForecast/db/mock"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	mockservice "github.com/mishvets/WeatherForecast/service/mock"
	"github.com/mishvets/WeatherForecast/util"
	mockworker "github.com/mishvets/WeatherForecast/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type subscribeTxParamsMatcher struct {
	expected db.SubscribeTxParams
}

func (m subscribeTxParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.SubscribeTxParams)
	if !ok {
		return false
	}
	return arg.CreateSubscriptionParams == m.expected.CreateSubscriptionParams
}

func (m subscribeTxParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %+v (ignoring AfterCreate)", m.expected)
}

func TestCreateSubscriptionAPI(t *testing.T) {
	subscription := randomSubscription()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Subscription successful",
			body: gin.H{
				"email":     subscription.Email,
				"city":      subscription.City,
				"frequency": subscription.Frequency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.SubscribeTxParams{
					CreateSubscriptionParams: db.CreateSubscriptionParams{
						Email:     subscription.Email,
						City:      subscription.City,
						Frequency: subscription.Frequency,
					},
				}

				store.EXPECT().
					SubscribeTx(gomock.Any(), subscribeTxParamsMatcher{expected: arg}).
					Times(1).
					Return(subscription, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSubscription(t, recorder.Body, subscription)
			},
		},
		{
			name: "Email already subscribed",
			body: gin.H{
				"email":     subscription.Email,
				"city":      subscription.City,
				"frequency": subscription.Frequency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				uniqueViolationErr := &pq.Error{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				}
				arg := db.SubscribeTxParams{
					CreateSubscriptionParams: db.CreateSubscriptionParams{
						Email:     subscription.Email,
						City:      subscription.City,
						Frequency: subscription.Frequency,
					},
				}
				store.EXPECT().
					SubscribeTx(gomock.Any(), subscribeTxParamsMatcher{expected: arg}).
					Times(1).
					Return(db.Subscription{}, uniqueViolationErr)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "DB Error",
			body: gin.H{
				"email":     subscription.Email,
				"city":      subscription.City,
				"frequency": subscription.Frequency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.SubscribeTxParams{
					CreateSubscriptionParams: db.CreateSubscriptionParams{
						Email:     subscription.Email,
						City:      subscription.City,
						Frequency: subscription.Frequency,
					},
				}
				store.EXPECT().
					SubscribeTx(gomock.Any(), subscribeTxParamsMatcher{expected: arg}).
					Times(1).
					Return(db.Subscription{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid input email",
			body: gin.H{
				"email":     "invalid",
				"city":      subscription.City,
				"frequency": subscription.Frequency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SubscribeTx(gomock.Any(), gomock.Any()).
					Times(0).
					Return(subscription, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid input city",
			body: gin.H{
				"email":     subscription.Email,
				"city":      "",
				"frequency": subscription.Frequency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SubscribeTx(gomock.Any(), gomock.Any()).
					Times(0).
					Return(subscription, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid input frequency",
			body: gin.H{
				"email":     subscription.Email,
				"city":      subscription.City,
				"frequency": "incorect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SubscribeTx(gomock.Any(), gomock.Any()).
					Times(0).
					Return(subscription, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			taskDistributor := mockworker.NewMockTaskDistributor(ctrl)
			weatherService := mockservice.NewMockService(ctrl)
			tc.buildStubs(store)

			server := NewServer(store, taskDistributor, weatherService)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/subscribe"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomFrequency() db.FrequencyEnum {
	avalFreqEnum := db.AllFrequencyEnumValues()
	i := util.RandomInt(0, int64(len(avalFreqEnum)-1))
	return avalFreqEnum[i]
}

func randomSubscription() db.Subscription {
	return db.Subscription{
		ID:        util.RandomInt(1, 1000),
		Email:     util.RandomEmail(),
		City:      util.RandomCity(),
		Frequency: randomFrequency(),
	}
}

func requireBodyMatchSubscription(t *testing.T, body *bytes.Buffer, subscription db.Subscription) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotSubscription db.Subscription
	err = json.Unmarshal(data, &gotSubscription)
	require.NoError(t, err)
	require.Equal(t, subscription, gotSubscription)
}
