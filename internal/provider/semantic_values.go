package provider

import (
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
)

// Structured parameter and link-setting values are sent to the API as JSON
// and come back re-encoded: object keys sorted, whitespace dropped. The
// string in the configuration is the source of truth, so a value the API
// returns that is the same JSON as the prior state's string keeps the prior
// string; anything else is a real change and is taken from the API.

// preferPriorValues returns current with every value that is JSON-equal to
// the prior value under the same key replaced by the prior value.
func preferPriorValues(prior, current map[string]string) map[string]string {
	if len(prior) == 0 {
		return current
	}
	out := make(map[string]string, len(current))
	for k, v := range current {
		if p, ok := prior[k]; ok && sameJSON(p, v) {
			out[k] = p
			continue
		}
		out[k] = v
	}
	return out
}

// sameJSON reports whether a and b are identical, or are both a JSON object
// or array with the same content.
func sameJSON(a, b string) bool {
	if a == b {
		return true
	}
	av, ok := structuredJSON(a)
	if !ok {
		return false
	}
	bv, ok := structuredJSON(b)
	if !ok {
		return false
	}
	return reflect.DeepEqual(canonicalNumbers(av), canonicalNumbers(bv))
}

// canonicalNumbers replaces every JSON number with its exact rational value,
// so 1, 1.0 and 1e0 compare equal.
func canonicalNumbers(v any) any {
	switch t := v.(type) {
	case json.Number:
		if r, ok := new(big.Rat).SetString(t.String()); ok {
			return r.RatString()
		}
		return t.String()
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = canonicalNumbers(e)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, e := range t {
			out[k] = canonicalNumbers(e)
		}
		return out
	default:
		return v
	}
}

func structuredJSON(s string) (any, bool) {
	trimmed := bytes.TrimSpace([]byte(s))
	if len(trimmed) == 0 || (trimmed[0] != '{' && trimmed[0] != '[') {
		return nil, false
	}
	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.UseNumber()
	var v any
	if err := decoder.Decode(&v); err != nil {
		return nil, false
	}
	if decoder.More() {
		return nil, false
	}
	return v, true
}
