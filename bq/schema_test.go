package bq

import (
	"encoding/json"
	"reflect"
	"testing"

	"cloud.google.com/go/bigquery"
)

func TestConvertSchema(t *testing.T) {
	type Item struct {
		Name     string
		Size     float64
		Count    int    `bigquery:"number"`
		Secret   []byte `bigquery:"-"`
		Optional bigquery.NullBool
		OptBytes []byte `bigquery:",nullable"`
	}
	schema, err := bigquery.InferSchema(Item{})
	if err != nil {
		t.Fatalf("when preparing test, InferSchema failed because of error = %v", err)
	}

	schemaOuter := bigquery.Schema{
		{Name: "Name", Required: true, Type: bigquery.StringFieldType},
		{Name: "Grades", Repeated: true, Type: bigquery.IntegerFieldType},
		{Name: "Optional", Required: false, Type: bigquery.IntegerFieldType},
		{Name: "Nest", Type: bigquery.RecordFieldType, Schema: schema},
	}

	_, err = convertSchema(schemaOuter)
	if err != nil {
		t.Fatalf("convertSchema() failed because of error = %v", err)
	}
}

func TestGetSchemaJSON(t *testing.T) {
	t.Run("getSchemaJSON should return []byte and nil error", func(t *testing.T) {
		fs := []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Mode string `json:"mode"`
		}{
			{Name: "Name", Type: "STRING", Mode: "REQUIRED"},
		}

		want, err := json.MarshalIndent(fs, "", "  ")
		if err != nil {
			t.Fatalf("when preparing test, InferSchema failed because of error = %v", err)
		}

		schema := bigquery.Schema{
			{Name: "Name", Required: true, Type: bigquery.StringFieldType},
		}
		got, err := getSchemaJSON(schema)

		if err != nil {
			t.Fatalf("getSchemaJSON() failed with error = %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("getSchemaJSON() = %v, want %v", got, want)
		}
	})
}
