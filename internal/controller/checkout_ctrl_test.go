package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/learn/api-shop/internal/controller"
	mockSvc "github.com/learn/api-shop/internal/generated/mock"
	"github.com/learn/api-shop/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCheckoutHandler(t *testing.T) {
	tests := []struct {
		name             string
		query            string
		expectedResponse map[string]interface{}
		mockSetupFunc    func(checkoutSvc *mockSvc.CheckoutUsecase)
	}{
		{
			name: "Checkout with one item",
			query: `
				mutation {
					checkout(items: [{product_id: 2, qty: 1}]) {
						items
						total_amount
					}
				}
			`,
			expectedResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"checkout": map[string]interface{}{
						"items":        []interface{}{"Item1"},
						"total_amount": 10.0,
					},
				},
			},
			mockSetupFunc: func(checkoutSvc *mockSvc.CheckoutUsecase) {
				checkoutSvc.On("Checkout", mock.Anything, mock.AnythingOfType("[]repo.OrderDetail")).Return(service.Checkout{
					Items:       []string{"Item1"},
					TotalAmount: 10.0,
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkoutSvc := new(mockSvc.CheckoutUsecase)
			tt.mockSetupFunc(checkoutSvc)

			mux := http.NewServeMux()
			controller.NewCheckoutHandler(mux, checkoutSvc)

			testServer := httptest.NewServer(mux)
			defer testServer.Close()

			requestBody := map[string]interface{}{
				"query": tt.query,
			}

			response, err := http.Post(testServer.URL+"/graphql", "application/json", toJSONRequestBody(requestBody))
			assert.NoError(t, err)
			assert.NotNil(t, response)

			var jsonResponse map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&jsonResponse)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedResponse, jsonResponse)
		})
	}
}

func toJSONRequestBody(data interface{}) io.Reader {
	jsonData, _ := json.Marshal(data)
	return bytes.NewBuffer(jsonData)
}

func TestCreateCheckoutSchema(t *testing.T) {
	testCases := []struct {
		name           string
		requestString  string
		checkoutResult service.Checkout
		expectedData   map[string]interface{}
	}{
		{
			name:          "Checkout with one item",
			requestString: `mutation { checkout(items: [{ product_id: 2, qty: 1 }]) { items total_amount }}`,
			checkoutResult: service.Checkout{
				Items:       []string{"Item1"},
				TotalAmount: 10,
			},
			expectedData: map[string]interface{}{
				"checkout": map[string]interface{}{
					"items":        []interface{}{"Item1"},
					"total_amount": 10.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCheckoutSvc := new(mockSvc.CheckoutUsecase)

			hc := &controller.CheckoutCntrlImpl{
				CheckoutSvc: mockCheckoutSvc,
			}

			schema, err := controller.CreateCheckoutSchema(hc)
			assert.NoError(t, err)

			mockCheckoutSvc.On("Checkout", mock.Anything, mock.Anything).Return(tc.checkoutResult, nil)

			params := graphql.Params{
				Schema:        schema,
				RequestString: tc.requestString,
			}

			result := graphql.Do(params)
			assert.False(t, result.HasErrors())
			assert.Equal(t, tc.expectedData, result.Data)
		})
	}
}
