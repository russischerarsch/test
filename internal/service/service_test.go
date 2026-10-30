package service_test

import (
	"PlataTest/internal/domain"
	"PlataTest/internal/service"
	"PlataTest/internal/service/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateQuote_success(t *testing.T) {
	m := mocks.NewRepository(t)
	m.On("CreateQuote", mock.Anything, "eurusd").Return(1, nil)

	svc := service.NewService(m)

	id, err := svc.CreateQuote(t.Context(), "EURUSD")

	require.NoError(t, err)

	assert.Equal(t, 1, id)
}

func TestCreateQuote_repoError(t *testing.T) {
	m := mocks.NewRepository(t)
	svc := service.NewService(m)

	repoErr := errors.New("db down")

	m.On("CreateQuote", mock.Anything, "eurusd").Return(0, repoErr).Once()

	id, err := svc.CreateQuote(t.Context(), "EURUSD")

	require.ErrorIs(t, err, repoErr)
	assert.Zero(t, id)
}

func TestCreateQuote_invalidPair(t *testing.T) {
	tests := []struct {
		name string
		pair string
	}{
		{name: "empty", pair: ""},
		{name: "too short", pair: "EUR"},
		{name: "too long", pair: "EURUSDT"},
		{name: "digits", pair: "EUR123"},
		{name: "special chars", pair: "EUR/USD"},
		{name: "spaces", pair: "EUR US"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			m := mocks.NewRepository(t)
			svc := service.NewService(m)

			id, err := svc.CreateQuote(t.Context(), testCase.pair)

			require.ErrorIs(t, err, domain.ErrInvalidPairFormat)
			assert.Zero(t, id)
		})
	}
}

func TestGetLatestQuote_success(t *testing.T) {
	m := mocks.NewRepository(t)
	svc := service.NewService(m)

	want := &domain.Quote{ID: 1, Pair: "eurusd"}

	m.On("GetLatestQuote", mock.Anything, "eurusd").Return(want, nil).Once()

	got, err := svc.GetLatestQuote(t.Context(), "EURUSD")

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetLatestQuote_invalidPairFormat(t *testing.T) {
	tests := []struct {
		name string
		pair string
	}{
		{name: "empty", pair: ""},
		{name: "too short", pair: "eurus"},
		{name: "too long", pair: "eurusdt"},
		{name: "digit", pair: "eur1sd"},
		{name: "space", pair: "eur sd"},
		{name: "symbol", pair: "eur-sd"},
		{name: "cyrillic", pair: "евроус"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewRepository(t)
			svc := service.NewService(m)

			got, err := svc.GetLatestQuote(t.Context(), tt.pair)

			require.ErrorIs(t, err, domain.ErrInvalidPairFormat)
			assert.Nil(t, got)

		})
	}
}
