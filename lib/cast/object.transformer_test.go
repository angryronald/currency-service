package cast

import (
	"reflect"
	"testing"
)

func TestTransformObject(t *testing.T) {
	// Define sample input and output structs
	type SourceStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	type ResultStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	// Create a sample source object
	source := SourceStruct{
		Field1: "value1",
		Field2: 42,
	}

	// Create a result object with the same structure
	var result ResultStruct

	// Test TransformObject function
	err := TransformObject(source, &result)
	if err != nil {
		t.Errorf("TransformObject returned an error: %v", err)
	}

	// Check if the result matches the expected output
	expectedResult := ResultStruct{
		Field1: "value1",
		Field2: 42,
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("TransformObject result doesn't match the expected result.")
	}
}

func TestFromBytes(t *testing.T) {
	// Define a sample JSON
	jsonData := []byte(`{"Field1": "value1", "Field2": 42}`)

	// Create a result struct
	var result struct {
		Field1 string `json:"Field1"`
		Field2 int    `json:"Field2"`
	}

	// Test FromBytes function
	err := FromBytes(jsonData, &result)
	if err != nil {
		t.Errorf("FromBytes returned an error: %v", err)
	}

	// Check if the result matches the expected output
	expectedResult := struct {
		Field1 string `json:"Field1"`
		Field2 int    `json:"Field2"`
	}{
		Field1: "value1",
		Field2: 42,
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("FromBytes result doesn't match the expected result.")
	}
}

func TestToBytes(t *testing.T) {
	// Define a sample struct
	data := struct {
		Field1 string `json:"Field1"`
		Field2 int    `json:"Field2"`
	}{
		Field1: "value1",
		Field2: 42,
	}

	// Test ToBytes function
	jsonData, err := ToBytes(data)
	if err != nil {
		t.Errorf("ToBytes returned an error: %v", err)
	}

	// Check if the JSON data matches the expected JSON
	expectedJSON := []byte(`{"Field1":"value1","Field2":42}`)

	if !reflect.DeepEqual(jsonData, expectedJSON) {
		t.Errorf("ToBytes result doesn't match the expected JSON data.")
	}
}
