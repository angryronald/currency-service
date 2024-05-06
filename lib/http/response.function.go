package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Response(w http.ResponseWriter, statusCode int, data interface{}, log *logrus.Logger) {
	response := &ResponseModel{}
	response.Code = fmt.Sprint(statusCode)
	response.Message = http.StatusText(statusCode)
	response.Status = http.StatusText(statusCode)
	response.Data = data

	w.WriteHeader(statusCode)
	var resultInBytes []byte
	var err error
	if resultInBytes, err = json.Marshal(response); err != nil {
		log.Warnf(`[ResponseError] error: %s StackTrace: %v`, err, errors.WithStack(err))

		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(resultInBytes)
}

func ResponseError(w http.ResponseWriter, statusCode int, message string, log *logrus.Logger) {
	response := &ResponseModel{}
	response.Code = fmt.Sprint(statusCode)
	response.Message = message
	response.Status = http.StatusText(statusCode)

	w.WriteHeader(statusCode)
	var resultInBytes []byte
	var err error
	if resultInBytes, err = json.Marshal(response); err != nil {
		log.Warnf(`[ResponseError] error: %s StackTrace: %v`, err, errors.WithStack(err))

		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(resultInBytes)
}
