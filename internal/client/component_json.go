package fractalCloud

import "encoding/json"

// MarshalJSON sends parameters as typed JSON; see encodeWireValue.
func (c Component) MarshalJSON() ([]byte, error) {
	type plain Component
	return json.Marshal(struct {
		plain
		Parameters map[string]json.RawMessage `json:"parameters,omitempty"`
	}{
		plain:      plain(c),
		Parameters: encodeWireMap(c.Parameters),
	})
}

// MarshalJSON sends link settings as typed JSON; see encodeWireValue.
func (l ComponentLink) MarshalJSON() ([]byte, error) {
	type plain ComponentLink
	return json.Marshal(struct {
		plain
		Settings map[string]json.RawMessage `json:"settings"`
	}{
		plain:    plain(l),
		Settings: encodeWireMap(l.Settings),
	})
}
