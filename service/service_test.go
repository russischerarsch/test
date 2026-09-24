package service

import (
	"PlataTest/service/mocks"
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchPrice(t *testing.T) {
	repository := mocks.NewQuote(t)
	service := CreateService(repository)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// url := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=%s&base=%s&symbols=%s",
	// 	"84343d3b62c279fb5c47f52fbdf034ae",
	// 	"EUR",
	// 	"USD",
	// )
	httpmock.RegisterResponder(http.MethodGet,
		"http://api.exchangeratesapi.io/v1/latest",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"success":   true,
			"timestamp": 1790177284,
			"base":      "EUR",
			"date":      "2026-09-23",
			"rates": map[string]interface{}{
				"USD": 1.139367,
			},
		}))
	price, err := service.FetchPrice(context.Background(), "EURUSD")
	require.NoError(t, err)
	assert.Equal(t, 1.139367, price)
}
