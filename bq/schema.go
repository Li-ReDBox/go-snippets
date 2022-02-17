package bq

import (
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// Field describes the schema of a field in BigQuery - this is for exporting in JSON.
// https://cloud.google.com/bigquery/docs/schemas#go
// https://github.com/GoogleCloudPlatform/protoc-gen-bq-schema/blob/master/main.go
type Field struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Mode        string   `json:"mode"`
	Description string   `json:"description,omitempty"`
	Fields      []*Field `json:"fields,omitempty"`
}

// convertSchema converts bigquery.Schema to []*bq.Field
func convertSchema(s bigquery.Schema) (schema []*Field, err error) {
	for _, f := range s {
		field, err := convertField(*f)

		if err != nil {
			err = fmt.Errorf("failed to convert field %s : %v", f.Name, err)
			return nil, err
		}

		schema = append(schema, field)
	}

	return
}

// convertField converts bigquery.FieldSchema to []*bq.Field, go deeper if needed
func convertField(fs bigquery.FieldSchema) (*Field, error) {
	field := &Field{
		Name:        fs.Name,
		Description: fs.Description,
		Type:        string(fs.Type),
		Mode:        "NULLABLE",
	}

	if fs.Repeated {
		field.Mode = "REPEATED"
	} else if fs.Required {
		field.Mode = "REQUIRED"
	}

	if field.Type != "RECORD" {
		return field, nil
	}

	fields, err := convertSchema(fs.Schema)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 { // discard RECORDs that would have zero fields
		return nil, nil
	}

	field.Fields = fields

	return field, nil
}

func getSchemaJSON(s bigquery.Schema) ([]byte, error) {
	fs, err := convertSchema(s)
	// TODO: error should say it is in exportSchemaJSON
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(fs, "", "  ")
	if err != nil {
		return nil, err
	}
	return out, nil
}
