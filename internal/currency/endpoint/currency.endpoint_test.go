package endpoint

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/angryronald/currency-service/internal/currency/application/command"
	applicationModel "github.com/angryronald/currency-service/internal/currency/application/model"
	"github.com/angryronald/currency-service/internal/currency/application/query"
	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/endpoint/model"
	"github.com/angryronald/currency-service/lib/cast"
	internalHttp "github.com/angryronald/currency-service/lib/http"
	httpUtils "github.com/angryronald/currency-service/lib/test/http"
)

func TestRegisterRoute(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("ALLOWED_CLIENTS", "{\"some_secret_key\":\"some_client\"}")

	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := NewCurrencyEndpoint(currencyCommand, currencyQuery, logrus.New())

	// Create a chi.Mux router
	r := chi.NewRouter()

	// Register routes using RegisterRoute
	r = endpoint.RegisterRoute(r)

	// Perform assertions on the router to ensure routes are registered correctly
	// For example, check if the expected routes exist in the router
	// Check if the expected routes exist in the router
	if !httpUtils.IsRouteExists(r, "/currencies") {
		t.Errorf("Expected route '%s' not found in the router", "/currencies")
	}

	if !httpUtils.IsRouteExists(r, "/{code}") {
		t.Errorf("Expected route '%s' not found in the router", "/{code}")
	}
}

// should return a CurrencyEndpointInterface object
func TestNewCurrencyEndpoint_ReturnsCurrencyEndpointInterfaceObject(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Act
	result := NewCurrencyEndpoint(currencyCommand, currencyQuery, logrus.New())

	// Assert
	assert.Implements(t, (*CurrencyEndpointInterface)(nil), result)
}

// should set the currencyCommand and currencyQuery fields of the CurrencyEndpoint object
func TestNewCurrencyEndpoint_SetsCurrencyCommandAndCurrencyQueryFields(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Act
	result := NewCurrencyEndpoint(currencyCommand, currencyQuery, logrus.New())

	// Assert
	assert.Equal(t, currencyCommand, result.(*CurrencyEndpoint).currencyCommand)
	assert.Equal(t, currencyQuery, result.(*CurrencyEndpoint).currencyQuery)
}

func TestListCurrencies(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := CurrencyEndpoint{currencyCommand, currencyQuery, logrus.New()}

	currenciesMock := []*applicationModel.CurrencyApplicationModel{
		{
			Name: "Indonesian Rupiah",
			Code: "IDR",
		},
		{
			Name: "United States of America Dollar",
			Code: "USD",
		},
		{
			Name: "Euro",
			Code: "EUR",
		},
		{
			Name: "Thailand Baht",
			Code: "THB",
		},
	}

	expectedCurrencies := []*model.CurrencyResponse{
		{
			Name: "Indonesian Rupiah",
			Code: "IDR",
		},
		{
			Name: "United States of America Dollar",
			Code: "USD",
		},
		{
			Name: "Euro",
			Code: "EUR",
		},
		{
			Name: "Thailand Baht",
			Code: "THB",
		},
	}

	currencyQuery.EXPECT().
		List(gomock.Any()).
		Return(currenciesMock, nil).
		Times(1)

	// Create a request with a valid payload
	req, err := http.NewRequest("GET", "/currencies", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the GenerateOTP handler
	handler := http.HandlerFunc(endpoint.listCurrencies)
	handler.ServeHTTP(rr, req)

	// Check the response status code (200 OK for this test)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ListCurrencies returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Parse the response body
	var response internalHttp.ResponseModel
	if err := cast.FromBytes(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error parsing response JSON: %v", err)
	}

	actualCurrencies := []*model.CurrencyResponse{}
	cast.TransformObject(response.Data, &actualCurrencies)
	// Add assertions to validate the response body and behavior as needed
	// Example assertion: assert response.AuthenticationType matches the expected value.
	if !reflect.DeepEqual(actualCurrencies, expectedCurrencies) {
		t.Errorf("got %v, want %v", actualCurrencies, expectedCurrencies)
	}
	if len(actualCurrencies) != len(expectedCurrencies) {
		t.Errorf("got %v, want %v", len(actualCurrencies), len(expectedCurrencies))
	}
}

func TestListCurrencies_QueryError(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := CurrencyEndpoint{currencyCommand, currencyQuery, logrus.New()}
	tests := []struct {
		name           string
		expectedErr    error
		mockHttpStatus int
	}{
		{
			name:           "ListCurrencies not found",
			expectedErr:    constant.ErrNotFound,
			mockHttpStatus: http.StatusNoContent,
		},
		{
			name:           "ListCurrencies query error",
			expectedErr:    errors.New("something went wrong"),
			mockHttpStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			currencyQuery.EXPECT().
				List(gomock.Any()).
				Return(nil, test.expectedErr).
				Times(1)

			// Create a request with a valid payload
			req, err := http.NewRequest("GET", "/currencies", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Call the GenerateOTP handler
			handler := http.HandlerFunc(endpoint.listCurrencies)
			handler.ServeHTTP(rr, req)

			// Check the response status code (500 Internal Server Error for command error)
			if status := rr.Code; status != test.mockHttpStatus {
				t.Errorf("ListCurrencies returned wrong status code: got %v, want %v", status, test.mockHttpStatus)
			}
		})
	}
}

func TestGetCurrency(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("ALLOWED_CLIENTS", "{\"some_secret_key\":\"some_client\"}")

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := CurrencyEndpoint{currencyCommand, currencyQuery, logrus.New()}

	currencyCode := "IDR"
	currencyMock := &applicationModel.CurrencyApplicationModel{
		Name: "Indonesian Rupiah",
		Code: "IDR",
	}

	expectedCurrency := &model.CurrencyResponse{
		Name: "Indonesian Rupiah",
		Code: "IDR",
	}

	currencyQuery.EXPECT().
		GetByCode(gomock.Any(), currencyCode).
		Return(currencyMock, nil).
		Times(1)

	// Create a request with a valid payload
	req, err := http.NewRequest("GET", fmt.Sprintf("/currencies/%s", currencyCode), nil)
	if err != nil {
		t.Fatal(err)
	}

	r := chi.NewRouter()
	r = endpoint.RegisterRoute(r)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the GenerateOTP handler
	r.ServeHTTP(rr, req)

	// Check the response status code (200 OK for this test)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetCurrency returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Parse the response body
	var response internalHttp.ResponseModel
	if err := cast.FromBytes(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error parsing response JSON: %v", err)
	}

	actualCurrency := &model.CurrencyResponse{}
	cast.TransformObject(response.Data, actualCurrency)
	// Add assertions to validate the response body and behavior as needed
	// Example assertion: assert response.AuthenticationType matches the expected value.
	if !reflect.DeepEqual(actualCurrency, expectedCurrency) {
		t.Errorf("got %v, want %v", actualCurrency, expectedCurrency)
	}
}

func TestGetCurrency_QueryError(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("ALLOWED_CLIENTS", "{\"some_secret_key\":\"some_client\"}")

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := CurrencyEndpoint{currencyCommand, currencyQuery, logrus.New()}

	currencyCode := "IDR"

	tests := []struct {
		name           string
		expectedErr    error
		mockHttpStatus int
	}{
		{
			name:           "GetCurrency not found",
			expectedErr:    constant.ErrNotFound,
			mockHttpStatus: http.StatusBadRequest,
		},
		{
			name:           "GetCurrency query error",
			expectedErr:    errors.New("something went wrong"),
			mockHttpStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			currencyQuery.EXPECT().
				GetByCode(gomock.Any(), currencyCode).
				Return(nil, test.expectedErr).
				Times(1)

			// Create a request with a valid payload
			req, err := http.NewRequest("GET", fmt.Sprintf("/currencies/%s", currencyCode), nil)
			if err != nil {
				t.Fatal(err)
			}

			r := chi.NewRouter()
			r = endpoint.RegisterRoute(r)

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Call the GenerateOTP handler
			r.ServeHTTP(rr, req)

			// Check the response status code (200 OK for this test)
			if status := rr.Code; status != test.mockHttpStatus {
				t.Errorf("GetCurrency returned wrong status code: got %v, want %v", status, test.mockHttpStatus)
			}
		})
	}
}

func TestAddCurrency(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("ALLOWED_CLIENTS", "{\"some_secret_key\":\"some_client\"}")

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := CurrencyEndpoint{currencyCommand, currencyQuery, logrus.New()}

	currencyMock := &applicationModel.CurrencyApplicationModel{
		Name: "Euro",
		Code: "EUR",
	}

	expectedCurrency := &model.CurrencyResponse{
		Name: "Euro",
		Code: "EUR",
	}

	currencyCommand.EXPECT().
		Add(gomock.Any(), currencyMock).
		Return(currencyMock, nil).
		Times(1)

	// Create a request with a valid payload
	payload := []byte(`{"Name": "Euro", "Code": "EUR"}`)
	req, err := http.NewRequest("POST", "/currencies", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	r := chi.NewRouter()
	r = endpoint.RegisterRoute(r)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the GenerateOTP handler
	r.ServeHTTP(rr, req)

	// Check the response status code (200 OK for this test)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AddCurrency returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Parse the response body
	var response internalHttp.ResponseModel
	if err := cast.FromBytes(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error parsing response JSON: %v", err)
	}

	actualCurrency := &model.CurrencyResponse{}
	cast.TransformObject(response.Data, actualCurrency)
	// Add assertions to validate the response body and behavior as needed
	// Example assertion: assert response.AuthenticationType matches the expected value.
	if !reflect.DeepEqual(actualCurrency, expectedCurrency) {
		t.Errorf("got %v, want %v", actualCurrency, expectedCurrency)
	}
}

func TestAddCurrency_CommandError(t *testing.T) {
	// Create a mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("ALLOWED_CLIENTS", "{\"some_secret_key\":\"some_client\"}")

	// Arrange
	currencyCommand := command.NewMockCurrencyCommandInterface(ctrl)
	currencyQuery := query.NewMockCurrencyQueryInterface(ctrl)

	// Initialize your AuthenticationEndpoint with mock dependencies
	endpoint := CurrencyEndpoint{currencyCommand, currencyQuery, logrus.New()}

	tests := []struct {
		name             string
		currencyShouldBe *applicationModel.CurrencyApplicationModel
		payload          []byte
		expectedErr      error
		mockHttpStatus   int
	}{
		{
			name: "AddCurrency not found",
			currencyShouldBe: &applicationModel.CurrencyApplicationModel{
				Name: "Euro",
				Code: "EUR",
			},
			payload:        []byte(`{"Name": "Euro", "Code": "EUR"}`),
			expectedErr:    constant.ErrConflict,
			mockHttpStatus: http.StatusConflict,
		},
		{
			name:           "AddCurrency command error",
			expectedErr:    errors.New("something went wrong"),
			mockHttpStatus: http.StatusInternalServerError,
			currencyShouldBe: &applicationModel.CurrencyApplicationModel{
				Name: "Euro",
				Code: "EUR",
			},
			payload: []byte(`{"Name": "Euro", "Code": "EUR"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			currencyCommand.EXPECT().
				Add(gomock.Any(), test.currencyShouldBe).
				Return(nil, test.expectedErr).
				Times(1)

			// Create a request with a valid payload
			req, err := http.NewRequest("POST", "/currencies", bytes.NewBuffer(test.payload))
			if err != nil {
				t.Fatal(err)
			}

			r := chi.NewRouter()
			r = endpoint.RegisterRoute(r)

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Call the GenerateOTP handler
			r.ServeHTTP(rr, req)

			// Check the response status code (200 OK for this test)
			if status := rr.Code; status != test.mockHttpStatus {
				t.Errorf("AddCurrency returned wrong status code: got %v, want %v", status, test.mockHttpStatus)
			}
		})
	}
}
