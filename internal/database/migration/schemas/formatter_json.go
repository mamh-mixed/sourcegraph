package schemas

import (
	"encoding/json"
)

type jsonFormatter struct{}

func NewJSONFormatter() SchemaFormatter {
	return jsonFormatter{}
}

func (f jsonFormatter) Format(schemaDescription SchemaDescription) string {
	serialized, _ := json.MarshalIndent(schemaDescription, "", "    ")
	return string(serialized)
}
