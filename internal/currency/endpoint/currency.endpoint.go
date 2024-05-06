package endpoint

import (
	"encoding/json"
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/config"
	"github.com/angryronald/currency-service/internal/currency/application/command"
	applicationModel "github.com/angryronald/currency-service/internal/currency/application/model"
	"github.com/angryronald/currency-service/internal/currency/application/query"
	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/endpoint/model"
	internalHttp "github.com/angryronald/currency-service/lib/http"
)

var requestGroup = singleflight.Group{}

type CurrencyEndpoint struct {
	currencyCommand command.CurrencyCommandInterface
	currencyQuery   query.CurrencyQueryInterface
	log             *logrus.Logger
}

func (e *CurrencyEndpoint) listCurrencies(w http.ResponseWriter, r *http.Request) {
	var err error
	var v interface{}
	var currencies []*applicationModel.CurrencyApplicationModel

	if v, err, _ = requestGroup.Do("list", func() (interface{}, error) {
		return e.currencyQuery.List(r.Context())
	}); err != nil {
		e.log.Warnf(`[CurrencyEndpoint.addCurrency] error: %s StackTrace: %v`, err, errors.WithStack(err))
		switch err {
		case constant.ErrNotFound:
			internalHttp.ResponseError(w, http.StatusNoContent, "", e.log)
			return

		default:
			internalHttp.ResponseError(w, http.StatusInternalServerError, "", e.log)
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	currencies = v.([]*applicationModel.CurrencyApplicationModel)

	result := []*model.CurrencyResponse{}
	for _, currency := range currencies {
		result = append(result, model.NewCurrencyResponse(currency))
	}
	internalHttp.Response(w, http.StatusOK, result, e.log)
}

func (e *CurrencyEndpoint) getCurrency(w http.ResponseWriter, r *http.Request) {
	var err error
	var v interface{}
	var currency *applicationModel.CurrencyApplicationModel
	currencyCode := chi.URLParam(r, "code")

	if v, err, _ = requestGroup.Do(currencyCode, func() (interface{}, error) {
		return e.currencyQuery.GetByCode(r.Context(), currencyCode)
	}); err != nil {
		e.log.Warnf(`[CurrencyEndpoint.addCurrency] error: %s StackTrace: %v`, err, errors.WithStack(err))
		// error handling
		switch err {
		case constant.ErrNotFound:
			internalHttp.ResponseError(w, http.StatusBadRequest, "", e.log)
			return

		default:
			internalHttp.ResponseError(w, http.StatusInternalServerError, "", e.log)
			return
		}
	}

	currency = v.(*applicationModel.CurrencyApplicationModel)
	result := model.NewCurrencyResponse(currency)
	internalHttp.Response(w, http.StatusOK, result, e.log)
}

func (e *CurrencyEndpoint) addCurrency(w http.ResponseWriter, r *http.Request) {
	var err error
	var currency *applicationModel.CurrencyApplicationModel

	var addCurrencyRequest model.AddCurrencyRequest
	err = json.NewDecoder(r.Body).Decode(&addCurrencyRequest)
	if err != nil {
		e.log.Warnf(`[CurrencyEndpoint.addCurrency] error: %s StackTrace: %v`, err, errors.WithStack(err))
		internalHttp.ResponseError(w, http.StatusBadRequest, "", e.log)
		return
	}

	if currency, err = e.currencyCommand.Add(r.Context(), addCurrencyRequest.ToApplicationModel()); err != nil {
		e.log.Warnf(`[CurrencyEndpoint.addCurrency] error: %s StackTrace: %v`, err, errors.WithStack(err))

		// error handling
		switch err {
		case constant.ErrConflict:
			internalHttp.ResponseError(w, http.StatusConflict, "", e.log)
			return

		default:
			internalHttp.ResponseError(w, http.StatusInternalServerError, "", e.log)
			return
		}
	}

	result := model.NewCurrencyResponse(currency)
	internalHttp.Response(w, http.StatusOK, result, e.log)
}

func (e *CurrencyEndpoint) RegisterRoute(r *chi.Mux) *chi.Mux {
	r.Route("/currencies", func(r chi.Router) {
		r.Get("/", e.listCurrencies)
		r.Get("/{code}", e.getCurrency)
		r.Post("/", e.addCurrency)
	})

	r.Group(func(r chi.Router) {
		// other services
		r.Route("/clients", func(r chi.Router) {
			r.Use(internalHttp.AuthenticateClient(config.GetAllowedClients()))

			r.Route("/currencies", func(r chi.Router) {
				r.Get("/", e.listCurrencies)
				r.Get("/{code}", e.getCurrency)
			})
		})
	})
	return r
}

func NewCurrencyEndpoint(
	currencyCommand command.CurrencyCommandInterface,
	currencyQuery query.CurrencyQueryInterface,
	log *logrus.Logger,
) CurrencyEndpointInterface {
	return &CurrencyEndpoint{
		currencyCommand: currencyCommand,
		currencyQuery:   currencyQuery,
		log:             log,
	}
}
