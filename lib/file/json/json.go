package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadJSON(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
