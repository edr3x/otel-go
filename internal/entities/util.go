package entities

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func addMapToSpan(span trace.Span, data map[string]any, prefix string) {
	if data == nil {
		return
	}
	var attrs []attribute.KeyValue
	for k, v := range data {
		attrKey := prefix + "." + k
		switch val := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(attrKey, val))
		case int:
			attrs = append(attrs, attribute.Int(attrKey, val))
		case int64:
			attrs = append(attrs, attribute.Int64(attrKey, val))
		case float64:
			attrs = append(attrs, attribute.Float64(attrKey, val))
		case bool:
			attrs = append(attrs, attribute.Bool(attrKey, val))
		default:
			jsonValue, _ := json.Marshal(val) // Convert unknown types to JSON
			attrs = append(attrs, attribute.String(attrKey, string(jsonValue)))
		}
	}
	span.SetAttributes(attrs...)
}
