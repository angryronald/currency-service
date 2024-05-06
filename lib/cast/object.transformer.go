package cast

import (
	"encoding/json"
)

// TransformObject used to transform source object to result object based on json tag
func TransformObject(source interface{}, result interface{}) error {
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}

	err = json.Unmarshal(sourceBytes, &result)
	if err != nil {
		return err
	}

	return nil
}

func FromBytes(source []byte, result interface{}) error {
	return json.Unmarshal(source, result)
}

func ToBytes(source interface{}) ([]byte, error) {
	return json.Marshal(source)
}
