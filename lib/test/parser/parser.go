package parser

import (
	"encoding/json"

	"github.com/angryronald/currency-service/lib/http"
)

func ParseSuccessObjectToBytes(source interface{}) ([]byte, error) {
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	var result http.ResponseModel

	err = json.Unmarshal(sourceBytes, &result)
	if err != nil {
		return nil, err
	}

	destinationBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, err
	}

	return destinationBytes, nil
}
