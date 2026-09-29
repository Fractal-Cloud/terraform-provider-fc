package fractalCloud

import (
	"bytes"
	"encoding/json"
)

// Parameter and link-setting values are strings in Terraform, because a
// component's parameters are a map(string). The API types them: an array or
// object parameter (sparkConf, nodePools, an environment-secret reference
// {"$envSecret": "..."}) must arrive as JSON, not as a string holding JSON.
//
// The convention is symmetric with mapAnyToMapStringJSON on the read path:
// a value that is a JSON object or array travels as that JSON, and every
// other value travels as the string it is. Scalars are deliberately left as
// strings, so a parameter that is the text "true" or "8080" is never
// reinterpreted.

// encodeWireValue returns the JSON a parameter or link-setting value is sent as.
func encodeWireValue(v string) json.RawMessage {
	trimmed := bytes.TrimSpace([]byte(v))
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') && json.Valid(trimmed) {
		return json.RawMessage(trimmed)
	}
	quoted, _ := json.Marshal(v) // marshaling a string cannot fail
	return quoted
}

// encodeWireMap applies encodeWireValue to every value of in.
func encodeWireMap(in map[string]string) map[string]json.RawMessage {
	if in == nil {
		return nil
	}
	out := make(map[string]json.RawMessage, len(in))
	for k, v := range in {
		out[k] = encodeWireValue(v)
	}
	return out
}
