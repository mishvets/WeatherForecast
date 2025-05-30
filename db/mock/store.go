// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mishvets/WeatherForecast/db/sqlc (interfaces: Store)
//
// Generated by this command:
//
//	mockgen --package mockdb --destination db/mock/store.go github.com/mishvets/WeatherForecast/db/sqlc Store
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
	isgomock struct{}
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// ConfirmSubscription mocks base method.
func (m *MockStore) ConfirmSubscription(ctx context.Context, arg db.ConfirmSubscriptionParams) (db.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmSubscription", ctx, arg)
	ret0, _ := ret[0].(db.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfirmSubscription indicates an expected call of ConfirmSubscription.
func (mr *MockStoreMockRecorder) ConfirmSubscription(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmSubscription", reflect.TypeOf((*MockStore)(nil).ConfirmSubscription), ctx, arg)
}

// ConfirmSubscriptionTx mocks base method.
func (m *MockStore) ConfirmSubscriptionTx(ctx context.Context, arg db.ConfirmSubscriptionTxParams) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmSubscriptionTx", ctx, arg)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfirmSubscriptionTx indicates an expected call of ConfirmSubscriptionTx.
func (mr *MockStoreMockRecorder) ConfirmSubscriptionTx(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmSubscriptionTx", reflect.TypeOf((*MockStore)(nil).ConfirmSubscriptionTx), ctx, arg)
}

// CreateNewWeatherTx mocks base method.
func (m *MockStore) CreateNewWeatherTx(ctx context.Context, arg db.CreateNewWeatherTxParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewWeatherTx", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewWeatherTx indicates an expected call of CreateNewWeatherTx.
func (mr *MockStoreMockRecorder) CreateNewWeatherTx(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewWeatherTx", reflect.TypeOf((*MockStore)(nil).CreateNewWeatherTx), ctx, arg)
}

// CreateSubscription mocks base method.
func (m *MockStore) CreateSubscription(ctx context.Context, arg db.CreateSubscriptionParams) (db.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSubscription", ctx, arg)
	ret0, _ := ret[0].(db.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubscription indicates an expected call of CreateSubscription.
func (mr *MockStoreMockRecorder) CreateSubscription(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubscription", reflect.TypeOf((*MockStore)(nil).CreateSubscription), ctx, arg)
}

// CreateWeather mocks base method.
func (m *MockStore) CreateWeather(ctx context.Context, arg db.CreateWeatherParams) (db.WeatherDatum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWeather", ctx, arg)
	ret0, _ := ret[0].(db.WeatherDatum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWeather indicates an expected call of CreateWeather.
func (mr *MockStoreMockRecorder) CreateWeather(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWeather", reflect.TypeOf((*MockStore)(nil).CreateWeather), ctx, arg)
}

// DeleteSubscription mocks base method.
func (m *MockStore) DeleteSubscription(ctx context.Context, token uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubscription", ctx, token)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSubscription indicates an expected call of DeleteSubscription.
func (mr *MockStoreMockRecorder) DeleteSubscription(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscription", reflect.TypeOf((*MockStore)(nil).DeleteSubscription), ctx, token)
}

// DeleteSubscriptionTx mocks base method.
func (m *MockStore) DeleteSubscriptionTx(ctx context.Context, arg db.DeleteSubscriptionTxParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubscriptionTx", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubscriptionTx indicates an expected call of DeleteSubscriptionTx.
func (mr *MockStoreMockRecorder) DeleteSubscriptionTx(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscriptionTx", reflect.TypeOf((*MockStore)(nil).DeleteSubscriptionTx), ctx, arg)
}

// DeleteWeather mocks base method.
func (m *MockStore) DeleteWeather(ctx context.Context, city string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWeather", ctx, city)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWeather indicates an expected call of DeleteWeather.
func (mr *MockStoreMockRecorder) DeleteWeather(ctx, city any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWeather", reflect.TypeOf((*MockStore)(nil).DeleteWeather), ctx, city)
}

// GetCitiesForUpdate mocks base method.
func (m *MockStore) GetCitiesForUpdate(ctx context.Context, frequency db.FrequencyEnum) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCitiesForUpdate", ctx, frequency)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCitiesForUpdate indicates an expected call of GetCitiesForUpdate.
func (mr *MockStoreMockRecorder) GetCitiesForUpdate(ctx, frequency any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCitiesForUpdate", reflect.TypeOf((*MockStore)(nil).GetCitiesForUpdate), ctx, frequency)
}

// GetEmailsForUpdate mocks base method.
func (m *MockStore) GetEmailsForUpdate(ctx context.Context, arg db.GetEmailsForUpdateParams) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailsForUpdate", ctx, arg)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmailsForUpdate indicates an expected call of GetEmailsForUpdate.
func (mr *MockStoreMockRecorder) GetEmailsForUpdate(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailsForUpdate", reflect.TypeOf((*MockStore)(nil).GetEmailsForUpdate), ctx, arg)
}

// GetSubscription mocks base method.
func (m *MockStore) GetSubscription(ctx context.Context, email string) (db.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscription", ctx, email)
	ret0, _ := ret[0].(db.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscription indicates an expected call of GetSubscription.
func (mr *MockStoreMockRecorder) GetSubscription(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscription", reflect.TypeOf((*MockStore)(nil).GetSubscription), ctx, email)
}

// GetSubscriptionForUpdate mocks base method.
func (m *MockStore) GetSubscriptionForUpdate(ctx context.Context, token uuid.UUID) (db.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionForUpdate", ctx, token)
	ret0, _ := ret[0].(db.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionForUpdate indicates an expected call of GetSubscriptionForUpdate.
func (mr *MockStoreMockRecorder) GetSubscriptionForUpdate(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionForUpdate", reflect.TypeOf((*MockStore)(nil).GetSubscriptionForUpdate), ctx, token)
}

// GetWeather mocks base method.
func (m *MockStore) GetWeather(ctx context.Context, city string) (db.WeatherDatum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeather", ctx, city)
	ret0, _ := ret[0].(db.WeatherDatum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWeather indicates an expected call of GetWeather.
func (mr *MockStoreMockRecorder) GetWeather(ctx, city any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeather", reflect.TypeOf((*MockStore)(nil).GetWeather), ctx, city)
}

// GetWeatherForUpdate mocks base method.
func (m *MockStore) GetWeatherForUpdate(ctx context.Context, city string) (db.WeatherDatum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeatherForUpdate", ctx, city)
	ret0, _ := ret[0].(db.WeatherDatum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWeatherForUpdate indicates an expected call of GetWeatherForUpdate.
func (mr *MockStoreMockRecorder) GetWeatherForUpdate(ctx, city any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeatherForUpdate", reflect.TypeOf((*MockStore)(nil).GetWeatherForUpdate), ctx, city)
}

// IsSubscriptionExist mocks base method.
func (m *MockStore) IsSubscriptionExist(ctx context.Context, id int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSubscriptionExist", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsSubscriptionExist indicates an expected call of IsSubscriptionExist.
func (mr *MockStoreMockRecorder) IsSubscriptionExist(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSubscriptionExist", reflect.TypeOf((*MockStore)(nil).IsSubscriptionExist), ctx, id)
}

// SubscribeTx mocks base method.
func (m *MockStore) SubscribeTx(ctx context.Context, arg db.SubscribeTxParams) (db.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeTx", ctx, arg)
	ret0, _ := ret[0].(db.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeTx indicates an expected call of SubscribeTx.
func (mr *MockStoreMockRecorder) SubscribeTx(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeTx", reflect.TypeOf((*MockStore)(nil).SubscribeTx), ctx, arg)
}

// UpdateWeather mocks base method.
func (m *MockStore) UpdateWeather(ctx context.Context, arg db.UpdateWeatherParams) (db.WeatherDatum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWeather", ctx, arg)
	ret0, _ := ret[0].(db.WeatherDatum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWeather indicates an expected call of UpdateWeather.
func (mr *MockStoreMockRecorder) UpdateWeather(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWeather", reflect.TypeOf((*MockStore)(nil).UpdateWeather), ctx, arg)
}
